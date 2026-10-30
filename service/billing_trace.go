package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

type BillingTraceEvent struct {
	ID                int    `json:"id"`
	Type              string `json:"type"`
	LogType           int    `json:"log_type"`
	QuotaDelta        int    `json:"quota_delta"`
	Quota             int    `json:"quota"`
	RequestID         string `json:"request_id,omitempty"`
	UpstreamRequestID string `json:"upstream_request_id,omitempty"`
	ModelName         string `json:"model_name,omitempty"`
	TokenID           int    `json:"token_id,omitempty"`
	ChannelID         int    `json:"channel_id,omitempty"`
	Content           string `json:"content,omitempty"`
	CreatedAt         int64  `json:"created_at"`
}

type BillingTraceSummary struct {
	BillingTraceID   string              `json:"billing_trace_id"`
	Status           string              `json:"status"`
	ConsumeQuota     int                 `json:"consume_quota"`
	AdjustQuota      int                 `json:"adjust_quota"`
	RefundQuota      int                 `json:"refund_quota"`
	FinalQuota       int                 `json:"final_quota"`
	PromptTokens     int                 `json:"prompt_tokens"`
	CompletionTokens int                 `json:"completion_tokens"`
	Events           []BillingTraceEvent `json:"events"`
}

func GetBillingTraceSummary(traceId string) (*BillingTraceSummary, error) {
	traceId = strings.TrimSpace(traceId)
	if traceId == "" {
		return nil, errors.New("billing_trace_id is required")
	}

	logs, err := model.GetLogsByBillingTraceId(traceId)
	if err != nil {
		return nil, err
	}
	summary := &BillingTraceSummary{
		BillingTraceID: traceId,
		Status:         "NOT_FOUND",
		Events:         []BillingTraceEvent{},
	}
	for _, item := range logs {
		eventType := "OTHER"
		quotaDelta := 0
		switch item.Type {
		case model.LogTypeConsume:
			eventType = "CONSUME"
			quotaDelta = item.Quota
			summary.ConsumeQuota += item.Quota
			summary.PromptTokens += item.PromptTokens
			summary.CompletionTokens += item.CompletionTokens
		case model.LogTypeRefund:
			eventType = "REFUND"
			quotaDelta = -item.Quota
			summary.RefundQuota += item.Quota
		}
		summary.Events = append(summary.Events, BillingTraceEvent{
			ID:                item.Id,
			Type:              eventType,
			LogType:           item.Type,
			QuotaDelta:        quotaDelta,
			Quota:             item.Quota,
			RequestID:         item.RequestId,
			UpstreamRequestID: item.UpstreamRequestId,
			ModelName:         item.ModelName,
			TokenID:           item.TokenId,
			ChannelID:         item.ChannelId,
			Content:           item.Content,
			CreatedAt:         item.CreatedAt,
		})
	}
	summary.FinalQuota = summary.ConsumeQuota + summary.AdjustQuota - summary.RefundQuota
	if summary.FinalQuota < 0 {
		summary.FinalQuota = 0
	}
	// Refund receipts share a transaction with the funds and survive usage-log
	// cleanup or a failed log write. Keep the original DTO while reflecting that
	// authoritative refund in its existing totals.
	if receipt, err := model.GetBillingTraceRefund(traceId); err != nil {
		return nil, err
	} else if receipt != nil {
		return refundedTraceSummary(summary, receipt), nil
	}
	if len(logs) == 0 {
		return summary, nil
	}
	if summary.FinalQuota == 0 && summary.RefundQuota > 0 {
		summary.Status = "REFUNDED"
	} else {
		summary.Status = "SETTLED"
	}
	return summary, nil
}

// RefundBillingTrace preserves the legacy POST contract. Task rows and permanent
// primary-database refund receipts authorize refunds; usage logs alone never
// decide whether a previously completed refund can be paid again.
func RefundBillingTrace(ctx context.Context, traceId string, source string, operator string, reason string) (*BillingTraceSummary, bool, error) {
	traceId = strings.TrimSpace(traceId)
	if traceId == "" || len(traceId) > 128 {
		return nil, false, errors.New("billing_trace_id must contain 1 to 128 bytes")
	}
	if model.DB == nil || model.LOG_DB == nil {
		return nil, false, errors.New("billing database is unavailable")
	}
	summary, err := GetBillingTraceSummary(traceId)
	if err != nil {
		return nil, false, err
	}
	receipt, err := model.GetBillingTraceRefund(traceId)
	if err != nil {
		return nil, false, err
	}
	tasks := []model.Task{}
	if err := model.DB.Where("billing_trace_id = ?", traceId).Limit(2).Find(&tasks).Error; err != nil {
		return nil, false, err
	}
	if len(tasks) > 1 {
		return nil, false, model.ErrBillingRefundConflict
	}
	if len(tasks) == 1 {
		task := &tasks[0]
		if receipt == nil && task.Quota > 0 {
			// Relay inserts an immediately-terminal task before initial funding
			// settlement. A consume log is written only after that settlement;
			// use it as a readiness check, never as the refund amount authority.
			logs, err := model.GetLogsByBillingTraceId(traceId)
			if err != nil {
				return nil, false, err
			}
			ready := false
			for _, entry := range logs {
				if entry.Type != model.LogTypeConsume && entry.Type != model.LogTypeRefund {
					continue
				}
				if entry.UserId != task.UserId || entry.TokenId != task.PrivateData.TokenId || entry.ChannelId != task.ChannelId || entry.Quota < 0 {
					return nil, false, model.ErrBillingRefundConflict
				}
				ready = ready || (entry.Type == model.LogTypeConsume && entry.Quota > 0)
			}
			if !ready {
				return nil, false, model.ErrBillingRefundPending
			}
			// Preserve the legacy no-op when historical logs already prove a
			// full refund but the old non-atomic task quota writeback failed.
			if summary.FinalQuota <= 0 && summary.RefundQuota > 0 {
				return summary, false, nil
			}
		}
		// A zero-quota task is authoritative even when its refund log was lost.
		// It must never fall through to the old consume-log wallet fallback.
		amount, err := refundTraceTaskQuota(ctx, task, refundReason(source, operator, reason), true)
		if err != nil {
			return nil, false, err
		}
		updated, readErr := GetBillingTraceSummary(traceId)
		if readErr == nil {
			summary = updated
		}
		summary.FinalQuota = 0
		summary.RefundQuota = max(summary.RefundQuota, summary.ConsumeQuota+summary.AdjustQuota)
		if summary.RefundQuota > 0 {
			summary.Status = "REFUNDED"
		}
		return summary, amount > 0, nil
	}
	// Return an already committed receipt before consulting mutable usage logs.
	if receipt != nil {
		return refundedTraceSummary(summary, receipt), false, nil
	}
	logs, err := model.GetLogsByBillingTraceId(traceId)
	if err != nil {
		return nil, false, err
	}
	var base *model.Log
	var identity traceRefundFunding
	consume, refunded := 0, 0
	requiresReconciliation := false
	for _, entry := range logs {
		if entry.Type != model.LogTypeConsume && entry.Type != model.LogTypeRefund {
			continue
		}
		if entry.Quota < 0 || entry.Quota > common.MaxQuota {
			return nil, false, model.ErrBillingRefundConflict
		}
		if entry.Type == model.LogTypeRefund {
			if refunded > common.MaxQuota-entry.Quota {
				return nil, false, model.ErrBillingRefundConflict
			}
			refunded += entry.Quota
			continue
		}
		if entry.Quota == 0 {
			continue
		}
		var metadata struct {
			traceRefundFunding
			AdminInfo struct {
				RefundProtocol string `json:"billing_refund_protocol"`
			} `json:"admin_info"`
		}
		if err := common.UnmarshalJsonStr(entry.Other, &metadata); err != nil || metadata.IsTask || metadata.TaskId != "" {
			return nil, false, model.ErrBillingRefundConflict
		}
		funding := metadata.traceRefundFunding
		// Unknown/missing funding cannot safely be guessed to be a wallet.
		switch funding.BillingSource {
		case BillingSourceWallet:
			if funding.SubscriptionId != 0 || funding.ReservationId != "" || funding.ItemKey != "" {
				return nil, false, model.ErrBillingRefundConflict
			}
		case BillingSourceSubscription:
			if funding.SubscriptionId <= 0 || funding.ReservationId != "" || funding.ItemKey != "" {
				return nil, false, model.ErrBillingRefundConflict
			}
		case BillingSourceWorkflowReservation:
			if funding.SubscriptionId != 0 || funding.ReservationId == "" || funding.ItemKey == "" {
				return nil, false, model.ErrBillingRefundConflict
			}
		default:
			return nil, false, model.ErrBillingRefundConflict
		}
		// Historical wallet/subscription logs cannot prove that an earlier
		// refund log was not lost. Only consumption from the receipt protocol
		// may start a new log-only refund. Workflow items have authoritative
		// state in the primary database and are checked in the transaction.
		if funding.BillingSource != BillingSourceWorkflowReservation && metadata.AdminInfo.RefundProtocol != model.BillingRefundReceiptProtocol {
			requiresReconciliation = true
		}
		if base == nil {
			base, identity = entry, funding
		} else if base.UserId != entry.UserId || base.TokenId != entry.TokenId || base.ChannelId != entry.ChannelId || identity != funding {
			return nil, false, model.ErrBillingRefundConflict
		}
		if consume > common.MaxQuota-entry.Quota {
			return nil, false, model.ErrBillingRefundConflict
		}
		consume += entry.Quota
	}
	if base == nil || consume <= refunded {
		return summary, false, nil
	}
	if requiresReconciliation {
		return nil, false, model.ErrBillingRefundReconciliation
	}
	for _, entry := range logs {
		if entry.Type == model.LogTypeRefund && (entry.UserId != base.UserId || entry.TokenId != base.TokenId || entry.ChannelId != base.ChannelId) {
			return nil, false, model.ErrBillingRefundConflict
		}
	}
	receipt, changed, err := model.RefundLoggedBillingTrace(ctx, model.BillingTraceRefund{
		TraceId: traceId, UserId: base.UserId, TokenId: base.TokenId, ChannelId: base.ChannelId,
		BillingSource: identity.BillingSource, SubscriptionId: identity.SubscriptionId,
		ReservationId: identity.ReservationId, ItemKey: identity.ItemKey,
		ConsumeQuota: consume, PreviousRefundQuota: refunded, Quota: consume - refunded,
		Source: source, Operator: operator, Reason: reason,
	})
	if err != nil {
		return nil, false, err
	}
	if changed {
		other := model.NewLogOther()
		other.SetPublic("billing_trace_id", traceId)
		other.SetPublic("billing_source", identity.BillingSource)
		other.SetPublic("subscription_id", identity.SubscriptionId)
		other.SetPublic("workflow_quota_reservation_id", identity.ReservationId)
		other.SetPublic("workflow_quota_item_key", identity.ItemKey)
		other.SetAdmin("refund_source", source)
		other.SetAdmin("refund_operator", operator)
		other.SetAdmin("refund_reason", reason)
		other.SetAdmin("billing_refund_receipt_id", receipt.Id)
		model.RecordTaskBillingLog(model.RecordTaskBillingLogParams{
			UserId: base.UserId, TokenId: base.TokenId, ChannelId: base.ChannelId,
			LogType: model.LogTypeRefund, Quota: receipt.Quota, ModelName: base.ModelName,
			Group: base.Group, BillingTraceId: traceId, Other: other,
			Content: []*common.Message{common.NewMessage("Billing trace refund")},
		})
	}
	return refundedTraceSummary(summary, receipt), changed, nil
}

type traceRefundFunding struct {
	BillingSource  string `json:"billing_source"`
	SubscriptionId int    `json:"subscription_id"`
	ReservationId  string `json:"workflow_quota_reservation_id"`
	ItemKey        string `json:"workflow_quota_item_key"`
	IsTask         bool   `json:"is_task"`
	TaskId         string `json:"task_id"`
}

func refundedTraceSummary(summary *BillingTraceSummary, receipt *model.BillingTraceRefund) *BillingTraceSummary {
	if receipt.TaskId > 0 {
		summary.ConsumeQuota = max(summary.ConsumeQuota, receipt.ConsumeQuota)
		summary.RefundQuota = max(summary.RefundQuota, summary.ConsumeQuota+summary.AdjustQuota)
		summary.FinalQuota = 0
		summary.Status = "REFUNDED"
		return summary
	}
	summary.ConsumeQuota = max(summary.ConsumeQuota, receipt.ConsumeQuota)
	summary.RefundQuota = max(summary.RefundQuota, receipt.PreviousRefundQuota+receipt.Quota)
	summary.FinalQuota = max(summary.ConsumeQuota-summary.RefundQuota, 0)
	if summary.FinalQuota == 0 {
		summary.Status = "REFUNDED"
	}
	return summary
}

func refundReason(source, operator, reason string) string {
	return fmt.Sprintf("billing trace refund, source=%s, operator=%s, reason=%s", source, operator, reason)
}

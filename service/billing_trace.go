package service

import (
	"context"
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

// ErrBillingTraceRefundDisabled marks the retired manual trace refund API.
var ErrBillingTraceRefundDisabled = errors.New("Billing trace refunds are no longer supported.")

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

type BillingTraceWorkflowQuotaItem struct {
	ReservationID string `json:"reservation_id"`
	ItemKey       string `json:"item_key"`
	QuotedQuota   int    `json:"quoted_quota"`
}

type BillingTraceSummary struct {
	BillingTraceID    string                         `json:"billing_trace_id"`
	Status            string                         `json:"status"`
	ConsumeQuota      int                            `json:"consume_quota"`
	AdjustQuota       int                            `json:"adjust_quota"`
	RefundQuota       int                            `json:"refund_quota"`
	FinalQuota        int                            `json:"final_quota"`
	PromptTokens      int                            `json:"prompt_tokens"`
	CompletionTokens  int                            `json:"completion_tokens"`
	Events            []BillingTraceEvent            `json:"events"`
	WorkflowQuotaItem *BillingTraceWorkflowQuotaItem `json:"workflow_quota_item,omitempty"`
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
	var workflowItem *BillingTraceWorkflowQuotaItem
	var workflowUserID, workflowTokenID int
	workflowConflict := false
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
		if item.Type != model.LogTypeConsume && item.Type != model.LogTypeRefund {
			continue
		}
		var other struct {
			BillingSource string `json:"billing_source"`
			ReservationID string `json:"workflow_quota_reservation_id"`
			ItemKey       string `json:"workflow_quota_item_key"`
		}
		if err := common.UnmarshalJsonStr(item.Other, &other); err != nil || other.BillingSource != BillingSourceWorkflowReservation {
			continue
		}
		other.ReservationID = strings.TrimSpace(other.ReservationID)
		other.ItemKey = strings.TrimSpace(other.ItemKey)
		if other.ReservationID == "" || other.ItemKey == "" || item.UserId <= 0 || item.TokenId <= 0 {
			workflowConflict = true
			continue
		}
		if workflowItem == nil {
			workflowItem = &BillingTraceWorkflowQuotaItem{ReservationID: other.ReservationID, ItemKey: other.ItemKey}
			workflowUserID, workflowTokenID = item.UserId, item.TokenId
		} else if workflowItem.ReservationID != other.ReservationID || workflowItem.ItemKey != other.ItemKey || workflowUserID != item.UserId || workflowTokenID != item.TokenId {
			workflowConflict = true
		}
	}
	if workflowItem != nil && !workflowConflict {
		for _, item := range logs {
			if (item.Type == model.LogTypeConsume || item.Type == model.LogTypeRefund) && (item.UserId != workflowUserID || item.TokenId != workflowTokenID) {
				workflowConflict = true
				break
			}
		}
		// Metadata is best-effort: missing historical rows or a primary DB
		// lookup failure must not block the existing trace amounts/events.
		if !workflowConflict {
			quote, quoteErr := model.GetWorkflowQuotaItemQuote(workflowItem.ReservationID, workflowItem.ItemKey, workflowUserID, workflowTokenID)
			if quoteErr == nil && quote.QuotedQuota >= 0 && quote.QuotedQuota <= common.MaxQuota {
				workflowItem.QuotedQuota = quote.QuotedQuota
				summary.WorkflowQuotaItem = workflowItem
			}
		}
	}
	summary.FinalQuota = summary.ConsumeQuota + summary.AdjustQuota - summary.RefundQuota
	if summary.FinalQuota < 0 {
		summary.FinalQuota = 0
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

// RefundBillingTrace preserves the retired API contract without consulting
// accounting logs or changing funds. Automatic task refunds and reservation
// releases use their own accounting paths.
//
// Deprecated: trace logs cannot authorize a financial refund.
func RefundBillingTrace(ctx context.Context, traceId string, source string, operator string, reason string) (*BillingTraceSummary, bool, error) {
	return nil, false, ErrBillingTraceRefundDisabled
}

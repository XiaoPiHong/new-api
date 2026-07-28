package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
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
	BillingTraceID  string              `json:"billing_trace_id"`
	Status          string              `json:"status"`
	ConsumeQuota    int                 `json:"consume_quota"`
	AdjustQuota     int                 `json:"adjust_quota"`
	RefundQuota     int                 `json:"refund_quota"`
	FinalQuota      int                 `json:"final_quota"`
	PromptTokens    int                 `json:"prompt_tokens"`
	CompletionTokens int                `json:"completion_tokens"`
	Events          []BillingTraceEvent `json:"events"`
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

func RefundBillingTrace(ctx context.Context, traceId string, source string, operator string, reason string) (*BillingTraceSummary, bool, error) {
	traceId = strings.TrimSpace(traceId)
	if traceId == "" {
		return nil, false, errors.New("billing_trace_id is required")
	}
	summary, err := GetBillingTraceSummary(traceId)
	if err != nil {
		return nil, false, err
	}
	if summary.FinalQuota <= 0 {
		return summary, false, nil
	}

	if task, ok, err := model.GetByBillingTraceId(traceId); err != nil {
		return nil, false, err
	} else if ok && task.Quota > 0 {
		RefundTaskQuota(ctx, task, refundReason(source, operator, reason))
		updated, err := GetBillingTraceSummary(traceId)
		if err == nil && updated.FinalQuota > 0 {
			err = fmt.Errorf("billing trace %s refund not settled, final quota=%d", traceId, updated.FinalQuota)
		}
		return updated, true, err
	}

	if err := refundTraceFromConsumeLogs(ctx, traceId, summary.FinalQuota, source, operator, reason); err != nil {
		return nil, false, err
	}
	updated, err := GetBillingTraceSummary(traceId)
	if err == nil && updated.FinalQuota > 0 {
		err = fmt.Errorf("billing trace %s refund not settled, final quota=%d", traceId, updated.FinalQuota)
	}
	return updated, true, err
}

func refundTraceFromConsumeLogs(ctx context.Context, traceId string, quota int, source string, operator string, reason string) error {
	logs, err := model.GetLogsByBillingTraceId(traceId)
	if err != nil {
		return err
	}
	var base *model.Log
	for _, item := range logs {
		if item.Type == model.LogTypeConsume && item.Quota > 0 {
			base = item
			break
		}
	}
	if base == nil {
		return fmt.Errorf("billing trace %s has no consume log to refund", traceId)
	}
	if err := refundTraceFunding(base, quota); err != nil {
		return err
	}
	// 预占节点退款由预占主单统一返还 Token，这里只处理旧钱包/订阅 trace。
	_, _, workflowReservation := workflowReservationFromLog(base)
	if base.TokenId > 0 && !workflowReservation {
		if token, err := model.GetTokenById(base.TokenId); err == nil {
			if err := model.IncreaseTokenQuota(base.TokenId, token.Key, quota); err != nil {
				logger.LogWarn(ctx, fmt.Sprintf("退还 trace %s 令牌额度失败: %s", traceId, err.Error()))
			}
		} else {
			logger.LogWarn(ctx, fmt.Sprintf("退还 trace %s 获取 token 失败: %s", traceId, err.Error()))
		}
	}

	other := map[string]interface{}{
		"billing_trace_id": traceId,
		"source":           source,
		"operator":         operator,
		"reason":           reason,
	}
	model.RecordTaskBillingLog(model.RecordTaskBillingLogParams{
		UserId:         base.UserId,
		LogType:        model.LogTypeRefund,
		Content:        refundReason(source, operator, reason),
		ChannelId:      base.ChannelId,
		ModelName:      base.ModelName,
		Quota:          quota,
		TokenId:        base.TokenId,
		Group:          base.Group,
		BillingTraceId: traceId,
		Other:          other,
	})
	return nil
}

func refundTraceFunding(base *model.Log, quota int) error {
	other, _ := common.StrToMap(base.Other)
	if reservationId, itemKey, ok := workflowReservationFromLog(base); ok {
		// 预占请求的 trace 退款只把节点消费归零，不能直接给钱包和 Token 加额度。
		return model.AdjustSettledWorkflowQuotaItem(reservationId, itemKey, base.UserId, base.TokenId, 0)
	}
	if other != nil && fmt.Sprint(other["billing_source"]) == BillingSourceSubscription {
		subscriptionId := intFromAny(other["subscription_id"])
		if subscriptionId > 0 {
			return model.PostConsumeUserSubscriptionDelta(subscriptionId, int64(-quota))
		}
	}
	return model.IncreaseUserQuota(base.UserId, quota, false)
}

// workflowReservationFromLog 从消费日志恢复预占身份，供人工 trace 退款安全复用。
func workflowReservationFromLog(base *model.Log) (string, string, bool) {
	if base == nil {
		return "", "", false
	}
	other, _ := common.StrToMap(base.Other)
	if other == nil || fmt.Sprint(other["billing_source"]) != BillingSourceWorkflowReservation {
		return "", "", false
	}
	reservationId := strings.TrimSpace(fmt.Sprint(other["workflow_quota_reservation_id"]))
	itemKey := strings.TrimSpace(fmt.Sprint(other["workflow_quota_item_key"]))
	return reservationId, itemKey, reservationId != "" && itemKey != ""
}

func intFromAny(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		var n int
		_, _ = fmt.Sscanf(v, "%d", &n)
		return n
	default:
		return 0
	}
}

func refundReason(source string, operator string, reason string) string {
	parts := []string{"billing trace refund"}
	if source != "" {
		parts = append(parts, "source="+source)
	}
	if operator != "" {
		parts = append(parts, "operator="+operator)
	}
	if reason != "" {
		parts = append(parts, "reason="+reason)
	}
	return strings.Join(parts, ", ")
}

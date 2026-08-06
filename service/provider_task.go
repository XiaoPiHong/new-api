package service

import (
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/model"
)

// ProviderTaskSnapshot 是供业务服务 reconcile 使用的最小异步任务快照。
// 不返回 PrivateData、Key 或完整请求体，避免内部查询接口泄露用户凭证。
type ProviderTaskSnapshot struct {
	BillingTraceID string `json:"billing_trace_id"`
	ProviderTaskID string `json:"provider_task_id"`
	UpstreamTaskID string `json:"upstream_task_id,omitempty"`
	Platform       string `json:"platform"`
	Status         string `json:"status"`
	Progress       string `json:"progress,omitempty"`
	ResultURL      string `json:"result_url,omitempty"`
	FailReason     string `json:"fail_reason,omitempty"`
	SubmitTime     int64  `json:"submit_time"`
	StartTime      int64  `json:"start_time"`
	FinishTime     int64  `json:"finish_time"`
	UpdatedAt      int64  `json:"updated_at"`
}

func GetProviderTaskByBillingTrace(
	traceID string,
) (*ProviderTaskSnapshot, bool, error) {
	traceID = strings.TrimSpace(traceID)
	if traceID == "" {
		return nil, false, errors.New("billing_trace_id is required")
	}

	task, exists, err := model.GetByBillingTraceId(traceID)
	if err != nil || !exists {
		return nil, exists, err
	}

	return &ProviderTaskSnapshot{
		BillingTraceID: traceID,
		ProviderTaskID: task.TaskID,
		UpstreamTaskID: task.GetUpstreamTaskID(),
		Platform:       string(task.Platform),
		Status:         string(task.Status),
		Progress:       task.Progress,
		ResultURL:      task.GetResultURL(),
		FailReason:     task.FailReason,
		SubmitTime:     task.SubmitTime,
		StartTime:      task.StartTime,
		FinishTime:     task.FinishTime,
		UpdatedAt:      task.UpdatedAt,
	}, true, nil
}

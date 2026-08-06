package controller

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

// GetProviderTaskByBillingTrace 只供 SSO 等受信任业务服务查询异步任务。
func GetProviderTaskByBillingTrace(c *gin.Context) {
	if !workflowQuotaRequestAuthorized(c) {
		workflowQuotaError(
			c,
			403,
			"access_denied",
			"provider task 查询鉴权失败",
		)
		return
	}

	traceID := strings.TrimSpace(c.Param("trace_id"))
	task, exists, err := service.GetProviderTaskByBillingTrace(traceID)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if !exists {
		common.ApiSuccess(c, gin.H{
			"status":           "NOT_FOUND",
			"billing_trace_id": traceID,
		})
		return
	}
	common.ApiSuccess(c, gin.H{
		"status": "FOUND",
		"task":   task,
	})
}

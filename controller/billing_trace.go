package controller

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

type refundBillingTraceRequest struct {
	Source   string `json:"source"`
	Operator string `json:"operator"`
	Reason   string `json:"reason"`
}

func GetBillingTrace(c *gin.Context) {
	traceId := strings.TrimSpace(c.Param("trace_id"))
	summary, err := service.GetBillingTraceSummary(traceId)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, summary)
}

func RefundBillingTrace(c *gin.Context) {
	traceId := strings.TrimSpace(c.Param("trace_id"))
	var req refundBillingTraceRequest
	_ = c.ShouldBindJSON(&req)
	if req.Source == "" {
		req.Source = "MANUAL"
	}
	if req.Operator == "" {
		req.Operator = c.GetString("username")
	}
	summary, refunded, err := service.RefundBillingTrace(c.Request.Context(), traceId, req.Source, req.Operator, req.Reason)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, gin.H{
		"refunded": refunded,
		"trace":    summary,
	})
}

package controller

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

func GetBillingTrace(c *gin.Context) {
	traceId := strings.TrimSpace(c.Param("trace_id"))
	summary, err := service.GetBillingTraceSummary(traceId)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, summary)
}

type refundBillingTraceRequest struct {
	Source   string `json:"source"`
	Operator string `json:"operator"`
	Reason   string `json:"reason"`
}

func RefundBillingTrace(c *gin.Context) {
	var req refundBillingTraceRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		common.ApiErrorStatus(c, http.StatusBadRequest, common.NewMessage("invalid refund request"))
		return
	}
	if req.Source == "" {
		req.Source = "MANUAL"
	}
	// The authenticated operator takes precedence over a client-supplied label.
	if operator := c.GetString("username"); operator != "" {
		req.Operator = operator
	}
	if len(req.Source) > 64 || len(req.Operator) > 128 || len(req.Reason) > 2000 {
		common.ApiErrorStatus(c, http.StatusBadRequest, common.NewMessage("refund metadata is too long"))
		return
	}
	summary, refunded, err := service.RefundBillingTrace(c.Request.Context(), c.Param("trace_id"), req.Source, req.Operator, req.Reason)
	if err != nil {
		common.ApiError(c, common.NewMessage(err.Error()))
		return
	}
	common.ApiSuccess(c, gin.H{"refunded": refunded, "trace": summary})
}

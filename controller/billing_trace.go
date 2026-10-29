package controller

import (
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

// RefundBillingTrace keeps the authenticated legacy route as a retirement response.
func RefundBillingTrace(c *gin.Context) {
	common.ApiErrorStatus(c, http.StatusGone,
		common.NewMessage(service.ErrBillingTraceRefundDisabled.Error()),
		gin.H{"code": "BILLING_TRACE_REFUND_DISABLED"})
}

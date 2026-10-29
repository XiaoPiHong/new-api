package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefundBillingTraceReturnsGoneWithoutDatabaseAccess(t *testing.T) {
	originalDB, originalLogDB := model.DB, model.LOG_DB
	model.DB, model.LOG_DB = nil, nil
	t.Cleanup(func() { model.DB, model.LOG_DB = originalDB, originalLogDB })
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.POST("/api/billing/traces/:trace_id/refund", RefundBillingTrace)
	for _, body := range []string{`{"source":"MANUAL","reason":"failed"}`, `{invalid-json`} {
		request := httptest.NewRequest(http.MethodPost, "/api/billing/traces/legacy-trace/refund", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		engine.ServeHTTP(response, request)
		require.Equal(t, http.StatusGone, response.Code)
		var result struct {
			Success    bool   `json:"success"`
			Code       string `json:"code"`
			Message    string `json:"message"`
			MessageKey string `json:"message_key"`
		}
		require.NoError(t, common.Unmarshal(response.Body.Bytes(), &result))
		assert.False(t, result.Success)
		assert.Equal(t, "BILLING_TRACE_REFUND_DISABLED", result.Code)
		assert.Equal(t, "Billing trace refunds are no longer supported.", result.Message)
		assert.Equal(t, result.Message, result.MessageKey)
	}
}

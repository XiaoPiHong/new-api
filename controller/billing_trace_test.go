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

func TestRefundBillingTraceValidatesRequestAndReachesRestoredService(t *testing.T) {
	originalDB, originalLogDB := model.DB, model.LOG_DB
	model.DB, model.LOG_DB = nil, nil
	t.Cleanup(func() { model.DB, model.LOG_DB = originalDB, originalLogDB })
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.POST("/api/billing/traces/:trace_id/refund", RefundBillingTrace)
	for _, testCase := range []struct {
		body    string
		status  int
		message string
	}{
		{`{"source":"MANUAL","reason":"failed"}`, http.StatusOK, "billing database is unavailable"},
		{``, http.StatusOK, "billing database is unavailable"},
		{`{invalid-json`, http.StatusBadRequest, "invalid refund request"},
	} {
		request := httptest.NewRequest(http.MethodPost, "/api/billing/traces/legacy-trace/refund", strings.NewReader(testCase.body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		engine.ServeHTTP(response, request)
		require.Equal(t, testCase.status, response.Code)
		var result struct {
			Success    bool   `json:"success"`
			Code       string `json:"code"`
			Message    string `json:"message"`
			MessageKey string `json:"message_key"`
		}
		require.NoError(t, common.Unmarshal(response.Body.Bytes(), &result))
		assert.False(t, result.Success)
		assert.Empty(t, result.Code)
		assert.Equal(t, testCase.message, result.Message)
		assert.Equal(t, testCase.message, result.MessageKey)
	}
}

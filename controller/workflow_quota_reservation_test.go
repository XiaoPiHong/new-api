package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/constant"
	taskrunninghub "github.com/QuantumNous/new-api/relay/channel/task/runninghub"
	tasksora "github.com/QuantumNous/new-api/relay/channel/task/sora"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestWorkflowVideoQuoteMultiplierRunningHubUsesPerCallAdaptor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	multiplier, ok := workflowVideoQuoteMultiplierFromAdaptor(
		ctx,
		workflowQuotaQuoteItem{
			Model:      "runninghub-rhart-video-g",
			Seconds:    3,
			Resolution: "720p",
			N:          1,
		},
		&relaycommon.RelayInfo{},
		constant.ChannelTypeRunningHub,
		&taskrunninghub.TaskAdaptor{},
	)

	require.True(t, ok)
	require.Equal(t, 1.0, multiplier)
}

func TestWorkflowVideoQuoteMultiplierKeepsAdaptorRatios(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	multiplier, ok := workflowVideoQuoteMultiplierFromAdaptor(
		ctx,
		workflowQuotaQuoteItem{
			Model:   "sora-2",
			Seconds: 3,
			N:       2,
		},
		&relaycommon.RelayInfo{},
		constant.ChannelTypeSora,
		&tasksora.TaskAdaptor{},
	)

	require.True(t, ok)
	require.Equal(t, 6.0, multiplier)
}

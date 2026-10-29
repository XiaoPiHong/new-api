package controller

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	hostdto "github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	taskkievideo "github.com/QuantumNous/new-api/relay/channel/task/kievideo"
	taskrunninghub "github.com/QuantumNous/new-api/relay/channel/task/runninghub"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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

func TestWorkflowVideoQuoteMultiplierKieVideoUsesPerCallAdaptor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	multiplier, ok := workflowVideoQuoteMultiplierFromAdaptor(
		ctx,
		workflowQuotaQuoteItem{
			Model:      "grok-imagine-video-1-5-preview",
			Seconds:    8,
			Resolution: "480p",
			N:          1,
		},
		&relaycommon.RelayInfo{},
		constant.ChannelTypeKieVideo,
		&taskkievideo.TaskAdaptor{},
	)

	require.True(t, ok)
	require.Equal(t, 1.0, multiplier)
}

func TestWorkflowVideoQuoteMultiplierSoraCompatibleUsesPerCallAdaptor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name        string
		channelType int
		model       string
	}{
		{"Sora", constant.ChannelTypeSora, "aihubcc-grok-imagine-video-1.5-preview"},
		{"OpenAI video", constant.ChannelTypeOpenAI, "yooart-Vidu-q2-turbo-720p"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			item := workflowQuotaQuoteItem{
				Model:      testCase.model,
				Seconds:    6,
				Resolution: "1080p",
				N:          1,
			}
			multiplier, ok := workflowVideoQuoteMultiplierFromAdaptor(
				ctx, item, &relaycommon.RelayInfo{}, testCase.channelType,
				workflowQuoteBillingAdaptor{},
			)
			require.True(t, ok)
			require.Equal(t, 1.0, multiplier)

			multiplier, ok = workflowVideoQuoteMultiplierFromAdaptor(
				ctx, item, &relaycommon.RelayInfo{}, testCase.channelType,
				workflowQuoteBillingAdaptor{ratios: map[string]float64{"seconds": 3}},
			)
			require.True(t, ok)
			require.Equal(t, 3.0, multiplier)
		})
	}
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
		workflowQuoteBillingAdaptor{ratios: map[string]float64{"seconds": 3}},
	)

	require.True(t, ok)
	require.Equal(t, 6.0, multiplier)
}

type workflowQuoteBillingAdaptor struct {
	ratios map[string]float64
}

func (a workflowQuoteBillingAdaptor) EstimateBilling(*gin.Context, *relaycommon.RelayInfo) map[string]float64 {
	return a.ratios
}

func TestWorkflowQuoteRejectsUnboundedBillingInputs(t *testing.T) {
	for _, item := range []workflowQuotaQuoteItem{
		{MaxTokens: -1},
		{MaxTokens: math.MaxInt32},
		{N: -1},
		{N: dto.MaxImageN + 1},
		{Seconds: -1},
		{Seconds: relaycommon.MaxTaskDurationSeconds + 1},
	} {
		_, err := normalizeWorkflowQuotaQuoteItem(item)
		require.Error(t, err)
	}
	item, err := normalizeWorkflowQuotaQuoteItem(workflowQuotaQuoteItem{
		Adapter: "newapi.video", N: dto.MaxImageN, Seconds: relaycommon.MaxTaskDurationSeconds,
	})
	require.NoError(t, err)
	require.Equal(t, dto.MaxImageN, item.N)
}

func TestWorkflowQuotaBufferRejectsOverflow(t *testing.T) {
	t.Setenv("WORKFLOW_QUOTA_RESERVE_BUFFER_RATIO", "1.1")
	_, err := reserveQuotaWithBuffer(common.MaxQuota)
	require.Error(t, err)

	quota, err := reserveQuotaWithBuffer(10)
	require.NoError(t, err)
	require.Equal(t, 11, quota)
}

func TestWorkflowQuoteSelectionRetainsTokenPinAndFiltersPerNode(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	previousDB := model.DB
	previousRedis, previousMemory := common.RedisEnabled, common.MemoryCacheEnabled
	model.DB = db
	common.RedisEnabled, common.MemoryCacheEnabled = false, false
	t.Cleanup(func() {
		model.DB = previousDB
		common.RedisEnabled, common.MemoryCacheEnabled = previousRedis, previousMemory
		require.NoError(t, sqlDB.Close())
	})
	require.NoError(t, db.AutoMigrate(&model.Channel{}, &model.Ability{}))
	channel := model.Channel{
		Id: 1, Name: "image-only", Type: constant.ChannelTypeAdvancedCustom, Key: "test-only",
		Status: common.ChannelStatusEnabled, Group: "default", Models: "workflow-model",
	}
	channel.SetOtherSettings(dto.ChannelOtherSettings{AdvancedCustom: &dto.AdvancedCustomConfig{
		Routes: []dto.AdvancedCustomRoute{{IncomingPath: "/v1/images/generations", UpstreamPath: "/images"}},
	}})
	require.NoError(t, db.Create(&channel).Error)
	require.NoError(t, channel.AddAbilities(db))
	fallback := model.Channel{
		Id: 2, Name: "fallback", Type: constant.ChannelTypeOpenAI, Key: "test-only",
		Status: common.ChannelStatusEnabled, Group: "default", Models: "workflow-model",
	}
	require.NoError(t, db.Create(&fallback).Error)
	require.NoError(t, fallback.AddAbilities(db))
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/api/workflow-quota/reservations", nil)
	common.SetContextKey(c, constant.ContextKeyUsingGroup, "default")
	common.SetContextKey(c, constant.ContextKeyUserGroup, "default")
	constraints := service.GetChannelConstraints(c)
	constraints.AddPin(hostdto.ChannelPin{ChannelId: 1, Source: hostdto.PinSourceToken, Rank: hostdto.PinRankToken})
	_, err = workflowQuotaFindAvailableChannel(c, workflowQuotaQuoteItem{Adapter: "newapi.chat", Model: "workflow-model"})
	require.ErrorIs(t, err, errWorkflowQuotaModelUnavailable, "must not fall back outside the token's pin")
	selected, err := workflowQuotaFindAvailableChannel(c, workflowQuotaQuoteItem{Adapter: "newapi.images", Model: "workflow-model"})
	require.NoError(t, err)
	require.NotNil(t, selected)
	require.Equal(t, 1, selected.Id)
	require.Empty(t, constraints.Filters, "one node's request path must not constrain later nodes")
	require.Len(t, constraints.Pins, 1)
	constraints.AddFilter(hostdto.ChannelFilter{
		Kind: hostdto.FilterTaskPluginIdentity, TaskPluginKey: "allowed-plugin",
		TaskPluginChannelTypes: []int{constant.ChannelTypeOpenAI},
	})
	_, err = workflowQuotaFindAvailableChannel(c, workflowQuotaQuoteItem{Adapter: "newapi.images", Model: "workflow-model"})
	require.ErrorIs(t, err, errWorkflowQuotaModelUnavailable, "the existing identity filter must remain enforced")
	require.Len(t, constraints.Filters, 1)
}

func TestWorkflowQuoteTokenModelRulesMatchRelay(t *testing.T) {
	for _, tc := range []struct {
		allowed string
		model   string
		ok      bool
	}{
		{"claude-3-7-sonnet", "claude-3-7-sonnet@thinking:on", true},
		{"claude-3-7-sonnet", "claude-3-7-sonnet-thinking", true},
		{"gemini-2.5-flash-thinking-*", "gemini-2.5-flash-thinking-8192", true},
		{"claude-3-7-sonnet-thinking", "claude-3-7-sonnet", false},
	} {
		t.Run(tc.model, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			common.SetContextKey(c, constant.ContextKeyTokenModelLimitEnabled, true)
			common.SetContextKey(c, constant.ContextKeyTokenModelLimit, map[string]bool{tc.allowed: true})
			err := validateWorkflowQuotaTokenModel(c, tc.model)
			if tc.ok {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, errWorkflowQuotaModelForbidden)
			}
		})
	}
}

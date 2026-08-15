package controller

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/pkg/billingexpr"
	"github.com/QuantumNous/new-api/relay"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relay/helper"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
)

const (
	// 预占有效期由 SSO 传入，但必须限制范围，避免额度永久占用或过早回收。
	workflowQuotaDefaultExpiresSeconds = 2 * 60 * 60
	workflowQuotaMinExpiresSeconds     = 5 * 60
	workflowQuotaMaxExpiresSeconds     = 24 * 60 * 60
	workflowQuotaMaxItems              = 64
)

var (
	// 报价前置校验使用哨兵错误区分“Token 无权使用”和“当前没有可用渠道”。
	errWorkflowQuotaModelForbidden   = errors.New("workflow quota model forbidden")
	errWorkflowQuotaModelUnavailable = errors.New("workflow quota model unavailable")
)

// workflowQuotaReservationRequest 使用任务 ID 派生幂等键，并一次提交全部模型节点。
type workflowQuotaReservationRequest struct {
	IdempotencyKey string                   `json:"idempotency_key" binding:"required"`
	ExpiresIn      int                      `json:"expires_in_seconds"`
	Items          []workflowQuotaQuoteItem `json:"items" binding:"required"`
}

// workflowQuotaQuoteItem 是跨项目的稳定计价描述，不接收客户端直接传入 quota。
// NewAPI 根据自己的模型价格、分组倍率和请求参数生成最终预占额度。
type workflowQuotaQuoteItem struct {
	ItemKey     string `json:"item_key"`
	Adapter     string `json:"adapter"`
	Model       string `json:"model"`
	Prompt      string `json:"prompt,omitempty"`
	MaxTokens   int    `json:"max_tokens,omitempty"`
	Size        string `json:"size,omitempty"`
	Quality     string `json:"quality,omitempty"`
	N           int    `json:"n,omitempty"`
	Seconds     int    `json:"seconds,omitempty"`
	FPS         int    `json:"fps,omitempty"`
	Resolution  string `json:"resolution,omitempty"`
	AspectRatio string `json:"aspect_ratio,omitempty"`
}

type releaseWorkflowQuotaReservationRequest struct {
	// 失败任务为 true，沿用 XphAI 失败任务整单退款策略；成功任务只释放未消费额度。
	RefundConsumed bool `json:"refund_consumed"`
}

// CreateWorkflowQuotaReservation 在任何模型节点执行前完成整条工作流报价和额度预占。
func CreateWorkflowQuotaReservation(c *gin.Context) {
	if !workflowQuotaRequestAuthorized(c) {
		workflowQuotaError(c, http.StatusForbidden, "access_denied", "工作流额度预占鉴权失败")
		return
	}

	var req workflowQuotaReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		workflowQuotaError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if len(req.Items) == 0 || len(req.Items) > workflowQuotaMaxItems {
		workflowQuotaError(c, http.StatusBadRequest, "invalid_request", fmt.Sprintf("items 数量必须在 1 到 %d 之间", workflowQuotaMaxItems))
		return
	}

	quotedItems := make([]model.WorkflowQuotaReservationItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		quota, err := quoteWorkflowQuotaItem(c, item)
		if err != nil {
			switch {
			case errors.Is(err, errWorkflowQuotaModelForbidden):
				workflowQuotaError(c, http.StatusForbidden, "access_denied", err.Error())
			case errors.Is(err, errWorkflowQuotaModelUnavailable):
				workflowQuotaError(c, http.StatusServiceUnavailable, "model_not_found", err.Error())
			default:
				workflowQuotaError(c, http.StatusBadRequest, "model_price_error", err.Error())
			}
			return
		}
		quotedItems = append(quotedItems, model.WorkflowQuotaReservationItemInput{
			ItemKey:     strings.TrimSpace(item.ItemKey),
			Adapter:     strings.TrimSpace(item.Adapter),
			ModelName:   strings.TrimSpace(item.Model),
			QuotedQuota: quota,
		})
	}

	expiresIn := req.ExpiresIn
	switch {
	case expiresIn <= 0:
		expiresIn = workflowQuotaDefaultExpiresSeconds
	case expiresIn < workflowQuotaMinExpiresSeconds:
		expiresIn = workflowQuotaMinExpiresSeconds
	case expiresIn > workflowQuotaMaxExpiresSeconds:
		expiresIn = workflowQuotaMaxExpiresSeconds
	}

	reservation, created, err := model.ReserveWorkflowQuota(model.ReserveWorkflowQuotaRequest{
		ReservationKey: strings.TrimSpace(req.IdempotencyKey),
		UserId:         c.GetInt("id"),
		TokenId:        c.GetInt("token_id"),
		ExpiresAt:      common.GetTimestamp() + int64(expiresIn),
		Items:          quotedItems,
	})
	if err != nil {
		writeWorkflowQuotaModelError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"reservation_id": reservation.PublicId,
			"status":         reservation.Status,
			"reserved_quota": reservation.ReservedQuota,
			"consumed_quota": reservation.ConsumedQuota,
			"released_quota": reservation.ReleasedQuota,
			"expires_at":     reservation.ExpiresAt,
			"created":        created,
			"items":          reservation.Items,
		},
	})
}

// ReleaseWorkflowQuotaReservationInternal 不依赖用户 Token，只允许持有共享密钥的
// SSO 调用。用户即使在任务期间删除 Token，失败任务仍然可以完成整单退款。
func ReleaseWorkflowQuotaReservationInternal(c *gin.Context) {
	if !workflowQuotaRequestAuthorized(c) {
		workflowQuotaError(c, http.StatusForbidden, "access_denied", "工作流额度预占鉴权失败")
		return
	}

	var req releaseWorkflowQuotaReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		workflowQuotaError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	reservation, changed, err := model.ReleaseWorkflowQuotaReservationInternal(
		strings.TrimSpace(c.Param("reservation_id")),
		req.RefundConsumed,
	)
	if err != nil {
		writeWorkflowQuotaModelError(c, err)
		return
	}
	writeWorkflowQuotaReleaseResponse(c, reservation, changed)
}

// writeWorkflowQuotaReleaseResponse 统一成功和失败释放的响应结构，changed 表示本次是否实际退款。
func writeWorkflowQuotaReleaseResponse(c *gin.Context, reservation *model.WorkflowQuotaReservation, changed bool) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"reservation_id": reservation.PublicId,
			"status":         reservation.Status,
			"reserved_quota": reservation.ReservedQuota,
			"consumed_quota": reservation.ConsumedQuota,
			"released_quota": reservation.ReleasedQuota,
			"changed":        changed,
			"items":          reservation.Items,
		},
	})
}

// quoteWorkflowQuotaItem 复用正式请求的价格计算器，按 adapter 补齐协议默认值。
func quoteWorkflowQuotaItem(c *gin.Context, item workflowQuotaQuoteItem) (int, error) {
	item.ItemKey = strings.TrimSpace(item.ItemKey)
	item.Adapter = strings.TrimSpace(item.Adapter)
	item.Model = strings.TrimSpace(item.Model)
	if item.ItemKey == "" || item.Model == "" {
		return 0, errors.New("工作流计费节点缺少 item_key 或 model")
	}
	if err := validateWorkflowQuotaTokenModel(c, item.Model); err != nil {
		return 0, err
	}
	if err := validateWorkflowQuotaModelAvailability(c, item.Model); err != nil {
		return 0, err
	}

	info, err := workflowQuotaRelayInfo(c, item)
	if err != nil {
		return 0, err
	}

	switch item.Adapter {
	case "newapi.chat":
		maxTokens := item.MaxTokens
		if maxTokens <= 0 {
			maxTokens = 4096
		}
		meta := &types.TokenCountMeta{
			CombineText: item.Prompt,
			MaxTokens:   maxTokens,
		}
		promptTokens, err := service.EstimateRequestToken(c, meta, info)
		if err != nil {
			return 0, err
		}
		price, err := helper.ModelPriceHelper(c, info, promptTokens, meta)
		if err != nil {
			return 0, err
		}
		return reserveQuotaWithBuffer(price.QuotaToPreConsume), nil
	case "newapi.images":
		n := item.N
		if n <= 0 {
			n = 1
		}
		imageN := uint(n)
		request := &dto.ImageRequest{
			Model:   item.Model,
			Prompt:  item.Prompt,
			N:       &imageN,
			Size:    item.Size,
			Quality: item.Quality,
		}
		meta := request.GetTokenCountMeta()
		promptTokens, err := service.EstimateRequestToken(c, meta, info)
		if err != nil {
			return 0, err
		}
		price, err := helper.ModelPriceHelper(c, info, promptTokens, meta)
		if err != nil {
			return 0, err
		}
		// 图片数量在真实响应结算阶段才乘入，这里必须提前计入整单报价。
		return reserveQuotaWithBuffer(price.QuotaToPreConsume * n), nil
	case "newapi.video":
		price, err := helper.ModelPriceHelperPerCall(c, info)
		if err != nil {
			return 0, err
		}
		if price.FreeModel {
			return 0, nil
		}
		seconds := item.Seconds
		if seconds <= 0 {
			seconds = 5
		}
		n := item.N
		if n <= 0 {
			n = 1
		}
		multiplier := workflowVideoQuoteMultiplier(c, item, info, seconds, n)
		return reserveQuotaWithBuffer(int(math.Ceil(float64(price.Quota) * multiplier))), nil
	default:
		return 0, fmt.Errorf("不支持的工作流计费 adapter: %s", item.Adapter)
	}
}

// validateWorkflowQuotaModelAvailability 提前确认每个节点至少存在一个可用渠道。
// 该检查不能保证渠道在真实执行时永不变化，但能拦住发布时就不可用的后续模型。
func validateWorkflowQuotaModelAvailability(c *gin.Context, modelName string) error {
	channel, err := workflowQuotaFindAvailableChannel(c, modelName)
	if err != nil {
		return err
	}
	if channel != nil {
		return nil
	}
	usingGroup := common.GetContextKeyString(c, constant.ContextKeyUsingGroup)
	return fmt.Errorf("%w: 分组 %s 没有模型 %s 的可用渠道", errWorkflowQuotaModelUnavailable, usingGroup, modelName)
}

func workflowQuotaFindAvailableChannel(c *gin.Context, modelName string) (*model.Channel, error) {
	usingGroup := common.GetContextKeyString(c, constant.ContextKeyUsingGroup)
	groups := []string{usingGroup}
	if usingGroup == "auto" {
		userGroup := common.GetContextKeyString(c, constant.ContextKeyUserGroup)
		groups = service.GetUserAutoGroup(userGroup)
	}
	for _, group := range groups {
		channel, err := model.GetRandomSatisfiedChannel(group, modelName, 0)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", errWorkflowQuotaModelUnavailable, err.Error())
		}
		if channel != nil {
			return channel, nil
		}
	}
	return nil, nil
}

// validateWorkflowQuotaTokenModel 在整单预占阶段提前执行 Token 模型白名单检查。
// 否则第二个节点无权限时，前面节点即使额度足够也会产生无效资源消耗。
func validateWorkflowQuotaTokenModel(c *gin.Context, modelName string) error {
	if !common.GetContextKeyBool(c, constant.ContextKeyTokenModelLimitEnabled) {
		return nil
	}
	value, ok := common.GetContextKey(c, constant.ContextKeyTokenModelLimit)
	if !ok {
		return fmt.Errorf("%w: API Token 未授权任何模型", errWorkflowQuotaModelForbidden)
	}
	limits, ok := value.(map[string]bool)
	if !ok {
		return fmt.Errorf("%w: API Token 模型权限配置格式错误", errWorkflowQuotaModelForbidden)
	}
	matchName := ratio_setting.FormatMatchingModelName(modelName)
	if !limits[matchName] {
		return fmt.Errorf("%w: API Token 无权使用模型 %s", errWorkflowQuotaModelForbidden, modelName)
	}
	return nil
}

// workflowQuotaRelayInfo 构造只用于报价的 RelayInfo，不会发起上游请求或产生消费日志。
func workflowQuotaRelayInfo(c *gin.Context, item workflowQuotaQuoteItem) (*relaycommon.RelayInfo, error) {
	requestBody, err := common.Marshal(map[string]interface{}{
		"model":        item.Model,
		"prompt":       item.Prompt,
		"max_tokens":   item.MaxTokens,
		"size":         item.Size,
		"quality":      item.Quality,
		"n":            item.N,
		"seconds":      item.Seconds,
		"fps":          item.FPS,
		"resolution":   item.Resolution,
		"aspect_ratio": item.AspectRatio,
	})
	if err != nil {
		return nil, err
	}
	userSetting, _ := common.GetContextKeyType[dto.UserSetting](c, constant.ContextKeyUserSetting)
	info := &relaycommon.RelayInfo{
		UserId:          c.GetInt("id"),
		TokenId:         c.GetInt("token_id"),
		TokenKey:        c.GetString("token_key"),
		TokenUnlimited:  c.GetBool("token_unlimited_quota"),
		UsingGroup:      common.GetContextKeyString(c, constant.ContextKeyUsingGroup),
		UserGroup:       common.GetContextKeyString(c, constant.ContextKeyUserGroup),
		OriginModelName: item.Model,
		UserSetting:     userSetting,
		RelayFormat:     types.RelayFormatOpenAI,
		BillingRequestInput: &billingexpr.RequestInput{
			Body: requestBody,
		},
		TaskRelayInfo: &relaycommon.TaskRelayInfo{},
	}
	return info, nil
}

func workflowVideoQuoteMultiplier(c *gin.Context, item workflowQuotaQuoteItem, info *relaycommon.RelayInfo, seconds int, n int) float64 {
	if multiplier, ok := workflowVideoQuoteMultiplierFromChannel(c, item, info); ok {
		return multiplier
	}
	// 不同视频渠道对分辨率倍率定义不同。无法定位到具体 task adaptor 时，保留旧的保守上界。
	return float64(seconds*n) * workflowVideoResolutionUpperRatio(item.Resolution)
}

func workflowVideoQuoteMultiplierFromChannel(c *gin.Context, item workflowQuotaQuoteItem, info *relaycommon.RelayInfo) (float64, bool) {
	channel, err := workflowQuotaFindAvailableChannel(c, item.Model)
	if err != nil || channel == nil {
		return 0, false
	}
	adaptor := relay.GetTaskAdaptor(constant.TaskPlatform(strconv.Itoa(channel.Type)))
	return workflowVideoQuoteMultiplierFromAdaptor(c, item, info, channel.Type, adaptor)
}

func workflowVideoQuoteMultiplierFromAdaptor(c *gin.Context, item workflowQuotaQuoteItem, info *relaycommon.RelayInfo, channelType int, adaptor interface {
	EstimateBilling(*gin.Context, *relaycommon.RelayInfo) map[string]float64
}) (float64, bool) {
	if adaptor == nil {
		return 0, false
	}
	if info == nil {
		info = &relaycommon.RelayInfo{}
	}
	if info.TaskRelayInfo == nil {
		info.TaskRelayInfo = &relaycommon.TaskRelayInfo{}
	}

	taskReq := workflowVideoTaskSubmitReq(item)
	previousTaskReq, hadPreviousTaskReq := c.Get("task_request")
	c.Set("task_request", taskReq)
	defer func() {
		if hadPreviousTaskReq {
			c.Set("task_request", previousTaskReq)
			return
		}
		delete(c.Keys, "task_request")
	}()

	n := item.N
	if n <= 0 {
		n = 1
	}
	ratios := adaptor.EstimateBilling(c, info)
	if len(ratios) > 0 {
		multiplier := float64(n)
		for _, ratio := range ratios {
			multiplier *= ratio
		}
		return multiplier, true
	}
	if workflowVideoChannelUsesBasePerCallBilling(channelType) {
		return float64(n), true
	}
	return 0, false
}

func workflowVideoTaskSubmitReq(item workflowQuotaQuoteItem) relaycommon.TaskSubmitReq {
	seconds := ""
	if item.Seconds > 0 {
		seconds = strconv.Itoa(item.Seconds)
	}
	metadata := map[string]interface{}{}
	if strings.TrimSpace(item.Resolution) != "" {
		metadata["resolution"] = item.Resolution
	}
	if strings.TrimSpace(item.AspectRatio) != "" {
		metadata["aspect_ratio"] = item.AspectRatio
	}
	return relaycommon.TaskSubmitReq{
		Prompt:   item.Prompt,
		Model:    item.Model,
		Size:     item.Resolution,
		Duration: item.Seconds,
		Seconds:  seconds,
		Metadata: metadata,
	}
}

func workflowVideoChannelUsesBasePerCallBilling(channelType int) bool {
	switch channelType {
	case constant.ChannelTypeGrokVideo,
		constant.ChannelTypeJunliai,
		constant.ChannelTypeKieVideo,
		constant.ChannelTypeRunningHub,
		constant.ChannelTypeLconVideo,
		constant.ChannelTypeMiniMax,
		constant.ChannelTypeKling,
		constant.ChannelTypeJimeng,
		constant.ChannelTypeVidu:
		return true
	default:
		return false
	}
}

// workflowVideoResolutionUpperRatio 使用当前渠道中的保守倍率上界，防止视频节点低估费用。
func workflowVideoResolutionUpperRatio(resolution string) float64 {
	switch strings.ToLower(strings.TrimSpace(resolution)) {
	case "480p", "540p", "720p", "":
		return 1
	case "1080p":
		return 5
	case "2k", "1440p":
		return 8
	case "4k", "2160p":
		return 16
	default:
		return 5
	}
}

// reserveQuotaWithBuffer 为响应后才能确认的实际费用预留缓冲，多余部分在终态释放。
func reserveQuotaWithBuffer(quota int) int {
	if quota <= 0 {
		return 0
	}
	ratio := 1.1
	configured, err := strconv.ParseFloat(
		strings.TrimSpace(os.Getenv("WORKFLOW_QUOTA_RESERVE_BUFFER_RATIO")),
		64,
	)
	if err == nil && configured >= 1 && configured <= 3 {
		ratio = configured
	}
	return int(math.Ceil(float64(quota) * ratio))
}

// workflowQuotaRequestAuthorized 使用常量时间比较校验 SSO 与 NewAPI 的共享密钥。
func workflowQuotaRequestAuthorized(c *gin.Context) bool {
	expected := strings.TrimSpace(common.GetEnvOrDefaultString("BILLING_TRACE_SECRET", ""))
	actual := strings.TrimSpace(c.GetHeader(common.BillingTraceSecretKey))
	if expected == "" || len(expected) != len(actual) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1
}

// writeWorkflowQuotaModelError 将事务错误统一映射为稳定 HTTP 状态和机器错误码。
func writeWorkflowQuotaModelError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrWorkflowQuotaUserInsufficient):
		workflowQuotaError(c, http.StatusForbidden, string(types.ErrorCodeInsufficientUserQuota), "用户额度不足，无法预占整个工作流费用")
	case errors.Is(err, model.ErrWorkflowQuotaTokenInsufficient):
		workflowQuotaError(c, http.StatusForbidden, string(types.ErrorCodePreConsumeTokenQuotaFailed), "API Token 额度不足，无法预占整个工作流费用")
	case errors.Is(err, model.ErrWorkflowQuotaReservationExpired):
		workflowQuotaError(c, http.StatusConflict, string(types.ErrorCodeWorkflowQuotaReservationExpired), "工作流额度预占单已过期")
	case errors.Is(err, model.ErrWorkflowQuotaItemExceeded):
		workflowQuotaError(c, http.StatusConflict, string(types.ErrorCodeWorkflowQuotaItemExceeded), "节点实际所需额度超过工作流预占额度")
	case errors.Is(err, model.ErrWorkflowQuotaBatchUpdateEnabled):
		workflowQuotaError(c, http.StatusServiceUnavailable, "workflow_quota_configuration_error", "工作流额度预占要求关闭 BATCH_UPDATE_ENABLED")
	case errors.Is(err, model.ErrWorkflowQuotaReservationNotFound),
		errors.Is(err, model.ErrWorkflowQuotaReservationConflict),
		errors.Is(err, model.ErrWorkflowQuotaReservationClosed),
		errors.Is(err, model.ErrWorkflowQuotaItemNotFound),
		errors.Is(err, model.ErrWorkflowQuotaItemConflict):
		workflowQuotaError(c, http.StatusConflict, string(types.ErrorCodeWorkflowQuotaReservationInvalid), err.Error())
	default:
		workflowQuotaError(c, http.StatusInternalServerError, "workflow_quota_internal_error", err.Error())
	}
}

// workflowQuotaError 保持与 NewAPI 现有 error 响应兼容，便于 SSO 统一解析。
func workflowQuotaError(c *gin.Context, status int, code string, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"type":    code,
			"code":    code,
			"message": message,
		},
	})
}

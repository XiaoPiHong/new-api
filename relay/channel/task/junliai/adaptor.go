package junliai

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/relay/channel"
	taskcommon "github.com/QuantumNous/new-api/relay/channel/task/taskcommon"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
)

const maxReferenceImageBytes = 20 * 1024 * 1024

type responseTask struct {
	ID          string         `json:"id"`
	TaskID      string         `json:"task_id,omitempty"`
	RequestID   string         `json:"request_id,omitempty"`
	Object      string         `json:"object,omitempty"`
	Model       string         `json:"model,omitempty"`
	Status      string         `json:"status"`
	Progress    int            `json:"progress,omitempty"`
	CreatedAt   int64          `json:"created_at,omitempty"`
	CompletedAt int64          `json:"completed_at,omitempty"`
	VideoURL    string         `json:"video_url,omitempty"`
	URL         string         `json:"url,omitempty"`
	Video       map[string]any `json:"video,omitempty"`
	Data        map[string]any `json:"data,omitempty"`
	Error       *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type TaskAdaptor struct {
	taskcommon.BaseBilling
	apiKey  string
	baseURL string
}

func (a *TaskAdaptor) Init(info *relaycommon.RelayInfo) {
	a.baseURL = strings.TrimRight(info.ChannelBaseUrl, "/")
	a.apiKey = info.ApiKey
}

func (a *TaskAdaptor) ValidateRequestAndSetAction(c *gin.Context, info *relaycommon.RelayInfo) *dto.TaskError {
	if taskErr := relaycommon.ValidateMultipartDirect(c, info); taskErr != nil {
		return taskErr
	}
	bodyMap, err := getRequestBodyMap(c)
	if err != nil {
		return service.TaskErrorWrapperLocal(err, "invalid_request", http.StatusBadRequest)
	}
	if hasImageValue(bodyMap["video"]) {
		return service.TaskErrorWrapperLocal(fmt.Errorf("Junliai 通用视频接口不支持 video 编辑参数，请使用 reference_videos 提交参考视频"), "invalid_request", http.StatusBadRequest)
	}
	info.Action = constant.TaskActionTextGenerate
	if hasReferenceImageInput(bodyMap) || hasImageValue(bodyMap["reference_videos"]) {
		info.Action = constant.TaskActionReferenceGenerate
	}

	req, err := relaycommon.GetTaskRequest(c)
	if err != nil {
		return service.TaskErrorWrapperLocal(err, "invalid_request", http.StatusBadRequest)
	}
	if strings.TrimSpace(req.Image) != "" && len(req.Images) == 0 {
		req.Images = []string{req.Image}
		c.Set("task_request", req)
	}
	return nil
}

func (a *TaskAdaptor) BuildRequestURL(info *relaycommon.RelayInfo) (string, error) {
	return fmt.Sprintf("%s%s", a.baseURL, VideoGenerationEndpoint), nil
}

func (a *TaskAdaptor) BuildRequestHeader(c *gin.Context, req *http.Request, info *relaycommon.RelayInfo) error {
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return nil
}

func (a *TaskAdaptor) BuildRequestBody(c *gin.Context, info *relaycommon.RelayInfo) (io.Reader, error) {
	storage, err := common.GetBodyStorage(c)
	if err != nil {
		return nil, errors.Wrap(err, "get_request_body_failed")
	}
	cachedBody, err := storage.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "read_body_bytes_failed")
	}

	bodyMap, err := parseRequestBodyMap(cachedBody)
	if err != nil {
		return nil, err
	}
	bodyMap["model"] = strings.TrimSpace(info.UpstreamModelName)
	if image := stringFromMap(bodyMap, "image"); image != "" {
		if _, ok := bodyMap["images"]; !ok {
			bodyMap["images"] = []string{image}
		}
	}
	if inputReference := stringFromMap(bodyMap, "input_reference"); inputReference != "" {
		if _, ok := bodyMap["images"]; !ok {
			bodyMap["images"] = []string{inputReference}
		}
	}
	if err := a.normalizeReferenceImagesDataURLObjects(bodyMap, info); err != nil {
		return nil, err
	}
	for _, field := range []string{"start_frame", "end_frame"} {
		value, exists := bodyMap[field]
		if !exists {
			continue
		}
		imageURL, err := mediaURLString(value)
		if err != nil {
			return nil, errors.Wrap(err, field)
		}
		dataURL, err := a.imageReferenceDataURL(imageURL, info)
		if err != nil {
			return nil, errors.Wrap(err, field)
		}
		bodyMap[field] = dataURL
	}
	normalizeJunliaiVideoRequestBody(bodyMap)

	newBody, err := common.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(newBody), nil
}

func (a *TaskAdaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error) {
	// 参数覆盖在 BuildRequestBody 之后执行，保留旧的 reference_images.0.url 取值路径；
	// 真正发送前再转成通用接口的字符串格式，首尾帧覆盖规则无需随协议切换而重写。
	body, err := io.ReadAll(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "read_request_body_failed")
	}
	bodyMap, err := parseRequestBodyMap(body)
	if err != nil {
		return nil, err
	}
	if err := normalizeJunliaiVideoMedia(bodyMap); err != nil {
		return nil, err
	}
	body, err = common.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}
	return channel.DoTaskApiRequest(a, c, info, bytes.NewReader(body))
}

func (a *TaskAdaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (taskID string, taskData []byte, taskErr *dto.TaskError) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		taskErr = service.TaskErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError)
		return
	}
	_ = resp.Body.Close()

	var dResp responseTask
	if err := common.Unmarshal(responseBody, &dResp); err != nil {
		taskErr = service.TaskErrorWrapper(errors.Wrapf(err, "body: %s", responseBody), "unmarshal_response_body_failed", http.StatusInternalServerError)
		return
	}

	upstreamID := firstNonEmpty(dResp.ID, dResp.TaskID, dResp.RequestID, stringFromMap(dResp.Data, "id"), stringFromMap(dResp.Data, "task_id"), stringFromMap(dResp.Data, "request_id"))
	if upstreamID == "" {
		taskErr = service.TaskErrorWrapper(fmt.Errorf("task_id is empty"), "invalid_response", http.StatusInternalServerError)
		return
	}

	dResp.ID = info.PublicTaskID
	dResp.TaskID = info.PublicTaskID
	c.JSON(http.StatusOK, dResp)
	return upstreamID, responseBody, nil
}

func (a *TaskAdaptor) FetchTask(baseURL, key string, body map[string]any, proxy string) (*http.Response, error) {
	taskID, ok := body["task_id"].(string)
	if !ok || strings.TrimSpace(taskID) == "" {
		return nil, fmt.Errorf("invalid task_id")
	}

	uri := fmt.Sprintf("%s%s/%s", strings.TrimRight(baseURL, "/"), QueryTaskEndpoint, taskID)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Accept", "application/json")

	client, err := service.GetHttpClientWithProxy(proxy)
	if err != nil {
		return nil, fmt.Errorf("new proxy http client failed: %w", err)
	}
	return client.Do(req)
}

func (a *TaskAdaptor) ParseTaskResult(respBody []byte) (*relaycommon.TaskInfo, error) {
	resTask := responseTask{}
	if err := common.Unmarshal(respBody, &resTask); err != nil {
		return nil, errors.Wrap(err, "unmarshal task result failed")
	}

	status := strings.ToLower(strings.TrimSpace(firstNonEmpty(resTask.Status, stringFromMap(resTask.Data, "status"))))
	taskResult := relaycommon.TaskInfo{Code: 0}

	switch status {
	case "queued", "pending", "submitted":
		taskResult.Status = model.TaskStatusQueued
		taskResult.Progress = taskcommon.ProgressQueued
	case "processing", "running", "in_progress":
		taskResult.Status = model.TaskStatusInProgress
		taskResult.Progress = taskcommon.ProgressInProgress
	case "completed", "succeeded", "success", "done":
		taskResult.Status = model.TaskStatusSuccess
		taskResult.Progress = taskcommon.ProgressComplete
		taskResult.Url = firstJunliaiVideoResultURL(respBody)
	case "failed", "fail", "cancelled", "canceled", "expired":
		taskResult.Status = model.TaskStatusFailure
		taskResult.Progress = taskcommon.ProgressComplete
		taskResult.Reason = failureReason(resTask)
		if status == "expired" && taskResult.Reason == "task failed" {
			taskResult.Reason = "Junliai 视频任务已过期"
		}
	default:
		taskResult.Status = model.TaskStatusInProgress
		taskResult.Progress = taskcommon.ProgressInProgress
	}

	if taskResult.Status != model.TaskStatusSuccess && taskResult.Status != model.TaskStatusFailure && resTask.Progress > 0 && resTask.Progress < 100 {
		taskResult.Progress = fmt.Sprintf("%d%%", resTask.Progress)
	}
	return &taskResult, nil
}

func (a *TaskAdaptor) ConvertToOpenAIVideo(task *model.Task) ([]byte, error) {
	data := task.Data
	var err error
	if data, err = sjson.SetBytes(data, "id", task.TaskID); err != nil {
		return nil, errors.Wrap(err, "set id failed")
	}
	if data, err = sjson.SetBytes(data, "task_id", task.TaskID); err != nil {
		return nil, errors.Wrap(err, "set task_id failed")
	}
	if url := task.GetResultURL(); url != "" {
		if data, err = sjson.SetBytes(data, "url", url); err != nil {
			return nil, errors.Wrap(err, "set url failed")
		}
	}
	return data, nil
}

func (a *TaskAdaptor) GetModelList() []string {
	return ModelList
}

func (a *TaskAdaptor) GetChannelName() string {
	return ChannelName
}

func getRequestBodyMap(c *gin.Context) (map[string]any, error) {
	storage, err := common.GetBodyStorage(c)
	if err != nil {
		return nil, errors.Wrap(err, "get_request_body_failed")
	}
	cachedBody, err := storage.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "read_body_bytes_failed")
	}
	return parseRequestBodyMap(cachedBody)
}

func parseRequestBodyMap(body []byte) (map[string]any, error) {
	var bodyMap map[string]any
	if err := common.Unmarshal(body, &bodyMap); err != nil {
		return nil, errors.Wrap(err, "unmarshal_request_body_failed")
	}
	return bodyMap, nil
}

func hasReferenceImageInput(bodyMap map[string]any) bool {
	return hasImageValue(bodyMap["reference_images"]) ||
		hasImageValue(bodyMap["images"]) ||
		hasImageValue(bodyMap["start_frame"]) ||
		hasImageValue(bodyMap["end_frame"]) ||
		stringFromMap(bodyMap, "image") != "" ||
		stringFromMap(bodyMap, "input_reference") != ""
}

func hasImageValue(value any) bool {
	switch images := value.(type) {
	case string:
		return strings.TrimSpace(images) != ""
	case []string:
		for _, image := range images {
			if strings.TrimSpace(image) != "" {
				return true
			}
		}
	case []any:
		for _, image := range images {
			if hasImageValue(image) {
				return true
			}
		}
	case map[string]any:
		return firstNonEmpty(
			stringFromMap(images, "url"),
			stringFromMap(images, "imageUrl"),
			stringFromMap(images, "image_url"),
			stringFromMap(images, "dataUrl"),
			stringFromMap(images, "data_url"),
		) != ""
	case map[string]string:
		for _, image := range images {
			if strings.TrimSpace(image) != "" {
				return true
			}
		}
	}
	return false
}

func normalizeJunliaiVideoRequestBody(bodyMap map[string]any) {
	allowed := map[string]struct{}{
		"model":            {},
		"prompt":           {},
		"duration":         {},
		"size":             {},
		"aspect_ratio":     {},
		"resolution":       {},
		"reference_mode":   {},
		"reference_images": {},
		"start_frame":      {},
		"end_frame":        {},
		"reference_videos": {},
		"audio_reference":  {},
		"audio_references": {},
		"audio":            {},
		"response_format":  {},
	}
	if duration, ok := positiveNumberFromAny(firstMapValue(bodyMap, "duration", "seconds")); ok {
		bodyMap["duration"] = duration
	}
	if _, ok := bodyMap["aspect_ratio"]; !ok {
		if aspectRatio := stringFromMap(bodyMap, "aspectRatio"); aspectRatio != "" {
			bodyMap["aspect_ratio"] = aspectRatio
		}
	}
	// 兼容原 Grok 请求的音频字段，通用接口使用 audio_references。
	if _, ok := bodyMap["audio_references"]; !ok {
		if value, exists := bodyMap["reference_audios"]; exists {
			bodyMap["audio_references"] = value
		}
	}
	// 通用接口需显式请求 URL，否则完成态只提供需要鉴权的 /content 下载入口。
	if _, ok := bodyMap["response_format"]; !ok {
		bodyMap["response_format"] = "url"
	}
	for key := range bodyMap {
		if _, ok := allowed[key]; !ok {
			delete(bodyMap, key)
		}
	}
}

// 通用视频协议使用媒体字符串，不接受 Grok 的 {url: ...} 图片对象。
// 此转换只发生在 Junliai 发送阶段，避免影响共享的参数覆盖引擎和其他渠道。
func normalizeJunliaiVideoMedia(bodyMap map[string]any) error {
	for _, field := range []string{"reference_images", "reference_videos", "audio_reference", "audio_references"} {
		value, exists := bodyMap[field]
		if !exists {
			continue
		}
		urls, err := mediaURLStrings(value)
		if err != nil {
			return errors.Wrap(err, field)
		}
		if len(urls) == 0 {
			delete(bodyMap, field)
		} else {
			bodyMap[field] = urls
		}
	}
	for _, field := range []string{"start_frame", "end_frame"} {
		value, exists := bodyMap[field]
		if !exists {
			continue
		}
		url, err := mediaURLString(value)
		if err != nil {
			return errors.Wrap(err, field)
		}
		if url == "" {
			delete(bodyMap, field)
		} else {
			bodyMap[field] = url
		}
	}
	// 覆盖规则可以把普通参考图改为首尾帧，互斥校验必须放在覆盖之后。
	hasStart := hasImageValue(bodyMap["start_frame"])
	hasEnd := hasImageValue(bodyMap["end_frame"])
	if hasEnd && !hasStart {
		return fmt.Errorf("Junliai end_frame 必须搭配 start_frame")
	}
	if (hasStart || hasEnd) && (hasImageValue(bodyMap["reference_images"]) || hasImageValue(bodyMap["reference_videos"])) {
		return fmt.Errorf("Junliai 首尾帧不能与普通图片或视频参考混用")
	}
	if stringFromMap(bodyMap, "response_format") == "" {
		bodyMap["response_format"] = "url"
	}
	return nil
}

func mediaURLStrings(value any) ([]string, error) {
	var values []any
	switch v := value.(type) {
	case []any:
		values = v
	case []string:
		for _, item := range v {
			values = append(values, item)
		}
	default:
		values = []any{value}
	}
	urls := make([]string, 0, len(values))
	for _, item := range values {
		url, err := mediaURLString(item)
		if err != nil {
			return nil, err
		}
		if url != "" {
			urls = append(urls, url)
		}
	}
	return urls, nil
}

func mediaURLString(value any) (string, error) {
	switch v := value.(type) {
	case nil:
		return "", nil
	case string:
		return strings.TrimSpace(v), nil
	case map[string]any:
		for _, key := range []string{"url", "imageUrl", "image_url", "dataUrl", "data_url"} {
			if url, ok := v[key].(string); ok && strings.TrimSpace(url) != "" {
				return strings.TrimSpace(url), nil
			}
		}
	}
	return "", fmt.Errorf("媒体参数必须为字符串或包含 url 的对象")
}

func firstJunliaiVideoResultURL(respBody []byte) string {
	var payload any
	if err := common.Unmarshal(respBody, &payload); err != nil {
		return ""
	}
	for _, candidate := range collectJunliaiVideoResultURLs(payload) {
		if isAbsoluteDownloadURL(candidate) {
			return candidate
		}
	}
	return ""
}

func collectJunliaiVideoResultURLs(value any) []string {
	switch v := value.(type) {
	case string:
		url := strings.TrimSpace(v)
		if isPotentialDownloadURL(url) {
			return []string{url}
		}
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			result = append(result, collectJunliaiVideoResultURLs(item)...)
		}
		return dedupeStrings(result)
	case map[string]any:
		result := make([]string, 0)
		for _, key := range []string{
			"data",
			"results",
			"result",
			"output",
			"video",
			"metadata",
			"url",
			"video_url",
			"videoUrl",
			"result_url",
			"resultUrl",
			"resultURL",
			"file_url",
			"fileUrl",
			"download_url",
			"downloadUrl",
		} {
			result = append(result, collectJunliaiVideoResultURLs(v[key])...)
		}
		return dedupeStrings(result)
	}
	return nil
}

func dedupeStrings(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func isPotentialDownloadURL(url string) bool {
	return strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://") ||
		strings.HasPrefix(url, "data:") ||
		strings.HasPrefix(url, "/")
}

func isAbsoluteDownloadURL(url string) bool {
	return strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://") ||
		strings.HasPrefix(url, "data:")
}

func firstMapValue(m map[string]any, keys ...string) any {
	for _, key := range keys {
		if value, ok := m[key]; ok && value != nil {
			return value
		}
	}
	return nil
}

func positiveNumberFromAny(value any) (any, bool) {
	switch v := value.(type) {
	case int:
		if v > 0 {
			return v, true
		}
	case int64:
		if v > 0 {
			return v, true
		}
	case float64:
		if v > 0 {
			return v, true
		}
	case float32:
		if v > 0 {
			return v, true
		}
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err == nil && parsed > 0 {
			return parsed, true
		}
	}
	return nil, false
}

func (a *TaskAdaptor) normalizeReferenceImagesDataURLObjects(bodyMap map[string]any, info *relaycommon.RelayInfo) error {
	if images, err := a.imageDataURLObjects(bodyMap["reference_images"], info); err != nil {
		return err
	} else if len(images) > 0 {
		bodyMap["reference_images"] = images
	} else if images, err := a.imageDataURLObjects(bodyMap["images"], info); err != nil {
		return err
	} else if len(images) > 0 {
		bodyMap["reference_images"] = images
	}
	delete(bodyMap, "images")
	delete(bodyMap, "image")
	delete(bodyMap, "input_reference")
	return nil
}

func (a *TaskAdaptor) imageDataURLObjects(value any, info *relaycommon.RelayInfo) ([]map[string]any, error) {
	switch images := value.(type) {
	case []string:
		result := make([]map[string]any, 0, len(images))
		for _, image := range images {
			imageObj, err := a.imageDataURLObject(image, info)
			if err != nil {
				return nil, err
			}
			if imageObj != nil {
				result = append(result, imageObj)
			}
		}
		return result, nil
	case []any:
		result := make([]map[string]any, 0, len(images))
		for _, image := range images {
			imageObj, err := a.imageDataURLObject(image, info)
			if err != nil {
				return nil, err
			}
			if imageObj != nil {
				result = append(result, imageObj)
			}
		}
		return result, nil
	default:
		imageObj, err := a.imageDataURLObject(value, info)
		if err != nil {
			return nil, err
		}
		if imageObj != nil {
			return []map[string]any{imageObj}, nil
		}
	}
	return nil, nil
}

func (a *TaskAdaptor) imageDataURLObject(value any, info *relaycommon.RelayInfo) (map[string]any, error) {
	switch image := value.(type) {
	case string:
		if image = strings.TrimSpace(image); image == "" {
			return nil, nil
		}
		dataURL, err := a.imageReferenceDataURL(image, info)
		if err != nil {
			return nil, err
		}
		return map[string]any{"url": dataURL}, nil
	case map[string]any:
		result := make(map[string]any, len(image))
		for key, value := range image {
			result[key] = value
		}
		imageURL := firstNonEmpty(
			stringFromMap(result, "url"),
			stringFromMap(result, "imageUrl"),
			stringFromMap(result, "image_url"),
			stringFromMap(result, "dataUrl"),
			stringFromMap(result, "data_url"),
		)
		if imageURL == "" {
			return result, nil
		}
		dataURL, err := a.imageReferenceDataURL(imageURL, info)
		if err != nil {
			return nil, err
		}
		result["url"] = dataURL
		delete(result, "imageUrl")
		delete(result, "image_url")
		delete(result, "dataUrl")
		delete(result, "data_url")
		return result, nil
	case map[string]string:
		result := make(map[string]any, len(image))
		for key, value := range image {
			result[key] = value
		}
		imageURL := firstNonEmpty(
			stringFromMap(result, "url"),
			stringFromMap(result, "imageUrl"),
			stringFromMap(result, "image_url"),
			stringFromMap(result, "dataUrl"),
			stringFromMap(result, "data_url"),
		)
		if imageURL == "" {
			return result, nil
		}
		dataURL, err := a.imageReferenceDataURL(imageURL, info)
		if err != nil {
			return nil, err
		}
		result["url"] = dataURL
		delete(result, "imageUrl")
		delete(result, "image_url")
		delete(result, "dataUrl")
		delete(result, "data_url")
		return result, nil
	}
	return nil, nil
}

func (a *TaskAdaptor) imageReferenceDataURL(image string, info *relaycommon.RelayInfo) (string, error) {
	image = strings.TrimSpace(image)
	if image == "" {
		return "", nil
	}
	lower := strings.ToLower(image)
	if strings.HasPrefix(lower, "data:") {
		return image, nil
	}
	if !strings.HasPrefix(lower, "http://") && !strings.HasPrefix(lower, "https://") {
		return image, nil
	}
	return a.downloadImageAsDataURL(image, info)
}

func (a *TaskAdaptor) downloadImageAsDataURL(imageURL string, info *relaycommon.RelayInfo) (string, error) {
	proxy := ""
	if info != nil {
		proxy = info.ChannelSetting.Proxy
	}
	client, err := service.GetHttpClientWithProxy(proxy)
	if err != nil {
		return "", fmt.Errorf("new proxy http client failed: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download reference image failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("download reference image failed: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxReferenceImageBytes+1))
	if err != nil {
		return "", err
	}
	if len(data) > maxReferenceImageBytes {
		return "", fmt.Errorf("reference image is too large")
	}

	contentType := referenceImageContentType(resp.Header.Get("Content-Type"), data)
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(data)), nil
}

func referenceImageContentType(header string, data []byte) string {
	if mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(header)); err == nil && mediaType != "" && mediaType != "application/octet-stream" {
		return mediaType
	}
	detected := http.DetectContentType(data)
	if mediaType, _, err := mime.ParseMediaType(detected); err == nil && mediaType != "" {
		return mediaType
	}
	if detected != "" {
		return detected
	}
	return "application/octet-stream"
}

func stringFromMap(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return strings.TrimSpace(val)
	case fmt.Stringer:
		return strings.TrimSpace(val.String())
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", val))
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func failureReason(task responseTask) string {
	if task.Error != nil && task.Error.Message != "" {
		return task.Error.Message
	}
	if task.ErrorMessage != "" {
		return task.ErrorMessage
	}
	if s := stringFromMap(task.Data, "error_message"); s != "" {
		return s
	}
	return "task failed"
}

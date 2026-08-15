package kievideo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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
)

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
	bodyMap, err := getRequestBodyMap(c)
	if err != nil {
		return service.TaskErrorWrapperLocal(err, "invalid_request", http.StatusBadRequest)
	}
	inputMap := mapFromMap(bodyMap, "input")

	prompt := firstNonEmpty(stringFromMap(inputMap, "prompt"), stringFromMap(bodyMap, "prompt"))
	if strings.TrimSpace(prompt) == "" {
		return service.TaskErrorWrapperLocal(fmt.Errorf("prompt is required"), "invalid_request", http.StatusBadRequest)
	}

	req := relaycommon.TaskSubmitReq{
		Prompt:   prompt,
		Model:    firstNonEmpty(stringFromMap(bodyMap, "model"), info.OriginModelName),
		Images:   imageURLsFromMaps(bodyMap, inputMap),
		Size:     firstNonEmpty(stringFromMap(bodyMap, "size"), stringFromMap(inputMap, "aspect_ratio"), stringFromMap(bodyMap, "aspect_ratio")),
		Duration: positiveIntFromAny(firstMapValue(inputMap, "duration"), positiveIntFromAny(firstMapValue(bodyMap, "duration"), 0)),
		Seconds:  stringFromMap(bodyMap, "seconds"),
	}
	if req.Duration <= 0 {
		req.Duration = positiveIntFromAny(firstMapValue(bodyMap, "seconds"), 0)
	}
	info.Action = constant.TaskActionGenerate
	c.Set("task_request", req)
	return nil
}

func (a *TaskAdaptor) ApplyParamOverrideBeforeBuildRequest() bool {
	return true
}

func (a *TaskAdaptor) BuildRequestURL(_ *relaycommon.RelayInfo) (string, error) {
	return fmt.Sprintf("%s%s", a.baseURL, CreateTaskEndpoint), nil
}

func (a *TaskAdaptor) BuildRequestHeader(_ *gin.Context, req *http.Request, _ *relaycommon.RelayInfo) error {
	req.Header.Set("Authorization", authorizationHeader(a.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return nil
}

func (a *TaskAdaptor) BuildRequestBody(c *gin.Context, info *relaycommon.RelayInfo) (io.Reader, error) {
	bodyMap, err := getRequestBodyMap(c)
	if err != nil {
		return nil, err
	}

	taskReq, err := relaycommon.GetTaskRequest(c)
	if err != nil {
		return nil, err
	}

	inputMap := mapFromMap(bodyMap, "input")
	modelName := firstNonEmpty(info.UpstreamModelName, stringFromMap(bodyMap, "model"), taskReq.Model)
	prompt := firstNonEmpty(stringFromMap(inputMap, "prompt"), stringFromMap(bodyMap, "prompt"), taskReq.Prompt)
	if modelName == "" {
		return nil, fmt.Errorf("model is required")
	}
	if prompt == "" {
		return nil, fmt.Errorf("prompt is required")
	}

	inputPayload := copyStringAnyMap(inputMap)
	inputPayload["prompt"] = prompt
	if images := imageURLsFromRequest(taskReq, bodyMap, inputMap); len(images) > 0 {
		inputPayload["image_urls"] = images
	}
	inputPayload["aspect_ratio"] = resolveAspectRatio(taskReq, bodyMap, inputMap)
	inputPayload["resolution"] = firstNonEmpty(stringFromMap(inputMap, "resolution"), stringFromMap(bodyMap, "resolution"), "480p")
	inputPayload["duration"] = resolveDuration(taskReq, bodyMap, inputMap)
	delete(inputPayload, "imageUrls")
	delete(inputPayload, "aspectRatio")
	delete(inputPayload, "seconds")

	payload := map[string]any{
		"model": modelName,
		"input": inputPayload,
	}
	if callbackURL := firstNonEmpty(
		stringFromMap(bodyMap, "callBackUrl"),
		stringFromMap(bodyMap, "callbackUrl"),
		stringFromMap(bodyMap, "callback_url"),
	); callbackURL != "" {
		payload["callBackUrl"] = callbackURL
	}

	data, err := common.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func (a *TaskAdaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error) {
	return channel.DoTaskApiRequest(a, c, info, requestBody)
}

func (a *TaskAdaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (taskID string, taskData []byte, taskErr *dto.TaskError) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		taskErr = service.TaskErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError)
		return
	}
	_ = resp.Body.Close()

	var bodyMap map[string]any
	if err := common.Unmarshal(responseBody, &bodyMap); err != nil {
		taskErr = service.TaskErrorWrapper(errors.Wrapf(err, "body: %s", responseBody), "unmarshal_response_body_failed", http.StatusInternalServerError)
		return
	}
	if !isSuccessEnvelope(bodyMap) {
		taskErr = service.TaskErrorWrapperLocal(fmt.Errorf("%s", failureReasonFromPayload(bodyMap)), "task_failed", http.StatusBadRequest)
		return
	}

	upstreamID := taskIDFromPayload(bodyMap)
	if upstreamID == "" {
		taskErr = service.TaskErrorWrapper(fmt.Errorf("task_id is empty"), "invalid_response", http.StatusInternalServerError)
		return
	}

	ov := dto.NewOpenAIVideo()
	ov.ID = info.PublicTaskID
	ov.TaskID = info.PublicTaskID
	ov.CreatedAt = time.Now().Unix()
	ov.Model = info.OriginModelName
	c.JSON(http.StatusOK, ov)
	return upstreamID, responseBody, nil
}

func (a *TaskAdaptor) FetchTask(baseURL, key string, body map[string]any, proxy string) (*http.Response, error) {
	taskID, ok := body["task_id"].(string)
	if !ok || strings.TrimSpace(taskID) == "" {
		return nil, fmt.Errorf("invalid task_id")
	}

	uri := fmt.Sprintf("%s%s?taskId=%s", strings.TrimRight(baseURL, "/"), QueryTaskEndpoint, url.QueryEscape(strings.TrimSpace(taskID)))
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authorizationHeader(key))
	req.Header.Set("Accept", "application/json")

	client, err := service.GetHttpClientWithProxy(proxy)
	if err != nil {
		return nil, fmt.Errorf("new proxy http client failed: %w", err)
	}
	return client.Do(req)
}

func (a *TaskAdaptor) ParseTaskResult(respBody []byte) (*relaycommon.TaskInfo, error) {
	var payload map[string]any
	if err := common.Unmarshal(respBody, &payload); err != nil {
		return nil, errors.Wrap(err, "unmarshal task result failed")
	}

	statusText := normalizedTaskState(payload)
	taskResult := relaycommon.TaskInfo{Code: 0}

	switch statusText {
	case "waiting", "queuing", "queued", "pending", "submitted":
		taskResult.Status = model.TaskStatusQueued
		taskResult.Progress = taskcommon.ProgressQueued
	case "generating", "processing", "running", "in_progress":
		taskResult.Status = model.TaskStatusInProgress
		taskResult.Progress = taskcommon.ProgressInProgress
	case "success", "succeeded", "completed", "complete", "done":
		taskResult.Status = model.TaskStatusSuccess
		taskResult.Progress = taskcommon.ProgressComplete
		taskResult.Url = firstVideoResultURL(payload)
	case "fail", "failed", "failure", "error", "cancelled", "canceled":
		taskResult.Status = model.TaskStatusFailure
		taskResult.Progress = taskcommon.ProgressComplete
		taskResult.Reason = failureReasonFromPayload(payload)
	default:
		if url := firstVideoResultURL(payload); url != "" {
			taskResult.Status = model.TaskStatusSuccess
			taskResult.Progress = taskcommon.ProgressComplete
			taskResult.Url = url
		} else if !isSuccessEnvelope(payload) {
			taskResult.Status = model.TaskStatusFailure
			taskResult.Progress = taskcommon.ProgressComplete
			taskResult.Reason = failureReasonFromPayload(payload)
		} else {
			taskResult.Status = model.TaskStatusInProgress
			taskResult.Progress = taskcommon.ProgressInProgress
		}
	}
	return &taskResult, nil
}

func (a *TaskAdaptor) ConvertToOpenAIVideo(task *model.Task) ([]byte, error) {
	openAIVideo := dto.NewOpenAIVideo()
	openAIVideo.ID = task.TaskID
	openAIVideo.TaskID = task.TaskID
	openAIVideo.Status = task.Status.ToVideoStatus()
	openAIVideo.SetProgressStr(task.Progress)
	openAIVideo.CreatedAt = task.CreatedAt
	openAIVideo.CompletedAt = task.UpdatedAt
	openAIVideo.Model = task.Properties.OriginModelName
	if url := task.GetResultURL(); url != "" {
		openAIVideo.SetMetadata("url", url)
		openAIVideo.SetMetadata("video_url", url)
	}
	if task.Status == model.TaskStatusFailure {
		openAIVideo.Error = &dto.OpenAIVideoError{
			Message: firstNonEmpty(task.FailReason, "task failed"),
			Code:    "task_failed",
		}
	}
	return common.Marshal(openAIVideo)
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
	var bodyMap map[string]any
	if err := common.Unmarshal(cachedBody, &bodyMap); err != nil {
		return nil, errors.Wrap(err, "unmarshal_request_body_failed")
	}
	return bodyMap, nil
}

func imageURLsFromRequest(req relaycommon.TaskSubmitReq, maps ...map[string]any) []string {
	values := make([]string, 0)
	values = append(values, req.Images...)
	values = append(values, req.Image, req.InputReference)
	for _, m := range maps {
		values = append(values, imageURLsFromMaps(m)...)
	}
	return compactStrings(values)
}

func imageURLsFromMaps(maps ...map[string]any) []string {
	values := make([]string, 0)
	for _, m := range maps {
		values = append(values, stringsFromAny(m["image_urls"])...)
		values = append(values, stringsFromAny(m["imageUrls"])...)
		values = append(values, stringsFromAny(m["images"])...)
		values = append(values, stringsFromAny(m["image"])...)
		values = append(values, stringsFromAny(m["input_reference"])...)
		values = append(values, stringsFromAny(m["reference_images"])...)
	}
	return compactStrings(values)
}

func resolveAspectRatio(req relaycommon.TaskSubmitReq, bodyMap map[string]any, inputMap map[string]any) string {
	if ratio := firstNonEmpty(
		stringFromMap(inputMap, "aspect_ratio"),
		stringFromMap(bodyMap, "aspect_ratio"),
		stringFromMap(inputMap, "aspectRatio"),
		stringFromMap(bodyMap, "aspectRatio"),
	); ratio != "" {
		return ratio
	}
	if ratio := ratioFromSize(firstNonEmpty(stringFromMap(bodyMap, "size"), req.Size)); ratio != "" {
		return ratio
	}
	return "16:9"
}

func resolveDuration(req relaycommon.TaskSubmitReq, bodyMap map[string]any, inputMap map[string]any) int {
	duration := positiveIntFromAny(firstMapValue(inputMap, "duration"), 0)
	if duration <= 0 {
		duration = positiveIntFromAny(firstMapValue(bodyMap, "duration"), 0)
	}
	if duration <= 0 {
		duration = positiveIntFromAny(firstMapValue(inputMap, "seconds"), 0)
	}
	if duration <= 0 {
		duration = positiveIntFromAny(firstMapValue(bodyMap, "seconds"), 0)
	}
	if duration <= 0 {
		duration = req.Duration
	}
	if duration <= 0 {
		duration = positiveIntFromAny(req.Seconds, 0)
	}
	if duration <= 0 {
		duration = 8
	}
	return duration
}

func normalizedTaskState(payload map[string]any) string {
	dataMap := mapFromMap(payload, "data")
	state := firstNonEmpty(
		stringFromMap(dataMap, "state"),
		stringFromMap(dataMap, "status"),
		stringFromMap(dataMap, "taskStatus"),
		stringFromMap(dataMap, "task_status"),
		stringFromMap(payload, "state"),
		stringFromMap(payload, "status"),
	)
	if state != "" {
		return strings.ToLower(strings.TrimSpace(state))
	}
	switch intFromAny(firstMapValue(dataMap, "successFlag", "success_flag"), -1) {
	case 0:
		return "generating"
	case 1:
		return "success"
	case 2, 3:
		return "fail"
	default:
		return ""
	}
}

func taskIDFromPayload(payload map[string]any) string {
	dataMap := mapFromMap(payload, "data")
	return firstNonEmpty(
		stringFromMap(dataMap, "taskId"),
		stringFromMap(dataMap, "task_id"),
		stringFromMap(payload, "taskId"),
		stringFromMap(payload, "task_id"),
		stringFromMap(payload, "id"),
	)
}

func firstVideoResultURL(payload map[string]any) string {
	for _, candidate := range collectVideoResultURLs(payload) {
		if isPotentialDownloadURL(candidate) {
			return candidate
		}
	}
	return ""
}

func collectVideoResultURLs(value any) []string {
	switch v := value.(type) {
	case string:
		text := strings.TrimSpace(v)
		if text == "" {
			return nil
		}
		if strings.HasPrefix(text, "{") || strings.HasPrefix(text, "[") {
			var nested any
			if err := common.Unmarshal([]byte(text), &nested); err == nil {
				return collectVideoResultURLs(nested)
			}
		}
		if isPotentialDownloadURL(text) {
			return []string{text}
		}
	case []string:
		return compactStrings(v)
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			result = append(result, collectVideoResultURLs(item)...)
		}
		return compactStrings(result)
	case map[string]any:
		result := make([]string, 0)
		for _, key := range []string{
			"resultJson",
			"resultJSON",
			"result",
			"results",
			"resultUrls",
			"result_urls",
			"videoUrls",
			"video_urls",
			"videos",
			"urls",
			"data",
			"output",
			"video",
			"url",
			"resultUrl",
			"result_url",
			"videoUrl",
			"video_url",
			"fileUrl",
			"file_url",
			"downloadUrl",
			"download_url",
		} {
			result = append(result, collectVideoResultURLs(v[key])...)
		}
		return compactStrings(result)
	}
	return nil
}

func failureReasonFromPayload(payload map[string]any) string {
	dataMap := mapFromMap(payload, "data")
	if errorMap := mapFromMap(payload, "error"); errorMap != nil {
		if message := stringFromMap(errorMap, "message"); message != "" {
			return message
		}
	}
	return firstNonEmpty(
		stringFromMap(dataMap, "failMsg"),
		stringFromMap(dataMap, "fail_msg"),
		stringFromMap(dataMap, "failReason"),
		stringFromMap(dataMap, "errorMessage"),
		stringFromMap(dataMap, "error_message"),
		stringFromMap(dataMap, "message"),
		stringFromMap(dataMap, "msg"),
		stringFromMap(payload, "errorMessage"),
		stringFromMap(payload, "error_message"),
		stringFromMap(payload, "message"),
		stringFromMap(payload, "msg"),
		"task failed",
	)
}

func isSuccessEnvelope(payload map[string]any) bool {
	if payload == nil {
		return true
	}
	code, hasCode := payload["code"]
	message := strings.ToLower(firstNonEmpty(stringFromMap(payload, "msg"), stringFromMap(payload, "message")))
	if message == "success" || message == "ok" {
		return true
	}
	if !hasCode || code == nil {
		return true
	}
	switch value := code.(type) {
	case float64:
		return value == 0 || value == 200
	case int:
		return value == 0 || value == 200
	case string:
		normalized := strings.ToLower(strings.TrimSpace(value))
		return normalized == "" || normalized == "0" || normalized == "200" || normalized == "success" || normalized == "ok"
	default:
		normalized := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", code)))
		return normalized == "" || normalized == "0" || normalized == "200" || normalized == "success" || normalized == "ok"
	}
}

func isPotentialDownloadURL(value string) bool {
	return strings.HasPrefix(value, "http://") ||
		strings.HasPrefix(value, "https://") ||
		strings.HasPrefix(value, "data:")
}

func ratioFromSize(size string) string {
	size = strings.TrimSpace(strings.ToLower(size))
	if size == "" || strings.Contains(size, ":") {
		return size
	}
	parts := strings.FieldsFunc(size, func(r rune) bool {
		return r == 'x' || r == '*'
	})
	if len(parts) != 2 {
		return ""
	}
	width, errW := strconv.Atoi(strings.TrimSpace(parts[0]))
	height, errH := strconv.Atoi(strings.TrimSpace(parts[1]))
	if errW != nil || errH != nil || width <= 0 || height <= 0 {
		return ""
	}
	gcd := gcdInt(width, height)
	return fmt.Sprintf("%d:%d", width/gcd, height/gcd)
}

func gcdInt(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func stringsFromAny(value any) []string {
	switch v := value.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			result = append(result, stringsFromAny(item)...)
		}
		return result
	case map[string]any:
		return []string{firstNonEmpty(stringFromMap(v, "url"), stringFromMap(v, "imageUrl"), stringFromMap(v, "image_url"))}
	}
	return nil
}

func compactStrings(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
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

func positiveIntFromAny(value any, fallback int) int {
	switch v := value.(type) {
	case int:
		if v > 0 {
			return v
		}
	case int64:
		if v > 0 {
			return int(v)
		}
	case float64:
		if v > 0 {
			return int(v)
		}
	case float32:
		if v > 0 {
			return int(v)
		}
	case string:
		if parsed, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func intFromAny(value any, fallback int) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case float32:
		return int(v)
	case string:
		if parsed, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return parsed
		}
	}
	return fallback
}

func firstMapValue(m map[string]any, keys ...string) any {
	for _, key := range keys {
		if m == nil {
			return nil
		}
		if value, ok := m[key]; ok && value != nil {
			return value
		}
	}
	return nil
}

func copyStringAnyMap(m map[string]any) map[string]any {
	result := make(map[string]any, len(m))
	for key, value := range m {
		result[key] = value
	}
	return result
}

func mapFromMap(m map[string]any, key string) map[string]any {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	if val, ok := v.(map[string]any); ok {
		return val
	}
	return nil
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

func authorizationHeader(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	lower := strings.ToLower(key)
	if strings.HasPrefix(lower, "bearer ") || strings.HasPrefix(lower, "basic ") {
		return key
	}
	return "Bearer " + key
}

package runninghub

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

type requestPayload struct {
	Prompt      string   `json:"prompt"`
	AspectRatio string   `json:"aspectRatio"`
	ImageUrls   []string `json:"imageUrls"`
	Resolution  string   `json:"resolution,omitempty"`
	Duration    int      `json:"duration"`
}

type submitResponse struct {
	Code         any    `json:"code,omitempty"`
	Message      string `json:"message,omitempty"`
	Msg          string `json:"msg,omitempty"`
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	TaskID       string `json:"taskId,omitempty"`
	TaskId       string `json:"task_id,omitempty"`
	Status       string `json:"status,omitempty"`
	Data         struct {
		TaskID       string `json:"taskId,omitempty"`
		TaskId       string `json:"task_id,omitempty"`
		Status       string `json:"status,omitempty"`
		Message      string `json:"message,omitempty"`
		Msg          string `json:"msg,omitempty"`
		ErrorCode    string `json:"errorCode,omitempty"`
		ErrorMessage string `json:"errorMessage,omitempty"`
	} `json:"data,omitempty"`
}

type queryResponse struct {
	Code         any    `json:"code,omitempty"`
	Message      string `json:"message,omitempty"`
	Msg          string `json:"msg,omitempty"`
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	TaskID       string `json:"taskId,omitempty"`
	TaskId       string `json:"task_id,omitempty"`
	Status       string `json:"status,omitempty"`
	TaskStatus   string `json:"taskStatus,omitempty"`
	Results      any    `json:"results,omitempty"`
	URL          string `json:"url,omitempty"`
	ResultURL    string `json:"resultUrl,omitempty"`
	ResultUrl    string `json:"result_url,omitempty"`
	FileURL      string `json:"fileUrl,omitempty"`
	VideoURL     string `json:"videoUrl,omitempty"`
	Error        string `json:"error,omitempty"`
	Data         struct {
		TaskID       string `json:"taskId,omitempty"`
		TaskId       string `json:"task_id,omitempty"`
		Status       string `json:"status,omitempty"`
		TaskStatus   string `json:"taskStatus,omitempty"`
		Results      any    `json:"results,omitempty"`
		URL          string `json:"url,omitempty"`
		ResultURL    string `json:"resultUrl,omitempty"`
		ResultUrl    string `json:"result_url,omitempty"`
		FileURL      string `json:"fileUrl,omitempty"`
		VideoURL     string `json:"videoUrl,omitempty"`
		Error        string `json:"error,omitempty"`
		ErrorCode    string `json:"errorCode,omitempty"`
		ErrorMessage string `json:"errorMessage,omitempty"`
		Message      string `json:"message,omitempty"`
		Msg          string `json:"msg,omitempty"`
	} `json:"data,omitempty"`
}

type resultItem struct {
	URL       string `json:"url,omitempty"`
	Url       string `json:"Url,omitempty"`
	ResultURL string `json:"resultUrl,omitempty"`
	ResultUrl string `json:"result_url,omitempty"`
	FileURL   string `json:"fileUrl,omitempty"`
	FileUrl   string `json:"file_url,omitempty"`
	VideoURL  string `json:"videoUrl,omitempty"`
	VideoUrl  string `json:"video_url,omitempty"`
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
	if taskErr := relaycommon.ValidateBasicTaskRequest(c, info, constant.TaskActionGenerate); taskErr != nil {
		return taskErr
	}

	bodyMap, err := getRequestBodyMap(c)
	if err != nil {
		return service.TaskErrorWrapperLocal(err, "invalid_request", http.StatusBadRequest)
	}
	if len(imageURLsFromMap(bodyMap)) == 0 {
		return service.TaskErrorWrapperLocal(fmt.Errorf("imageUrls/images is required"), "invalid_request", http.StatusBadRequest)
	}
	info.Action = constant.TaskActionGenerate
	return nil
}

func (a *TaskAdaptor) BuildRequestURL(info *relaycommon.RelayInfo) (string, error) {
	modelSlug, err := modelSlugFromRelayInfo(info)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", a.baseURL, fmt.Sprintf(ImageToVideoEndpointFormat, modelSlug)), nil
}

func (a *TaskAdaptor) BuildRequestHeader(_ *gin.Context, req *http.Request, _ *relaycommon.RelayInfo) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	return nil
}

func (a *TaskAdaptor) BuildRequestBody(c *gin.Context, _ *relaycommon.RelayInfo) (io.Reader, error) {
	bodyMap, err := getRequestBodyMap(c)
	if err != nil {
		return nil, err
	}

	images := imageURLsFromMap(bodyMap)
	if len(images) == 0 {
		return nil, fmt.Errorf("imageUrls/images is required")
	}

	aspectRatio := firstNonEmpty(
		stringFromMap(bodyMap, "aspectRatio"),
		stringFromMap(bodyMap, "aspect_ratio"),
		aspectRatioFromSize(stringFromMap(bodyMap, "size")),
		"16:9",
	)
	duration := positiveIntFromMap(bodyMap, "duration", positiveIntFromMap(bodyMap, "seconds", 5))

	payload := requestPayload{
		Prompt:      stringFromMap(bodyMap, "prompt"),
		AspectRatio: aspectRatio,
		ImageUrls:   images,
		Resolution:  firstNonEmpty(stringFromMap(bodyMap, "resolution"), "720p"),
		Duration:    duration,
	}

	if strings.TrimSpace(payload.Prompt) == "" {
		return nil, fmt.Errorf("prompt is required")
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

	var result submitResponse
	if err := common.Unmarshal(responseBody, &result); err != nil {
		taskErr = service.TaskErrorWrapper(errors.Wrapf(err, "body: %s", responseBody), "unmarshal_response_body_failed", http.StatusInternalServerError)
		return
	}

	if !isSuccessCode(result.Code) {
		message := submitFailureReason(result, "task submit failed")
		taskErr = service.TaskErrorWrapperLocal(fmt.Errorf("%s", message), "task_failed", http.StatusBadRequest)
		return
	}
	if isFailureStatus(firstNonEmpty(result.Status, result.Data.Status)) || hasSubmitError(result) {
		message := submitFailureReason(result, "task submit failed")
		taskErr = service.TaskErrorWrapperLocal(fmt.Errorf("%s", message), "task_failed", http.StatusBadRequest)
		return
	}

	upstreamID := firstNonEmpty(result.TaskID, result.TaskId, result.Data.TaskID, result.Data.TaskId)
	if upstreamID == "" {
		message := firstNonEmpty(result.Message, result.Data.Message, result.Msg, result.Data.Msg, "task_id is empty")
		taskErr = service.TaskErrorWrapperLocal(fmt.Errorf("%s", message), "invalid_response", http.StatusInternalServerError)
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

	requestBody := map[string]any{
		"taskId": taskID,
	}
	data, err := common.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s%s", strings.TrimRight(baseURL, "/"), QueryTaskEndpoint)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	client, err := service.GetHttpClientWithProxy(proxy)
	if err != nil {
		return nil, fmt.Errorf("new proxy http client failed: %w", err)
	}
	return client.Do(req)
}

func (a *TaskAdaptor) ParseTaskResult(respBody []byte) (*relaycommon.TaskInfo, error) {
	resTask := queryResponse{}
	if err := common.Unmarshal(respBody, &resTask); err != nil {
		return nil, errors.Wrap(err, "unmarshal task result failed")
	}

	status := strings.ToLower(firstNonEmpty(resTask.Status, resTask.TaskStatus, resTask.Data.Status, resTask.Data.TaskStatus))
	taskResult := relaycommon.TaskInfo{Code: 0}

	switch status {
	case "queued", "pending", "submitted", "waiting":
		taskResult.Status = model.TaskStatusQueued
		taskResult.Progress = taskcommon.ProgressQueued
	case "running", "processing", "in_progress":
		taskResult.Status = model.TaskStatusInProgress
		taskResult.Progress = taskcommon.ProgressInProgress
	case "success", "succeeded", "completed", "complete", "done":
		taskResult.Status = model.TaskStatusSuccess
		taskResult.Progress = taskcommon.ProgressComplete
		taskResult.Url = resultURL(resTask)
	case "failed", "fail", "failure", "error", "cancelled", "canceled":
		taskResult.Status = model.TaskStatusFailure
		taskResult.Progress = taskcommon.ProgressComplete
		taskResult.Reason = failureReason(resTask)
	default:
		if !isSuccessCode(resTask.Code) {
			taskResult.Status = model.TaskStatusFailure
			taskResult.Progress = taskcommon.ProgressComplete
			taskResult.Reason = failureReason(resTask)
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

func modelSlugFromRelayInfo(info *relaycommon.RelayInfo) (string, error) {
	modelName := ""
	if info != nil {
		modelName = firstNonEmpty(info.UpstreamModelName, info.OriginModelName)
	}
	modelSlug := runningHubModelSlug(modelName)
	if modelSlug == "" {
		return "", fmt.Errorf("invalid runninghub model: %s", modelName)
	}
	return modelSlug, nil
}

func runningHubModelSlug(modelName string) string {
	modelName = strings.TrimSpace(modelName)
	if modelName == "" {
		return ""
	}
	if strings.HasPrefix(modelName, ModelNamePrefix) {
		modelName = strings.TrimPrefix(modelName, ModelNamePrefix)
	}
	if !isValidRunningHubModelSlug(modelName) {
		return ""
	}
	return modelName
}

func isValidRunningHubModelSlug(slug string) bool {
	if slug == "" {
		return false
	}
	if strings.HasPrefix(slug, "/") || strings.HasSuffix(slug, "/") {
		return false
	}
	for _, segment := range strings.Split(slug, "/") {
		if !isValidRunningHubModelSlugSegment(segment) {
			return false
		}
	}
	return true
}

func isValidRunningHubModelSlugSegment(segment string) bool {
	if segment == "" || segment == "." || segment == ".." {
		return false
	}
	for _, r := range segment {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		switch r {
		case '-', '_', '.':
			continue
		default:
			return false
		}
	}
	return true
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
		return nil, err
	}
	return bodyMap, nil
}

func imageURLsFromMap(bodyMap map[string]any) []string {
	if urls := stringSliceFromAny(bodyMap["imageUrls"]); len(urls) > 0 {
		return urls
	}
	if urls := stringSliceFromAny(bodyMap["image_urls"]); len(urls) > 0 {
		return urls
	}
	if urls := stringSliceFromAny(bodyMap["images"]); len(urls) > 0 {
		return urls
	}
	if image := firstNonEmpty(stringFromMap(bodyMap, "image"), stringFromMap(bodyMap, "input_reference")); image != "" {
		return []string{image}
	}
	return nil
}

func stringSliceFromAny(value any) []string {
	switch list := value.(type) {
	case []string:
		return compactStrings(list)
	case []any:
		result := make([]string, 0, len(list))
		for _, item := range list {
			switch v := item.(type) {
			case string:
				result = append(result, v)
			case map[string]any:
				result = append(result, firstNonEmpty(stringFromMap(v, "url"), stringFromMap(v, "imageUrl"), stringFromMap(v, "image_url")))
			}
		}
		return compactStrings(result)
	}
	return nil
}

func resultURL(task queryResponse) string {
	if url := firstNonEmpty(task.Data.URL, task.Data.ResultURL, task.Data.ResultUrl, task.Data.FileURL, task.Data.VideoURL, task.URL, task.ResultURL, task.ResultUrl, task.FileURL, task.VideoURL); url != "" {
		return url
	}
	if url := resultURLFromAny(task.Data.Results); url != "" {
		return url
	}
	return resultURLFromAny(task.Results)
}

func resultURLFromAny(value any) string {
	switch results := value.(type) {
	case []any:
		for _, item := range results {
			if url := resultURLFromAny(item); url != "" {
				return url
			}
		}
	case []resultItem:
		for _, item := range results {
			if url := firstNonEmpty(item.URL, item.Url, item.ResultURL, item.ResultUrl, item.FileURL, item.FileUrl, item.VideoURL, item.VideoUrl); url != "" {
				return url
			}
		}
	case map[string]any:
		return firstNonEmpty(
			stringFromMap(results, "url"),
			stringFromMap(results, "Url"),
			stringFromMap(results, "resultUrl"),
			stringFromMap(results, "result_url"),
			stringFromMap(results, "fileUrl"),
			stringFromMap(results, "file_url"),
			stringFromMap(results, "videoUrl"),
			stringFromMap(results, "video_url"),
		)
	case string:
		if url := strings.TrimSpace(results); url != "" {
			return url
		}
	}
	return ""
}

func failureReason(task queryResponse) string {
	return firstNonEmpty(task.Data.ErrorMessage, task.Data.Message, task.Data.Msg, task.Data.Error, task.ErrorMessage, task.Message, task.Msg, task.Error, "task failed")
}

func hasSubmitError(result submitResponse) bool {
	return firstNonEmpty(result.ErrorCode, result.ErrorMessage, result.Data.ErrorCode, result.Data.ErrorMessage) != ""
}

func submitFailureReason(result submitResponse, fallback string) string {
	message := firstNonEmpty(result.ErrorMessage, result.Data.ErrorMessage, result.Message, result.Data.Message, result.Msg, result.Data.Msg, fallback)
	errorCode := firstNonEmpty(result.ErrorCode, result.Data.ErrorCode)
	if errorCode == "" || strings.Contains(message, errorCode) {
		return message
	}
	return fmt.Sprintf("%s (errorCode=%s)", message, errorCode)
}

func isFailureStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "failed", "fail", "failure", "error", "cancelled", "canceled":
		return true
	default:
		return false
	}
}

func isSuccessCode(value any) bool {
	if value == nil {
		return true
	}
	switch code := value.(type) {
	case float64:
		return code == 0 || code == 200
	case int:
		return code == 0 || code == 200
	case string:
		normalized := strings.ToLower(strings.TrimSpace(code))
		return normalized == "" || normalized == "0" || normalized == "200" || normalized == "success" || normalized == "ok"
	default:
		normalized := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", value)))
		return normalized == "" || normalized == "0" || normalized == "200" || normalized == "success" || normalized == "ok"
	}
}

func positiveIntFromMap(bodyMap map[string]any, key string, fallback int) int {
	value := bodyMap[key]
	switch v := value.(type) {
	case float64:
		if v > 0 {
			return int(v)
		}
	case int:
		if v > 0 {
			return v
		}
	case string:
		if parsed, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func aspectRatioFromSize(size string) string {
	size = strings.TrimSpace(strings.ToLower(size))
	if size == "" {
		return ""
	}
	parts := strings.FieldsFunc(size, func(r rune) bool {
		return r == 'x' || r == '*' || r == ':'
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

func compactStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}
	return result
}

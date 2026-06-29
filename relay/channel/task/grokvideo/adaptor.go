package grokvideo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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
	info.Action = constant.TaskActionTextGenerate
	if stringFromMap(bodyMap, "video") != "" {
		info.Action = constant.TaskActionGenerate
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
	if info.Action == constant.TaskActionGenerate {
		return fmt.Sprintf("%s%s", a.baseURL, VideoEditEndpoint), nil
	}
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
	normalizeVideoURLObject(bodyMap)
	normalizeReferenceImagesURLObjects(bodyMap)

	newBody, err := common.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(newBody), nil
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
		taskResult.Url = firstNonEmpty(
			resTask.VideoURL,
			resTask.URL,
			stringFromMap(resTask.Video, "url"),
			stringFromMap(resTask.Data, "video_url"),
			stringFromMap(resTask.Data, "url"),
			stringFromMap(mapFromMap(resTask.Data, "video"), "url"),
		)
	case "failed", "fail", "cancelled", "canceled":
		taskResult.Status = model.TaskStatusFailure
		taskResult.Progress = taskcommon.ProgressComplete
		taskResult.Reason = failureReason(resTask)
	default:
		taskResult.Status = model.TaskStatusInProgress
		taskResult.Progress = taskcommon.ProgressInProgress
	}

	if resTask.Progress > 0 && resTask.Progress < 100 {
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

func normalizeVideoURLObject(bodyMap map[string]any) {
	if value := stringFromMap(bodyMap, "video"); value != "" {
		bodyMap["video"] = map[string]any{"url": value}
	}
}

func normalizeReferenceImagesURLObjects(bodyMap map[string]any) {
	if images := imageURLObjects(bodyMap["reference_images"]); len(images) > 0 {
		bodyMap["reference_images"] = images
	} else if images := imageURLObjects(bodyMap["images"]); len(images) > 0 {
		bodyMap["reference_images"] = images
	}
	delete(bodyMap, "images")
	delete(bodyMap, "image")
	delete(bodyMap, "input_reference")
}

func imageURLObjects(value any) []map[string]any {
	switch images := value.(type) {
	case []string:
		result := make([]map[string]any, 0, len(images))
		for _, image := range images {
			if image = strings.TrimSpace(image); image != "" {
				result = append(result, map[string]any{"url": image})
			}
		}
		return result
	case []any:
		result := make([]map[string]any, 0, len(images))
		for _, image := range images {
			if imageObj := imageURLObject(image); imageObj != nil {
				result = append(result, imageObj)
			}
		}
		return result
	default:
		if imageObj := imageURLObject(value); imageObj != nil {
			return []map[string]any{imageObj}
		}
	}
	return nil
}

func imageURLObject(value any) map[string]any {
	switch image := value.(type) {
	case string:
		if image = strings.TrimSpace(image); image != "" {
			return map[string]any{"url": image}
		}
	case map[string]any:
		return image
	case map[string]string:
		result := make(map[string]any, len(image))
		for key, value := range image {
			result[key] = value
		}
		return result
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

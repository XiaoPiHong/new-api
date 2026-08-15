package lconvideo

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/relay/channel"
	taskcommon "github.com/QuantumNous/new-api/relay/channel/task/taskcommon"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const maxInputReferenceBytes = 32 << 20

type responseTask struct {
	Code         any            `json:"code,omitempty"`
	Message      string         `json:"message,omitempty"`
	Msg          string         `json:"msg,omitempty"`
	ID           string         `json:"id,omitempty"`
	TaskID       string         `json:"task_id,omitempty"`
	TaskId       string         `json:"taskId,omitempty"`
	RequestID    string         `json:"request_id,omitempty"`
	Object       string         `json:"object,omitempty"`
	Model        string         `json:"model,omitempty"`
	Status       string         `json:"status,omitempty"`
	Progress     int            `json:"progress,omitempty"`
	CreatedAt    int64          `json:"created_at,omitempty"`
	CompletedAt  int64          `json:"completed_at,omitempty"`
	VideoURL     string         `json:"video_url,omitempty"`
	VideoUrl     string         `json:"videoUrl,omitempty"`
	URL          string         `json:"url,omitempty"`
	Video        map[string]any `json:"video,omitempty"`
	Output       map[string]any `json:"output,omitempty"`
	Data         map[string]any `json:"data,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	Error        *struct {
		Message string `json:"message,omitempty"`
		Code    string `json:"code,omitempty"`
	} `json:"error,omitempty"`
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
	return relaycommon.ValidateMultipartDirect(c, info)
}

func (a *TaskAdaptor) ApplyParamOverrideBeforeBuildRequest() bool {
	return true
}

func (a *TaskAdaptor) BuildRequestURL(_ *relaycommon.RelayInfo) (string, error) {
	return fmt.Sprintf("%s%s", a.baseURL, VideoEndpoint), nil
}

func (a *TaskAdaptor) BuildRequestHeader(c *gin.Context, req *http.Request, _ *relaycommon.RelayInfo) error {
	req.Header.Set("Authorization", authorizationHeader(a.apiKey))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", c.Request.Header.Get("Content-Type"))
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

	var bodyMap map[string]any
	_ = common.Unmarshal(cachedBody, &bodyMap)

	taskReq, err := relaycommon.GetTaskRequest(c)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writeFormField(writer, "model", firstNonEmpty(info.UpstreamModelName, taskReq.Model, stringFromMap(bodyMap, "model")))
	writeFormField(writer, "prompt", firstNonEmpty(taskReq.Prompt, stringFromMap(bodyMap, "prompt")))
	writeFormField(writer, "size", resolveSize(taskReq, bodyMap))
	writeFormField(writer, "seconds", resolveSeconds(taskReq, bodyMap))

	if err := a.writeInputReference(c, writer, taskReq, bodyMap, info); err != nil {
		_ = writer.Close()
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	return &buf, nil
}

func (a *TaskAdaptor) writeInputReference(c *gin.Context, writer *multipart.Writer, taskReq relaycommon.TaskSubmitReq, bodyMap map[string]any, info *relaycommon.RelayInfo) error {
	if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
		formData, err := common.ParseMultipartFormReusable(c)
		if err == nil {
			if files := formData.File["input_reference"]; len(files) > 0 {
				return copyMultipartFile(writer, "input_reference", files[0])
			}
		}
	}

	imageURL := firstImageURL(taskReq, bodyMap)
	if imageURL == "" {
		return nil
	}
	return a.writeRemoteFile(writer, imageURL, info)
}

func (a *TaskAdaptor) writeRemoteFile(writer *multipart.Writer, imageURL string, info *relaycommon.RelayInfo) error {
	client, err := service.GetHttpClientWithProxy(info.ChannelSetting.Proxy)
	if err != nil {
		return fmt.Errorf("new proxy http client failed: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download input_reference failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download input_reference failed: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxInputReferenceBytes+1))
	if err != nil {
		return err
	}
	if len(data) > maxInputReferenceBytes {
		return fmt.Errorf("input_reference is too large")
	}

	contentType := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = http.DetectContentType(data)
	}
	filename := inputReferenceFilename(imageURL, contentType)
	return writeMultipartFileBytes(writer, "input_reference", filename, contentType, data)
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
	if !isSuccessCode(dResp.Code) {
		taskErr = service.TaskErrorWrapperLocal(fmt.Errorf("%s", failureReason(dResp)), "task_failed", http.StatusBadRequest)
		return
	}

	upstreamID := taskIDFromResponse(dResp)
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

	uri := fmt.Sprintf("%s%s/%s", strings.TrimRight(baseURL, "/"), QueryTaskEndpoint, taskID)
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
	resTask := responseTask{}
	if err := common.Unmarshal(respBody, &resTask); err != nil {
		return nil, errors.Wrap(err, "unmarshal task result failed")
	}

	status := strings.ToLower(strings.TrimSpace(firstNonEmpty(resTask.Status, stringFromMap(resTask.Data, "status"))))
	taskResult := relaycommon.TaskInfo{Code: 0}

	switch status {
	case "queued", "pending", "submitted", "waiting":
		taskResult.Status = model.TaskStatusQueued
		taskResult.Progress = taskcommon.ProgressQueued
	case "processing", "running", "in_progress":
		taskResult.Status = model.TaskStatusInProgress
		taskResult.Progress = taskcommon.ProgressInProgress
	case "completed", "succeeded", "success", "done":
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

	if resTask.Progress > 0 && resTask.Progress < 100 {
		taskResult.Progress = fmt.Sprintf("%d%%", resTask.Progress)
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

func copyMultipartFile(writer *multipart.Writer, fieldName string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileHeader.Filename))
	header.Set("Content-Type", contentType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	return err
}

func writeMultipartFileBytes(writer *multipart.Writer, fieldName, filename, contentType string, data []byte) error {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filename))
	header.Set("Content-Type", contentType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}
	_, err = part.Write(data)
	return err
}

func writeFormField(writer *multipart.Writer, key, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	_ = writer.WriteField(key, strings.TrimSpace(value))
}

func resolveSize(req relaycommon.TaskSubmitReq, bodyMap map[string]any) string {
	size := firstNonEmpty(stringFromMap(bodyMap, "size"), req.Size)
	if ratio := ratioFromSize(size); ratio != "" {
		return ratio
	}
	return firstNonEmpty(
		stringFromMap(bodyMap, "aspect_ratio"),
		stringFromMap(bodyMap, "aspectRatio"),
		"16:9",
	)
}

func resolveSeconds(req relaycommon.TaskSubmitReq, bodyMap map[string]any) string {
	if value := firstNonEmpty(stringFromMap(bodyMap, "seconds"), req.Seconds, stringFromMap(bodyMap, "duration")); value != "" {
		return value
	}
	if req.Duration > 0 {
		return strconv.Itoa(req.Duration)
	}
	return "5"
}

func firstImageURL(req relaycommon.TaskSubmitReq, bodyMap map[string]any) string {
	if req.InputReference != "" {
		return req.InputReference
	}
	if req.Image != "" {
		return req.Image
	}
	if len(req.Images) > 0 {
		return firstNonEmpty(req.Images...)
	}
	return firstNonEmpty(
		stringFromMap(bodyMap, "input_reference"),
		stringFromMap(bodyMap, "image"),
		firstStringFromAny(bodyMap["images"]),
		firstStringFromAny(bodyMap["imageUrls"]),
		firstStringFromAny(bodyMap["image_urls"]),
		firstStringFromAny(bodyMap["reference_images"]),
	)
}

func firstStringFromAny(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case []string:
		return firstNonEmpty(v...)
	case []any:
		for _, item := range v {
			if s := firstStringFromAny(item); s != "" {
				return s
			}
		}
	case map[string]any:
		return firstNonEmpty(stringFromMap(v, "url"), stringFromMap(v, "imageUrl"), stringFromMap(v, "image_url"))
	}
	return ""
}

func inputReferenceFilename(imageURL, contentType string) string {
	ext := strings.ToLower(filepath.Ext(strings.Split(imageURL, "?")[0]))
	if ext == "" {
		if exts, err := mime.ExtensionsByType(contentType); err == nil && len(exts) > 0 {
			ext = exts[0]
		}
	}
	if ext == "" {
		ext = ".png"
	}
	return "input-reference" + ext
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

func taskIDFromResponse(task responseTask) string {
	return firstNonEmpty(
		task.ID,
		task.TaskID,
		task.TaskId,
		task.RequestID,
		stringFromMap(task.Data, "id"),
		stringFromMap(task.Data, "task_id"),
		stringFromMap(task.Data, "taskId"),
		stringFromMap(task.Data, "request_id"),
	)
}

func resultURL(task responseTask) string {
	return firstNonEmpty(
		task.VideoURL,
		task.VideoUrl,
		task.URL,
		stringFromMap(task.Video, "url"),
		stringFromMap(task.Output, "url"),
		stringFromMap(task.Output, "video_url"),
		stringFromMap(task.Output, "videoUrl"),
		stringFromMap(task.Data, "url"),
		stringFromMap(task.Data, "video_url"),
		stringFromMap(task.Data, "videoUrl"),
		stringFromMap(mapFromMap(task.Data, "video"), "url"),
		stringFromMap(mapFromMap(task.Data, "output"), "url"),
	)
}

func failureReason(task responseTask) string {
	if task.Error != nil && task.Error.Message != "" {
		return task.Error.Message
	}
	return firstNonEmpty(
		task.ErrorMessage,
		task.Message,
		task.Msg,
		stringFromMap(task.Data, "error_message"),
		stringFromMap(task.Data, "message"),
		stringFromMap(task.Data, "msg"),
		"task failed",
	)
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

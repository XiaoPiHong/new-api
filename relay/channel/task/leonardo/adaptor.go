package leonardotask

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
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

var dimensionsRE = regexp.MustCompile(`^(\d+)\s*[xX*]\s*(\d+)$`)

type TaskAdaptor struct {
	taskcommon.BaseBilling
	baseURL string
	apiKey  string
}

type submitResponse struct {
	ID string `json:"id"`
}

type adminResult struct {
	URLs         []string `json:"urls"`
	Status       string   `json:"status"`
	Error        any      `json:"error"`
	Ok           *bool    `json:"ok"`
	GenerationID string   `json:"generationId"`
}

type adminJob struct {
	ID      string        `json:"id"`
	Status  string        `json:"status"`
	Error   any           `json:"error"`
	Results []adminResult `json:"results"`
	Output  struct {
		Results []adminResult `json:"results"`
	} `json:"output"`
}

func (a *TaskAdaptor) Init(info *relaycommon.RelayInfo) {
	a.baseURL = strings.TrimRight(info.ChannelBaseUrl, "/")
	a.apiKey = strings.TrimSpace(info.ApiKey)
}

func (a *TaskAdaptor) ValidateRequestAndSetAction(c *gin.Context, info *relaycommon.RelayInfo) *dto.TaskError {
	if taskErr := relaycommon.ValidateMultipartDirect(c, info); taskErr != nil {
		return taskErr
	}
	// Reference images are forwarded to Leonardo Admin. The admin service is
	// responsible for uploading URL/data references to Leonardo first and then
	// submitting the resulting image IDs in the guidance slots selected by
	// reference_frame_types (for example StartFrame/EndFrame).
	// Do not reject images here: xphai-web sends them in the standard `images`
	// field and channel parameter overrides may add richer reference fields.
	return nil
}

func (a *TaskAdaptor) BuildRequestURL(_ *relaycommon.RelayInfo) (string, error) {
	if a.baseURL == "" {
		return "", errors.New("Leonardo Admin base URL is empty")
	}
	return a.baseURL + "/api/video/generate", nil
}

func (a *TaskAdaptor) BuildRequestHeader(_ *gin.Context, req *http.Request, _ *relaycommon.RelayInfo) error {
	req.Header.Set("Authorization", bearerHeader(a.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return nil
}

func (a *TaskAdaptor) BuildRequestBody(c *gin.Context, info *relaycommon.RelayInfo) (io.Reader, error) {
	req, err := relaycommon.GetTaskRequest(c)
	if err != nil {
		return nil, err
	}
	body := map[string]any{
		"prompt":  strings.TrimSpace(req.Prompt),
		"model":   strings.TrimSpace(info.UpstreamModelName),
		"count":   1,
		"strategy": "rr",
		"dryRun":  false,
	}
	if body["model"] == "" {
		body["model"] = strings.TrimSpace(req.Model)
	}
	if resolution := normalizeResolution(req.Size); resolution != "" {
		body["resolution"] = resolution
	}
	if seconds := parseSeconds(req); seconds > 0 {
		body["seconds"] = seconds
	}
	// Multipart requests are parsed into TaskSubmitReq before this method is
	// called, so the original body is not JSON and cannot be unmarshaled below.
	// Copy the normalized reference/extension fields explicitly to keep image
	// references and future model parameters working for both JSON and
	// multipart OpenAI video clients.
	if req.Mode != "" {
		body["mode"] = req.Mode
	}
	if req.Image != "" {
		body["image"] = req.Image
	}
	if len(req.Images) > 0 {
		body["images"] = req.Images
	}
	if req.InputReference != "" {
		body["input_reference"] = req.InputReference
	}
	for key, value := range req.Metadata {
		if _, reserved := body[key]; !reserved && key != "dryRun" {
			body[key] = value
		}
	}

	// Preserve explicitly supplied extension fields while forcing security-
	// sensitive routing fields above. This keeps future Leonardo options
	// configurable through new-api's parameter override mechanism.
	if storage, storageErr := common.GetBodyStorage(c); storageErr == nil {
		if raw, readErr := storage.Bytes(); readErr == nil {
			var incoming map[string]any
			if common.Unmarshal(raw, &incoming) == nil {
				for key, value := range incoming {
					if _, reserved := body[key]; !reserved && key != "dryRun" {
						body[key] = value
					}
				}
			}
		}
	}
	if resolution, ok := body["resolution"].(string); ok {
		body["resolution"] = normalizeResolution(resolution)
	}
	body["model"] = strings.TrimSpace(fmt.Sprintf("%v", body["model"]))
	body["dryRun"] = false
	encoded, err := common.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(encoded), nil
}

func (a *TaskAdaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error) {
	resp, err := channel.DoTaskApiRequest(a, c, info, requestBody)
	if err != nil {
		return nil, err
	}
	// Leonardo Admin intentionally returns 202 for the asynchronous enqueue,
	// while RelayTaskSubmit expects a normalized successful HTTP response before
	// invoking DoResponse.
	if resp != nil && resp.StatusCode == http.StatusAccepted {
		resp.StatusCode = http.StatusOK
		resp.Status = "200 OK"
	}
	return resp, nil
}

func (a *TaskAdaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (taskID string, taskData []byte, taskErr *dto.TaskError) {
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return "", nil, service.TaskErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError)
	}
	var submitted submitResponse
	if err = common.Unmarshal(body, &submitted); err != nil || strings.TrimSpace(submitted.ID) == "" {
		if err == nil {
			err = errors.New("missing job id")
		}
		return "", body, service.TaskErrorWrapper(errors.Wrap(err, "invalid Leonardo Admin video response"), "invalid_response", http.StatusBadGateway)
	}
	video := dto.NewOpenAIVideo()
	video.ID = info.PublicTaskID
	video.TaskID = info.PublicTaskID
	video.Model = info.OriginModelName
	video.Progress = 20
	video.CreatedAt = time.Now().Unix()
	c.JSON(http.StatusOK, video)
	return submitted.ID, body, nil
}

func (a *TaskAdaptor) FetchTask(baseURL, key string, body map[string]any, proxy string) (*http.Response, error) {
	taskID, ok := body["task_id"].(string)
	if !ok || strings.TrimSpace(taskID) == "" {
		return nil, errors.New("invalid task_id")
	}
	url := fmt.Sprintf("%s/api/jobs/%s", strings.TrimRight(baseURL, "/"), taskID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerHeader(key))
	req.Header.Set("Accept", "application/json")
	client, err := service.GetHttpClientWithProxy(proxy)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func (a *TaskAdaptor) ParseTaskResult(body []byte) (*relaycommon.TaskInfo, error) {
	var job adminJob
	if err := common.Unmarshal(body, &job); err != nil {
		return nil, errors.Wrap(err, "unmarshal Leonardo Admin job")
	}
	result := &relaycommon.TaskInfo{Code: 0}
	topStatus := strings.ToLower(strings.TrimSpace(job.Status))
	if topStatus == "" && job.Error != nil {
		result.Status = model.TaskStatusFailure
		result.Progress = taskcommon.ProgressComplete
		result.Reason = errorMessage(job.Error)
		return result, nil
	}
	switch topStatus {
	case "queued", "pending":
		result.Status = model.TaskStatusQueued
		result.Progress = taskcommon.ProgressQueued
	case "running", "processing", "in_progress":
		result.Status = model.TaskStatusInProgress
		result.Progress = taskcommon.ProgressInProgress
	case "complete", "completed", "success", "succeeded":
		result.Url = firstVideoURL(job)
		// Leonardo Admin wraps per-account failures in an outer job with
		// status=complete. Do not turn a failed result (or a completed result
		// without any media URL) into a false SUCCESS; otherwise new-api will
		// expose /content and the caller receives a misleading 409/502.
		if result.Url != "" {
			result.Status = model.TaskStatusSuccess
			result.Progress = taskcommon.ProgressComplete
			break
		}
		if failed, reason := adminJobFailure(job); failed {
			result.Status = model.TaskStatusFailure
			result.Progress = taskcommon.ProgressComplete
			result.Reason = reason
			break
		}
		if adminJobPending(job) {
			result.Status = model.TaskStatusInProgress
			result.Progress = taskcommon.ProgressInProgress
			break
		}
		result.Status = model.TaskStatusFailure
		result.Progress = taskcommon.ProgressComplete
		result.Reason = "Leonardo Admin completed without a video URL"
	case "error", "failed", "failure", "cancelled", "canceled":
		result.Status = model.TaskStatusFailure
		result.Progress = taskcommon.ProgressComplete
		result.Reason = errorMessage(job.Error)
	default:
		result.Status = model.TaskStatusInProgress
		result.Progress = taskcommon.ProgressInProgress
	}
	return result, nil
}

// adminJobResults returns both result locations used by the Admin API. The
// output location is used by some historical job serializers.
func adminJobResults(job adminJob) []adminResult {
	return append(append([]adminResult{}, job.Results...), job.Output.Results...)
}

func adminJobFailure(job adminJob) (bool, string) {
	if job.Error != nil {
		return true, errorMessage(job.Error)
	}
	for _, item := range adminJobResults(job) {
		status := strings.ToLower(strings.TrimSpace(item.Status))
		failed := item.Ok != nil && !*item.Ok
		failed = failed || status == "failed" || status == "failure" || status == "error" || status == "cancelled" || status == "canceled"
		if item.Error != nil {
			failed = true
		}
		if failed {
			reason := errorMessage(item.Error)
			if item.Error == nil {
				reason = firstNonEmpty(status, "Leonardo Admin video generation failed")
			}
			return true, reason
		}
	}
	return false, ""
}

func adminJobPending(job adminJob) bool {
	for _, item := range adminJobResults(job) {
		switch strings.ToLower(strings.TrimSpace(item.Status)) {
		case "queued", "pending", "running", "processing", "in_progress", "not_start", "submitted":
			return true
		}
	}
	return false
}

func (a *TaskAdaptor) ConvertToOpenAIVideo(task *model.Task) ([]byte, error) {
	video := dto.NewOpenAIVideo()
	video.ID = task.TaskID
	video.TaskID = task.TaskID
	video.Model = task.Properties.OriginModelName
	video.Status = task.Status.ToVideoStatus()
	video.SetProgressStr(task.Progress)
	video.CreatedAt = task.CreatedAt
	video.CompletedAt = task.UpdatedAt
	if url := task.GetResultURL(); url != "" {
		video.SetMetadata("url", taskcommon.BuildProxyURL(task.TaskID))
		video.SetMetadata("video_url", url)
	}
	if task.Status == model.TaskStatusFailure {
		video.Error = &dto.OpenAIVideoError{Message: firstNonEmpty(task.FailReason, "task failed"), Code: "task_failed"}
	}
	return common.Marshal(video)
}

func (a *TaskAdaptor) GetModelList() []string { return ModelList }
func (a *TaskAdaptor) GetChannelName() string { return ChannelName }

func parseSeconds(req relaycommon.TaskSubmitReq) int {
	if value, err := strconv.Atoi(strings.TrimSpace(req.Seconds)); err == nil && value > 0 {
		return value
	}
	return req.Duration
}

func normalizeResolution(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}
	switch value {
	case "480p", "720p", "768p", "768sq", "1080p", "1080sq", "1440p", "2160p", "2160sq":
		return value
	case "16:9", "9:16", "1:1", "4:3", "3:4":
		return "768p"
	}
	parts := dimensionsRE.FindStringSubmatch(value)
	if len(parts) != 3 {
		return value
	}
	height, _ := strconv.Atoi(parts[2])
	if height <= 500 {
		return "480p"
	}
	if height >= 1000 {
		return "1080p"
	}
	return "768p"
}

func firstVideoURL(job adminJob) string {
	for _, result := range append(job.Results, job.Output.Results...) {
		for _, url := range result.URLs {
			if strings.TrimSpace(url) != "" {
				return strings.TrimSpace(url)
			}
		}
	}
	return ""
}

func errorMessage(value any) string {
	if value == nil {
		return "Leonardo Admin video generation failed"
	}
	if message, ok := value.(string); ok && strings.TrimSpace(message) != "" {
		return strings.TrimSpace(message)
	}
	if object, ok := value.(map[string]any); ok {
		if message, ok := object["message"].(string); ok && strings.TrimSpace(message) != "" {
			return strings.TrimSpace(message)
		}
	}
	return fmt.Sprintf("%v", value)
}

func bearerHeader(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(key), "bearer ") {
		return key
	}
	return "Bearer " + key
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

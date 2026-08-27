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
	URLs []string `json:"urls"`
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
	req, err := relaycommon.GetTaskRequest(c)
	if err != nil {
		return service.TaskErrorWrapperLocal(err, "invalid_request", http.StatusBadRequest)
	}
	if req.HasImage() || strings.TrimSpace(req.InputReference) != "" {
		return service.TaskErrorWrapperLocal(errors.New("Leonardo Admin video image references are not supported yet"), "invalid_request", http.StatusBadRequest)
	}
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
	if strings.TrimSpace(job.Status) == "" && job.Error != nil {
		result.Status = model.TaskStatusFailure
		result.Progress = taskcommon.ProgressComplete
		result.Reason = errorMessage(job.Error)
		return result, nil
	}
	switch strings.ToLower(strings.TrimSpace(job.Status)) {
	case "queued", "pending":
		result.Status = model.TaskStatusQueued
		result.Progress = taskcommon.ProgressQueued
	case "running", "processing", "in_progress":
		result.Status = model.TaskStatusInProgress
		result.Progress = taskcommon.ProgressInProgress
	case "complete", "completed", "success", "succeeded":
		result.Status = model.TaskStatusSuccess
		result.Progress = taskcommon.ProgressComplete
		result.Url = firstVideoURL(job)
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

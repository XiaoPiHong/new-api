package leonardo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/relay/channel"
	"github.com/QuantumNous/new-api/relay/channel/openai"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Adaptor struct {
	baseURL string
	apiKey  string
}

type generateRequest struct {
	Prompt     string         `json:"prompt"`
	Model      string         `json:"model"`
	Size       string         `json:"size,omitempty"`
	Quality    string         `json:"quality,omitempty"`
	Count      int            `json:"count"`
	DryRun     bool           `json:"dryRun"`
	Parameters map[string]any `json:"parameters,omitempty"`
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

func (a *Adaptor) Init(info *relaycommon.RelayInfo) {
	a.baseURL = strings.TrimRight(info.ChannelBaseUrl, "/")
	a.apiKey = strings.TrimSpace(info.ApiKey)
}

func (a *Adaptor) GetRequestURL(_ *relaycommon.RelayInfo) (string, error) {
	if a.baseURL == "" {
		return "", errors.New("Leonardo Admin base URL is empty")
	}
	return a.baseURL + "/api/generate", nil
}

func (a *Adaptor) SetupRequestHeader(_ *gin.Context, header *http.Header, _ *relaycommon.RelayInfo) error {
	header.Set("Authorization", bearerHeader(a.apiKey))
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	return nil
}

func (a *Adaptor) ConvertImageRequest(_ *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	if info.RelayMode == relayconstant.RelayModeImagesEdits {
		return nil, errors.New("Leonardo Admin image edits are not supported yet; use image generations")
	}
	if request.N != nil && *request.N != 1 {
		return nil, errors.New("Leonardo Admin currently supports n=1 only")
	}
	modelName := strings.TrimSpace(info.UpstreamModelName)
	if modelName == "" {
		modelName = strings.TrimSpace(request.Model)
	}
	prompt := strings.TrimSpace(request.Prompt)
	if prompt == "" {
		return nil, errors.New("prompt is required")
	}
	size := strings.TrimSpace(request.Size)
	if size == "" {
		size = "1024x1024"
	}
	return generateRequest{
		Prompt:  prompt,
		Model:   modelName,
		Size:    size,
		Quality: strings.TrimSpace(request.Quality),
		Count:   1,
		DryRun:  false,
	}, nil
}

func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error) {
	resp, err := channel.DoApiRequest(a, c, info, requestBody)
	if err != nil {
		return nil, err
	}
	client, err := service.GetHttpClientWithProxy(info.ChannelSetting.Proxy)
	if err != nil {
		return nil, err
	}
	submitBody, readErr := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return responseWithStatus(resp.StatusCode, submitBody, "application/json"), nil
	}

	var submitted adminJob
	if err = common.Unmarshal(submitBody, &submitted); err != nil || strings.TrimSpace(submitted.ID) == "" {
		return responseWithStatus(http.StatusBadGateway, submitBody, "application/json"), nil
	}
	return a.pollImageJob(c, info, submitted.ID, client)
}

func (a *Adaptor) pollImageJob(c *gin.Context, info *relaycommon.RelayInfo, jobID string, client *http.Client) (*http.Response, error) {
	timeout := durationFromEnv("LEONARDO_ADMIN_IMAGE_POLL_TIMEOUT_MS", 90*time.Second)
	interval := durationFromEnv("LEONARDO_ADMIN_IMAGE_POLL_INTERVAL_MS", 2*time.Second)
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		body, status, err := a.fetchJob(c, info, jobID, client)
		if err != nil {
			return nil, err
		}
		var job adminJob
		if err = common.Unmarshal(body, &job); err != nil {
			return responseWithStatus(http.StatusBadGateway, body, "application/json"), nil
		}
		if status < http.StatusOK || status >= http.StatusMultipleChoices {
			return responseWithStatus(status, body, "application/json"), nil
		}
		state := strings.ToLower(strings.TrimSpace(job.Status))
		if state == "complete" || state == "error" {
			if state == "error" {
				return responseWithStatus(http.StatusBadGateway, body, "application/json"), nil
			}
			url := firstImageURL(job)
			if url == "" {
				return responseWithStatus(http.StatusBadGateway, []byte(`{"error":"Leonardo job completed without an image URL"}`), "application/json"), nil
			}
			payload, _ := common.Marshal(map[string]any{
				"created": time.Now().Unix(),
				"data":    []map[string]string{{"url": url}},
			})
			return responseWithStatus(http.StatusOK, payload, "application/json"), nil
		}

		select {
		case <-c.Request.Context().Done():
			return nil, c.Request.Context().Err()
		case <-deadline.C:
			return responseWithStatus(http.StatusGatewayTimeout, []byte(`{"error":"Leonardo image generation timed out"}`), "application/json"), nil
		case <-ticker.C:
		}
	}
}

func (a *Adaptor) fetchJob(c *gin.Context, info *relaycommon.RelayInfo, jobID string, client *http.Client) ([]byte, int, error) {
	url := fmt.Sprintf("%s/api/jobs/%s", a.baseURL, jobID)
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", bearerHeader(a.apiKey))
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (any, *types.NewAPIError) {
	return openai.OpenaiHandlerWithUsage(c, info, resp)
}

func (a *Adaptor) GetModelList() []string { return ModelList }
func (a *Adaptor) GetChannelName() string { return ChannelName }

func (a *Adaptor) ConvertOpenAIRequest(*gin.Context, *relaycommon.RelayInfo, *dto.GeneralOpenAIRequest) (any, error) {
	return nil, errors.New("Leonardo Admin channel only supports image and video endpoints")
}
func (a *Adaptor) ConvertRerankRequest(*gin.Context, int, dto.RerankRequest) (any, error) {
	return nil, errors.New("Leonardo Admin does not support rerank")
}
func (a *Adaptor) ConvertEmbeddingRequest(*gin.Context, *relaycommon.RelayInfo, dto.EmbeddingRequest) (any, error) {
	return nil, errors.New("Leonardo Admin does not support embeddings")
}
func (a *Adaptor) ConvertAudioRequest(*gin.Context, *relaycommon.RelayInfo, dto.AudioRequest) (io.Reader, error) {
	return nil, errors.New("Leonardo Admin does not support audio")
}
func (a *Adaptor) ConvertOpenAIResponsesRequest(*gin.Context, *relaycommon.RelayInfo, dto.OpenAIResponsesRequest) (any, error) {
	return nil, errors.New("Leonardo Admin does not support responses")
}
func (a *Adaptor) ConvertClaudeRequest(*gin.Context, *relaycommon.RelayInfo, *dto.ClaudeRequest) (any, error) {
	return nil, errors.New("Leonardo Admin does not support Claude messages")
}
func (a *Adaptor) ConvertGeminiRequest(*gin.Context, *relaycommon.RelayInfo, *dto.GeminiChatRequest) (any, error) {
	return nil, errors.New("Leonardo Admin does not support Gemini chat")
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

func firstImageURL(job adminJob) string {
	for _, result := range append(job.Results, job.Output.Results...) {
		for _, url := range result.URLs {
			if strings.TrimSpace(url) != "" {
				return strings.TrimSpace(url)
			}
		}
	}
	return ""
}

func responseWithStatus(status int, body []byte, contentType string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     http.Header{"Content-Type": []string{contentType}},
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}

func durationFromEnv(name string, fallback time.Duration) time.Duration {
	value, err := strconv.Atoi(strings.TrimSpace(os.Getenv(name)))
	if err != nil || value <= 0 {
		return fallback
	}
	return time.Duration(value) * time.Millisecond
}

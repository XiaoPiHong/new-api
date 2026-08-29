package leonardo

import (
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/relay/channel"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Adaptor is kept for the generic channel registry. Leonardo Admin requests
// are handled by the task adaptor (relay/channel/task/leonardo); this adapter
// intentionally exposes no image-generation path.
type Adaptor struct {
	baseURL string
	apiKey  string
}

func (a *Adaptor) Init(info *relaycommon.RelayInfo) {
	a.baseURL = strings.TrimRight(info.ChannelBaseUrl, "/")
	a.apiKey = strings.TrimSpace(info.ApiKey)
}

func (a *Adaptor) GetRequestURL(_ *relaycommon.RelayInfo) (string, error) {
	if a.baseURL == "" {
		return "", errors.New("Leonardo Admin base URL is empty")
	}
	return a.baseURL + "/api/video/generate", nil
}

func (a *Adaptor) SetupRequestHeader(_ *gin.Context, header *http.Header, _ *relaycommon.RelayInfo) error {
	header.Set("Authorization", bearerHeader(a.apiKey))
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	return nil
}

func (a *Adaptor) ConvertImageRequest(*gin.Context, *relaycommon.RelayInfo, dto.ImageRequest) (any, error) {
	return nil, errors.New("Leonardo Admin only supports the OpenAI video endpoint")
}

func (a *Adaptor) DoRequest(*gin.Context, *relaycommon.RelayInfo, io.Reader) (any, error) {
	return nil, errors.New("Leonardo Admin generic requests are disabled; use the video task endpoint")
}

func (a *Adaptor) DoResponse(*gin.Context, *http.Response, *relaycommon.RelayInfo) (any, *types.NewAPIError) {
	return nil, types.NewErrorWithStatusCode(
		errors.New("Leonardo Admin generic responses are disabled; use the video task endpoint"),
		types.ErrorCodeBadResponse,
		http.StatusNotImplemented,
	)
}

func (a *Adaptor) GetModelList() []string { return ModelList }
func (a *Adaptor) GetChannelName() string { return ChannelName }

func (a *Adaptor) ConvertOpenAIRequest(*gin.Context, *relaycommon.RelayInfo, *dto.GeneralOpenAIRequest) (any, error) {
	return nil, errors.New("Leonardo Admin only supports the OpenAI video endpoint")
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

// Keep the import below explicit in the generic registry contract. The task
// adaptor performs the actual HTTP request for video jobs.
var _ channel.Adaptor = (*Adaptor)(nil)

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

package gemini

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// convertGeminiGenerateContentImageRequest converts the OpenAI image request
// used by xphai-web into a native Gemini generateContent request. The
// GeminiImage channel type identifies the protocol, so no provider/model name
// is hard-coded here.
func convertGeminiGenerateContentImageRequest(
	c *gin.Context,
	info *relaycommon.RelayInfo,
	request dto.ImageRequest,
) (any, error) {
	if strings.TrimSpace(request.Prompt) == "" {
		return nil, errors.New("prompt is required")
	}

	parts := []dto.GeminiPart{{Text: request.Prompt}}
	imageParts, err := readGeminiMultipartImageParts(c)
	if err != nil {
		return nil, err
	}
	if len(imageParts) == 0 {
		imageParts, err = readGeminiRawImageParts(request.Image, request.Images)
		if err != nil {
			return nil, err
		}
	}
	if info != nil && info.RelayMode == relayconstant.RelayModeImagesEdits && len(imageParts) == 0 {
		return nil, errors.New("image is required for image edits")
	}
	parts = append(parts, imageParts...)

	imageConfig, err := buildGeminiImageConfig(request.Size)
	if err != nil {
		return nil, err
	}

	candidateCount := int(lo.FromPtrOr(request.N, uint(1)))
	if candidateCount <= 0 {
		candidateCount = 1
	}

	return &dto.GeminiChatRequest{
		Contents: []dto.GeminiChatContent{{
			Role:  "user",
			Parts: parts,
		}},
		GenerationConfig: dto.GeminiChatGenerationConfig{
			ResponseModalities: []string{"IMAGE"},
			CandidateCount:     lo.ToPtr(candidateCount),
			ImageConfig:        imageConfig,
		},
	}, nil
}

func readGeminiMultipartImageParts(c *gin.Context) ([]dto.GeminiPart, error) {
	if c == nil || c.Request == nil || !strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
		return nil, nil
	}

	multipartForm := c.Request.MultipartForm
	if multipartForm == nil {
		var err error
		multipartForm, err = c.MultipartForm()
		if err != nil {
			return nil, fmt.Errorf("failed to parse multipart form: %w", err)
		}
	}

	files := make([]*multipart.FileHeader, 0)
	// xphai-web uses image for one image and image[] for multiple images. Keep
	// the slice order supplied by net/http for image[].
	for _, fieldName := range []string{"image", "image[]"} {
		files = append(files, multipartForm.File[fieldName]...)
	}

	// Also accept image[0], image[1], ... for OpenAI-compatible clients. Map
	// iteration is sorted to keep the request deterministic.
	arrayFields := make([]string, 0)
	for fieldName := range multipartForm.File {
		if strings.HasPrefix(fieldName, "image[") && fieldName != "image[]" {
			arrayFields = append(arrayFields, fieldName)
		}
	}
	sort.Strings(arrayFields)
	for _, fieldName := range arrayFields {
		files = append(files, multipartForm.File[fieldName]...)
	}

	parts := make([]dto.GeminiPart, 0, len(files))
	for i, fileHeader := range files {
		if fileHeader == nil {
			continue
		}
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open image file %d: %w", i, err)
		}
		data, readErr := io.ReadAll(file)
		_ = file.Close()
		if readErr != nil {
			return nil, fmt.Errorf("failed to read image file %d: %w", i, readErr)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("image file %d is empty", i)
		}

		mimeType := fileHeader.Header.Get("Content-Type")
		if mimeType == "" || mimeType == "application/octet-stream" {
			mimeType = http.DetectContentType(data)
		}
		if !strings.HasPrefix(strings.ToLower(mimeType), "image/") {
			return nil, fmt.Errorf("unsupported image mime type %q", mimeType)
		}

		parts = append(parts, dto.GeminiPart{
			InlineData: &dto.GeminiInlineData{
				MimeType: mimeType,
				Data:     base64.StdEncoding.EncodeToString(data),
			},
		})
	}

	return parts, nil
}

func readGeminiRawImageParts(imageValue, imagesValue json.RawMessage) ([]dto.GeminiPart, error) {
	parts := make([]dto.GeminiPart, 0)
	if len(imageValue) > 0 && string(imageValue) != "null" {
		if err := appendGeminiRawImageParts(&parts, imageValue); err != nil {
			return nil, err
		}
	}
	if len(imagesValue) > 0 && string(imagesValue) != "null" {
		if err := appendGeminiRawImageParts(&parts, imagesValue); err != nil {
			return nil, err
		}
	}
	return parts, nil
}

func appendGeminiRawImageParts(parts *[]dto.GeminiPart, raw json.RawMessage) error {
	var value any
	if err := common.Unmarshal(raw, &value); err != nil {
		return fmt.Errorf("invalid image input: %w", err)
	}

	var appendValue func(any) error
	appendValue = func(item any) error {
		switch typed := item.(type) {
		case string:
			if strings.TrimSpace(typed) == "" {
				return nil
			}
			mimeType, data, err := service.DecodeBase64FileData(typed)
			if err != nil {
				return fmt.Errorf("invalid base64 image input: %w", err)
			}
			*parts = append(*parts, dto.GeminiPart{
				InlineData: &dto.GeminiInlineData{MimeType: mimeType, Data: data},
			})
			return nil
		case []any:
			for _, child := range typed {
				if err := appendValue(child); err != nil {
					return err
				}
			}
			return nil
		default:
			return errors.New("image input must be a data URL, base64 string, or array of those values")
		}
	}

	return appendValue(value)
}

func buildGeminiImageConfig(size string) (json.RawMessage, error) {
	size = strings.TrimSpace(size)
	imageSize := "1K"
	aspectRatio := ""

	switch size {
	case "", "1K", "1024x1024":
		// defaults
	case "2K", "2048x2048":
		imageSize = "2K"
	case "1536x1024":
		aspectRatio = "3:2"
	case "1024x1536":
		aspectRatio = "2:3"
	case "1792x1024":
		aspectRatio = "16:9"
	case "1024x1792":
		aspectRatio = "9:16"
	default:
		return nil, fmt.Errorf("unsupported Gemini image size %q", size)
	}

	config := map[string]string{"imageSize": imageSize}
	if aspectRatio != "" {
		config["aspectRatio"] = aspectRatio
	}
	encoded, err := common.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Gemini image config: %w", err)
	}
	return encoded, nil
}

// GeminiNativeImageHandler converts Gemini generateContent image parts to
// the OpenAI ImageResponse consumed by xphai-web.
func GeminiNativeImageHandler(
	c *gin.Context,
	info *relaycommon.RelayInfo,
	resp *http.Response,
) (*dto.Usage, *types.NewAPIError) {
	defer service.CloseResponseBodyGracefully(resp)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
	}

	var geminiResponse dto.GeminiChatResponse
	if err = common.Unmarshal(responseBody, &geminiResponse); err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
	}

	if len(geminiResponse.Candidates) == 0 {
		if geminiResponse.PromptFeedback != nil && geminiResponse.PromptFeedback.BlockReason != nil {
			common.SetContextKey(c, constant.ContextKeyAdminRejectReason, fmt.Sprintf("gemini_block_reason=%s", *geminiResponse.PromptFeedback.BlockReason))
			return nil, types.NewOpenAIError(
				errors.New("request blocked by Gemini API: "+*geminiResponse.PromptFeedback.BlockReason),
				types.ErrorCodePromptBlocked,
				http.StatusBadRequest,
			)
		}
		return nil, types.NewOpenAIError(errors.New("empty response from Gemini API"), types.ErrorCodeEmptyResponse, http.StatusInternalServerError)
	}

	data := make([]dto.ImageData, 0)
	for _, candidate := range geminiResponse.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.InlineData == nil || part.InlineData.Data == "" {
				continue
			}

			// xphai-web's artifact pipeline stores image responses as .png and
			// advertises them as image/png. Gemini may return JPEG bytes even
			// though the OpenAI-compatible response only has b64_json and no
			// per-item MIME field. Normalize the actual bytes here so the
			// downstream PNG contract remains valid.
			normalizedData, normalizeErr := normalizeGeminiImageToPNG(
				part.InlineData.MimeType,
				part.InlineData.Data,
			)
			if normalizeErr != nil {
				return nil, types.NewOpenAIError(normalizeErr, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
			}
			data = append(data, dto.ImageData{B64Json: normalizedData})
		}
	}
	if len(data) == 0 {
		return nil, types.NewOpenAIError(errors.New("no images generated"), types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
	}

	openAIResponse := dto.ImageResponse{
		Created: common.GetTimestamp(),
		Data:    data,
	}
	jsonResponse, err := common.Marshal(openAIResponse)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(resp.StatusCode)
	_, _ = c.Writer.Write(jsonResponse)

	usage := buildUsageFromGeminiMetadata(geminiResponse.UsageMetadata, info.GetEstimatePromptTokens())
	if usage.TotalTokens <= 0 {
		// Gemini may omit usage metadata for image responses. Keep billing
		// non-zero while preserving any prompt estimate when available.
		if usage.PromptTokens <= 0 {
			usage.PromptTokens = info.GetEstimatePromptTokens()
		}
		if usage.PromptTokens <= 0 {
			usage.PromptTokens = 1
		}
		usage.CompletionTokens = len(data) * 1400
		usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
	}

	return &usage, nil
}

// normalizeGeminiImageToPNG converts Gemini's inline image bytes to a real
// PNG. The xphai artifact pipeline currently persists all image results with
// a .png extension and image/png content type, so forwarding JPEG bytes
// unchanged would produce an invalid PNG artifact.
func normalizeGeminiImageToPNG(mimeType, encodedData string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("invalid Gemini image base64 (%s): %w", mimeType, err)
	}

	decoded, format, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("unable to decode Gemini image (%s): %w", mimeType, err)
	}

	var encoded bytes.Buffer
	if err := png.Encode(&encoded, decoded); err != nil {
		return "", fmt.Errorf("unable to encode Gemini %s image as PNG: %w", format, err)
	}

	return base64.StdEncoding.EncodeToString(encoded.Bytes()), nil
}

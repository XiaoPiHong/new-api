package relay

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/pkg/billingexpr"
	pluginruntime "github.com/QuantumNous/new-api/pkg/jsplugin"
	"github.com/QuantumNous/new-api/relay/channel"
	jspluginadaptor "github.com/QuantumNous/new-api/relay/channel/task/jsplugin"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/relay/helper"
	kitdto "github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTaskPlatformPriority(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)
	c.Set("platform", "fallback")
	assert.Equal(t, "fallback", string(GetTaskPlatform(c)))
	c.Set("channel_type", 59)
	assert.Equal(t, "59", string(GetTaskPlatform(c)))
	c.Set("task_plugin_key", "document-parser")
	assert.Equal(t, "document-parser", string(GetTaskPlatform(c)))
}

func TestLegacySoraMultipartRetainsVendorFieldsAndUploadedReference(t *testing.T) {
	service.InitHttpClient()
	source, err := os.ReadFile("../plugins/tasks/sora/plugin.js")
	require.NoError(t, err)
	plugin, err := pluginruntime.NewRegistry().RegisterFactory(string(source), pluginruntime.Options{Key: "sora"})
	require.NoError(t, err)
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	for key, value := range map[string]string{"model": "vendor-alias", "prompt": "animate", "seconds": "5", "video": "https://mock-media.invalid/source.mp4", "enabled": "false", "seed": "0", "metadata": `{"custom":true}`} {
		require.NoError(t, writer.WriteField(key, value))
	}
	require.NoError(t, writer.WriteField("vendor_options", "a"))
	require.NoError(t, writer.WriteField("vendor_options", "b"))
	file, err := writer.CreateFormFile("input_reference", "reference.png")
	require.NoError(t, err)
	_, err = file.Write([]byte("reference image bytes"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	captured := make(chan *multipart.Form, 1)
	referenceBytes := make(chan []byte, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer r.MultipartForm.RemoveAll()
		files := r.MultipartForm.File["input_reference"]
		if len(files) != 1 {
			http.Error(w, "missing reference", 400)
			return
		}
		file, err := files[0].Open()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		data, err := io.ReadAll(file)
		_ = file.Close()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		captured <- r.MultipartForm
		referenceBytes <- data
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, `{"id":"provider-id"}`)
	}))
	t.Cleanup(server.Close)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos", bytes.NewReader(requestBody.Bytes()))
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	t.Cleanup(func() {
		if storage, err := common.GetBodyStorage(c); err == nil {
			_ = storage.Close()
		}
	})
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{ChannelType: constant.ChannelTypeSora, ChannelBaseUrl: server.URL, ApiKey: "mock-key", UpstreamModelName: "vendor-upstream"}, TaskRelayInfo: &relaycommon.TaskRelayInfo{PublicTaskID: "public-task"}, OriginModelName: "vendor-alias"}
	adaptor := jspluginadaptor.New(plugin)
	adaptor.Init(info)
	require.Nil(t, prepareLegacySoraTaskRequest(c, info, "sora"))
	require.Nil(t, adaptor.ValidateRequestAndSetAction(c, info))
	body, err := adaptor.BuildRequestBody(c, info)
	require.NoError(t, err)
	response, err := adaptor.DoRequest(c, info, body)
	require.NoError(t, err)
	_, taskErr := adaptor.ParseResponse(c, response, info)
	require.Nil(t, taskErr)
	require.Len(t, captured, 1)
	form := <-captured
	assert.Equal(t, []string{"vendor-upstream"}, form.Value["model"])
	assert.Equal(t, []string{"a", "b"}, form.Value["vendor_options"])
	assert.Equal(t, []string{"false"}, form.Value["enabled"])
	assert.Equal(t, []string{"0"}, form.Value["seed"])
	assert.Equal(t, []string{"https://mock-media.invalid/source.mp4"}, form.Value["video"])
	assert.Equal(t, []string{`{"custom":true}`}, form.Value["metadata"])
	assert.Equal(t, []byte("reference image bytes"), <-referenceBytes)
}

func TestCustomVideoChannelsSubmitAndPollMockUpstream(t *testing.T) {
	service.InitHttpClient()
	for _, tc := range []struct {
		channelType    int
		model          string
		input          string
		submitPath     string
		submitResponse string
		queryPath      string
		queryResponse  string
	}{
		{58, "grok-imagine-video", `{}`, "/v1/videos/generations", `{"id":"provider-id"}`, "/v1/videos/provider-id", `{"status":"completed","video_url":"https://cdn.example/video.mp4"}`},
		{59, "runninghub-rhart-video-g", `{"images":["https://cdn.example/reference.png"],"lora1":"example.safetensors","lora1_strength_model":0}`, "/openapi/v2/rhart-video-g/image-to-video", `{"code":0,"data":{"taskId":"provider-id"}}`, "/openapi/v2/query", `{"code":0,"status":"success","url":"https://cdn.example/video.mp4"}`},
		{60, "grok-video-3-pro", `{}`, "/v1/videos", `{"id":"provider-id"}`, "/v1/videos/provider-id", `{"status":"completed","video_url":"https://cdn.example/video.mp4"}`},
		{61, "junliai-grok-imagine-video", `{}`, "/v1/videos", `{"id":"provider-id"}`, "/v1/videos/provider-id", `{"status":"completed","video_url":"https://cdn.example/video.mp4"}`},
		{62, "veo3_fast", `{"input":{"seed":0,"enableFallback":false}}`, "/api/v1/jobs/createTask", `{"code":200,"data":{"taskId":"provider-id"}}`, "/api/v1/jobs/recordInfo?taskId=provider-id", `{"code":200,"data":{"state":"success","resultJson":"{\"resultUrls\":[\"https://cdn.example/video.mp4\"]}"}}`},
		{64, "leonardo-wan-2.7", `{"dryRun":true,"seed":0,"enablePromptEnhance":false}`, "/api/video/generate", `{"id":"provider-id"}`, "/api/jobs/provider-id", `{"status":"complete","results":[{"status":"complete","urls":["https://cdn.example/video.mp4"]}]}`},
	} {
		t.Run(constant.GetChannelTypeName(tc.channelType), func(t *testing.T) {
			var submitted map[string]any
			var submitHeader http.Header
			var queryURI, queryBody, queryAuthorization string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == tc.submitPath && r.Method == http.MethodPost {
					submitHeader = r.Header.Clone()
					if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
						if err := r.ParseMultipartForm(1 << 20); err != nil {
							http.Error(w, err.Error(), 400)
							return
						}
						submitted = make(map[string]any)
						for key, values := range r.MultipartForm.Value {
							submitted[key] = values[0]
						}
					} else if err := common.DecodeJson(r.Body, &submitted); err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusAccepted)
					_, _ = io.WriteString(w, tc.submitResponse)
					return
				}
				queryURI = r.URL.RequestURI()
				queryAuthorization = r.Header.Get("Authorization")
				body, _ := io.ReadAll(r.Body)
				queryBody = string(body)
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, tc.queryResponse)
			}))
			t.Cleanup(server.Close)
			var input map[string]any
			require.NoError(t, common.UnmarshalJsonStr(tc.input, &input))
			input["model"] = "public-model"
			input["prompt"] = "original prompt"
			input["seconds"] = "5"
			raw, err := common.Marshal(input)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos", bytes.NewReader(raw))
			c.Request.Header.Set("Content-Type", "application/json")
			storage, err := common.CreateBodyStorage(raw)
			require.NoError(t, err)
			c.Set(common.KeyBodyStorage, storage)
			t.Cleanup(func() {
				if current, err := common.GetBodyStorage(c); err == nil {
					_ = current.Close()
				}
			})
			info := &relaycommon.RelayInfo{
				ChannelMeta: &relaycommon.ChannelMeta{ChannelType: tc.channelType, ChannelBaseUrl: server.URL, ApiKey: "mock-channel-key", UpstreamModelName: tc.model,
					ParamOverride: map[string]any{"prompt": "overridden prompt"}, HeadersOverride: map[string]any{"X-Compatibility": "retained"}},
				TaskRelayInfo: &relaycommon.TaskRelayInfo{PublicTaskID: "public-task"}, OriginModelName: "public-model", RequestURLPath: "/v1/videos",
			}
			adaptor := GetTaskAdaptor(constant.TaskPlatform(strconv.Itoa(tc.channelType)))
			require.NotNil(t, adaptor)
			adaptor.Init(info)
			require.Nil(t, adaptor.ValidateRequestAndSetAction(c, info))
			require.Nil(t, applyTaskParamOverrideBeforeBuild(c, info, adaptor))
			body, err := adaptor.BuildRequestBody(c, info)
			require.NoError(t, err)
			body, err = applyTaskParamOverride(c, body, info)
			require.NoError(t, err)
			response, err := adaptor.DoRequest(c, info, body)
			require.NoError(t, err)
			assert.Equal(t, http.StatusAccepted, response.StatusCode)
			parsed, taskErr := adaptor.ParseResponse(c, response, info)
			require.Nil(t, taskErr)
			require.NotNil(t, parsed)
			assert.Equal(t, "provider-id", parsed.UpstreamTaskID)
			assert.JSONEq(t, tc.submitResponse, string(parsed.TaskData), "polling must retain the provider response")
			clientBody, err := common.Marshal(parsed.ClientResponse)
			require.NoError(t, err)
			var clientResponse map[string]any
			require.NoError(t, common.Unmarshal(clientBody, &clientResponse))
			assert.Equal(t, "public-task", clientResponse["id"])
			assert.False(t, c.Writer.Written(), "parsing must not expose a task before durable host settlement")
			assert.Empty(t, recorder.Body.String())
			assert.Equal(t, "Bearer mock-channel-key", submitHeader.Get("Authorization"))
			assert.Equal(t, "retained", submitHeader.Get("X-Compatibility"))
			if tc.channelType == constant.ChannelTypeKieVideo {
				input, ok := submitted["input"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "overridden prompt", input["prompt"])
				assert.Equal(t, float64(0), input["seed"])
				assert.Equal(t, false, input["enableFallback"])
			} else {
				assert.Equal(t, "overridden prompt", submitted["prompt"])
			}
			if tc.channelType != constant.ChannelTypeRunningHub {
				assert.Equal(t, tc.model, submitted["model"])
			}
			if tc.channelType == constant.ChannelTypeRunningHub {
				assert.Equal(t, "example.safetensors", submitted["lora1"])
				assert.Equal(t, float64(0), submitted["lora1_strength_model"])
			}
			if tc.channelType == constant.ChannelTypeLeonardoAdmin {
				assert.Equal(t, false, submitted["dryRun"])
				assert.Equal(t, float64(0), submitted["seed"])
				assert.Equal(t, false, submitted["enablePromptEnhance"])
			}
			task := &model.Task{TaskID: "public-task", PrivateData: model.TaskPrivateData{UpstreamTaskID: parsed.UpstreamTaskID}}
			pollResponse, err := adaptor.FetchTask(server.URL, "saved-task-key", task, "")
			require.NoError(t, err)
			pollBody, err := io.ReadAll(pollResponse.Body)
			require.NoError(t, err)
			require.NoError(t, pollResponse.Body.Close())
			assert.Equal(t, tc.queryPath, queryURI)
			assert.Equal(t, "Bearer saved-task-key", queryAuthorization)
			if tc.channelType == constant.ChannelTypeRunningHub {
				assert.JSONEq(t, `{"taskId":"provider-id"}`, queryBody)
			}
			result, err := adaptor.ParseTaskResult(task, pollResponse, pollBody)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, model.TaskStatusSuccess, result.Status)
			assert.Equal(t, "https://cdn.example/video.mp4", result.Url)
		})
	}
}

func TestGeminiImageChannelMockUpstreamPreservesPNGAndActualImageCount(t *testing.T) {
	service.InitHttpClient()
	var jpegData bytes.Buffer
	require.NoError(t, jpeg.Encode(&jpegData, image.NewRGBA(image.Rect(0, 0, 2, 2)), nil))
	encodedJPEG := base64.StdEncoding.EncodeToString(jpegData.Bytes())
	var upstreamBody map[string]any
	var upstreamPath, upstreamKey, upstreamContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamPath, upstreamKey, upstreamContentType = r.URL.Path, r.Header.Get("x-goog-api-key"), r.Header.Get("Content-Type")
		if err := common.DecodeJson(r.Body, &upstreamBody); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		_, _ = io.WriteString(w, `{"candidates":[{"content":{"parts":[{"text":"not an image"},{"inlineData":{"mimeType":"image/jpeg","data":"`+encodedJPEG+`"}}]}}],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":20,"totalTokenCount":30}}`)
	}))
	t.Cleanup(server.Close)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	estimatedCount := 2
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{ChannelType: constant.ChannelTypeGeminiImage, ChannelBaseUrl: server.URL, ApiKey: "mock-image-key", UpstreamModelName: "gemini-3-pro-image-preview"}, RelayMode: relayconstant.RelayModeImagesGenerations, TieredBillingSnapshot: &billingexpr.BillingSnapshot{EstimatedImageCount: &estimatedCount}}
	apiType, ok := common.ChannelType2APIType(constant.ChannelTypeGeminiImage)
	require.True(t, ok)
	adaptor := GetAdaptor(apiType)
	require.NotNil(t, adaptor)
	adaptor.Init(info)
	n := uint(2)
	converted, err := adaptor.ConvertImageRequest(c, info, kitdto.ImageRequest{Prompt: "draw", N: &n, Size: "1792x1024"})
	require.NoError(t, err)
	raw, err := common.Marshal(converted)
	require.NoError(t, err)
	upstreamResponse, err := adaptor.DoRequest(c, info, bytes.NewReader(raw))
	require.NoError(t, err)
	response, ok := upstreamResponse.(*http.Response)
	require.True(t, ok)
	usage, responseErr := adaptor.DoResponse(c, response, info)
	require.Nil(t, responseErr)
	assert.NotNil(t, usage)
	assert.Equal(t, "/v1beta/models/gemini-3-pro-image-preview:generateContent", upstreamPath)
	assert.Equal(t, "mock-image-key", upstreamKey)
	assert.Equal(t, "application/json", upstreamContentType)
	config, ok := upstreamBody["generationConfig"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(2), config["candidateCount"])
	require.NotNil(t, info.BillingImageCount)
	assert.Equal(t, 1, *info.BillingImageCount, "text parts must not increase the billable image count")
	var output kitdto.ImageResponse
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &output))
	require.Len(t, output.Data, 1)
	pngData, err := base64.StdEncoding.DecodeString(output.Data[0].B64Json)
	require.NoError(t, err)
	_, err = png.Decode(bytes.NewReader(pngData))
	require.NoError(t, err, "JPEG from the provider must become a valid PNG for stored artifacts")
}

// The fixture is a sanitized, read-only export supplied for upgrade verification.
// It contains no upstream URLs or credentials and is deliberately not committed.
func TestConfiguredVideoChannelOverridesAtMockUpstream(t *testing.T) {
	fixturePath := os.Getenv("NEWAPI_CHANNEL_COMPAT_FIXTURE")
	if fixturePath == "" {
		t.Skip("set NEWAPI_CHANNEL_COMPAT_FIXTURE to a sanitized channel export")
	}
	rawFixture, err := os.ReadFile(fixturePath)
	require.NoError(t, err)
	var fixtures []struct {
		ID            int    `json:"id"`
		Type          int    `json:"type"`
		Models        string `json:"models"`
		ModelMapping  string `json:"model_mapping"`
		ParamOverride string `json:"param_override"`
	}
	require.NoError(t, common.Unmarshal(rawFixture, &fixtures))
	service.InitHttpClient()
	soraSource, err := os.ReadFile("../plugins/tasks/sora/plugin.js")
	require.NoError(t, err)
	sora, err := pluginruntime.NewRegistry().RegisterFactory(string(soraSource), pluginruntime.Options{Key: "sora"})
	require.NoError(t, err)
	tested := 0
	for _, fixture := range fixtures {
		if strings.TrimSpace(fixture.ParamOverride) == "" {
			continue
		}
		for originalModel := range strings.SplitSeq(fixture.Models, ",") {
			var mapping map[string]string
			require.NoError(t, common.UnmarshalJsonStr(fixture.ModelMapping, &mapping))
			upstreamModel := mapping[originalModel]
			imageCounts := []int{1}
			if fixture.Type == constant.ChannelTypeRunningHub && strings.Contains(upstreamModel, "rhart-video-g/") {
				imageCounts = []int{9}
			}
			if fixture.Type == constant.ChannelTypeKieVideo && strings.Contains(upstreamModel, "seedance-2-mini") {
				imageCounts = []int{1, 2, 9}
			}
			for _, imageCount := range imageCounts {
				tested++
				t.Run(strconv.Itoa(fixture.ID)+"/"+originalModel+"/"+strconv.Itoa(imageCount), func(t *testing.T) {
					var submitted map[string]any
					var requestPath, contentType string
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						requestPath, contentType = r.URL.Path, r.Header.Get("Content-Type")
						if strings.HasPrefix(contentType, "multipart/form-data") {
							if err := r.ParseMultipartForm(1 << 20); err != nil {
								http.Error(w, err.Error(), 400)
								return
							}
							submitted = make(map[string]any)
							for key, values := range r.MultipartForm.Value {
								submitted[key] = values[0]
							}
						} else if err := common.DecodeJson(r.Body, &submitted); err != nil {
							http.Error(w, err.Error(), 400)
							return
						}
						w.WriteHeader(http.StatusAccepted)
						_, _ = io.WriteString(w, `{"id":"provider-id","code":0,"data":{"taskId":"provider-id"}}`)
					}))
					t.Cleanup(server.Close)
					images := make([]string, imageCount)
					for i := range images {
						images[i] = "https://mock-media.invalid/" + strconv.Itoa(i) + ".png"
					}
					input := map[string]any{"model": originalModel, "prompt": "fixture prompt", "seconds": "5", "duration": 5, "size": "16:9", "aspect_ratio": "16:9", "images": images}
					if fixture.Type == constant.ChannelTypeLconVideo && !strings.Contains(upstreamModel, "seedance-2-0-mini") {
						delete(input, "images")
					}
					if fixture.Type == constant.ChannelTypeJunliai && upstreamModel != "minimax-h3-optimized" {
						delete(input, "images")
					}
					if fixture.Type == constant.ChannelTypeSora && strings.Contains(upstreamModel, "v2v") {
						input["video"] = "https://mock-media.invalid/source.mp4"
					}
					raw, err := common.Marshal(input)
					require.NoError(t, err)
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos", bytes.NewReader(raw))
					c.Request.Header.Set("Content-Type", "application/json")
					c.Set("model_mapping", fixture.ModelMapping)
					storage, err := common.CreateBodyStorage(raw)
					require.NoError(t, err)
					c.Set(common.KeyBodyStorage, storage)
					t.Cleanup(func() {
						if current, err := common.GetBodyStorage(c); err == nil {
							_ = current.Close()
						}
					})
					var overrides map[string]any
					require.NoError(t, common.UnmarshalJsonStr(fixture.ParamOverride, &overrides))
					info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{ChannelType: fixture.Type, ChannelBaseUrl: server.URL, ApiKey: "mock-fixture-key", UpstreamModelName: originalModel, ParamOverride: overrides}, TaskRelayInfo: &relaycommon.TaskRelayInfo{PublicTaskID: "public-task"}, OriginModelName: originalModel, RequestURLPath: "/v1/videos"}
					require.NoError(t, helper.ModelMappedHelper(c, info, nil))
					platform := constant.TaskPlatform(strconv.Itoa(fixture.Type))
					var adaptor channel.TaskAdaptor
					if fixture.Type == constant.ChannelTypeSora {
						platform = "sora"
						adaptor = jspluginadaptor.New(sora)
					} else {
						adaptor = GetTaskAdaptor(platform)
					}
					require.NotNil(t, adaptor)
					adaptor.Init(info)
					if bodyAdaptor, ok := adaptor.(channel.TaskRequestBodyOverrideAdaptor); ok {
						info.SupportsRequestBodyOverride = bodyAdaptor.SupportsRequestBodyOverride()
					}
					require.Nil(t, prepareLegacySoraTaskRequest(c, info, platform))
					require.Nil(t, adaptor.ValidateRequestAndSetAction(c, info))
					require.Nil(t, applyTaskParamOverrideBeforeBuild(c, info, adaptor))
					body, err := adaptor.BuildRequestBody(c, info)
					require.NoError(t, err)
					body, err = applyTaskParamOverride(c, body, info)
					require.NoError(t, err)
					response, err := adaptor.DoRequest(c, info, body)
					require.NoError(t, err)
					_, taskErr := adaptor.ParseResponse(c, response, info)
					require.Nil(t, taskErr)
					if fixture.Type != constant.ChannelTypeKieVideo {
						assert.Equal(t, "fixture prompt", submitted["prompt"])
					}
					if fixture.Type != constant.ChannelTypeRunningHub {
						assert.Equal(t, upstreamModel, submitted["model"])
					}
					switch fixture.Type {
					case constant.ChannelTypeRunningHub:
						assert.Equal(t, strings.TrimPrefix(upstreamModel, "path:"), requestPath)
						if strings.Contains(upstreamModel, "rhart-video-g/") {
							assert.Equal(t, float64(6), submitted["duration"])
							assert.Equal(t, []any{images[0], images[2], images[3], images[4], images[5], images[6], images[7]}, submitted["imageUrls"])
						}
						if strings.Contains(upstreamModel, "wan-2.7/") {
							assert.Equal(t, float64(3), submitted["duration"])
							assert.Equal(t, images[0], submitted["firstImageUrl"])
							assert.Equal(t, images[0], submitted["lastImageUrl"])
							assert.NotContains(t, submitted, "imageUrls")
						}
					case constant.ChannelTypeKieVideo:
						fields, ok := submitted["input"].(map[string]any)
						require.True(t, ok)
						assert.Equal(t, "fixture prompt", fields["prompt"])
						if strings.Contains(upstreamModel, "seedance-2-mini") {
							assert.Equal(t, float64(4), fields["duration"])
							assert.Equal(t, false, fields["web_search"])
							assert.Equal(t, true, fields["generate_audio"])
							assert.Equal(t, true, fields["nsfw_checker"])
							if imageCount == 1 {
								assert.Equal(t, images[0], fields["first_frame_url"])
								assert.Equal(t, images[0], fields["last_frame_url"])
							} else {
								references, ok := fields["reference_image_urls"].([]any)
								require.True(t, ok)
								require.Len(t, references, imageCount)
								assert.Equal(t, images[0], references[0])
								if imageCount == 9 {
									assert.Equal(t, images[0], references[1])
									assert.Equal(t, images[0], references[8])
								}
							}
						}
					case constant.ChannelTypeLconVideo:
						if strings.Contains(upstreamModel, "seedance-2-0-mini") {
							assert.Equal(t, "application/json", contentType)
							assert.Equal(t, float64(4), submitted["duration"])
							assert.Equal(t, "16:9", submitted["ratio"])
							assert.Equal(t, []any{images[0], images[0]}, submitted["images"])
							assert.NotContains(t, submitted, "seconds")
							assert.NotContains(t, submitted, "size")
						}
					case constant.ChannelTypeJunliai:
						if upstreamModel == "minimax-h3-optimized" {
							assert.Equal(t, float64(5), submitted["duration"])
							assert.Equal(t, "768p", submitted["resolution"])
							assert.Equal(t, images[0], submitted["start_frame"])
							assert.Equal(t, images[0], submitted["end_frame"])
							assert.NotContains(t, submitted, "reference_images")
						} else {
							assert.Equal(t, float64(6), submitted["duration"])
						}
					case constant.ChannelTypeLeonardoAdmin:
						if upstreamModel == "wan-2.7" {
							fields, ok := submitted["parameters"].(map[string]any)
							require.True(t, ok)
							assert.Equal(t, float64(5), fields["duration"])
							assert.Equal(t, float64(960), fields["width"])
							assert.Equal(t, float64(960), fields["height"])
							assert.Equal(t, "720p", fields["resolution"])
							assert.Equal(t, []any{"StartFrame", "EndFrame"}, submitted["reference_frame_types"])
						}
					case constant.ChannelTypeSora:
						if strings.HasPrefix(upstreamModel, "grok-imagine-video") {
							assert.Equal(t, "6", submitted["seconds"])
							assert.Equal(t, []any{images[0]}, submitted["reference_images"])
							assert.NotContains(t, submitted, "images")
						}
						if strings.Contains(upstreamModel, "doubao-seedance") {
							assert.Equal(t, "5", submitted["seconds"])
							assert.Equal(t, "720p", submitted["resolution"])
							assert.Equal(t, []any{images[0], images[0]}, submitted["images"])
						}
						if strings.Contains(upstreamModel, "v2v") {
							assert.Equal(t, "https://mock-media.invalid/source.mp4", submitted["video"])
						}
					}
				})
			}
		}
	}
	require.Positive(t, tested, "fixture must contain at least one configured video override")
}

func TestLconConfiguredBodyFormatPreservesOverridesAtMockUpstream(t *testing.T) {
	service.InitHttpClient()
	for _, format := range []string{"json", "multipart"} {
		t.Run(format, func(t *testing.T) {
			var actual map[string]any
			var contentType string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				contentType = r.Header.Get("Content-Type")
				if format == "multipart" {
					if err := r.ParseMultipartForm(1 << 20); err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
					actual = make(map[string]any)
					for key, values := range r.MultipartForm.Value {
						actual[key] = values[0]
					}
				} else if err := common.DecodeJson(r.Body, &actual); err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
				w.WriteHeader(http.StatusAccepted)
				_, _ = io.WriteString(w, `{"id":"provider-id"}`)
			}))
			t.Cleanup(server.Close)
			raw := []byte(`{"model":"public-model","prompt":"draw","seconds":"5","seed":"0","enabled":"false","removed":"must disappear","nullable":null}`)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos", bytes.NewReader(raw))
			c.Request.Header.Set("Content-Type", "application/json")
			storage, err := common.CreateBodyStorage(raw)
			require.NoError(t, err)
			c.Set(common.KeyBodyStorage, storage)
			t.Cleanup(func() {
				if current, err := common.GetBodyStorage(c); err == nil {
					_ = current.Close()
				}
			})
			info := &relaycommon.RelayInfo{
				ChannelMeta: &relaycommon.ChannelMeta{ChannelType: constant.ChannelTypeLconVideo, ChannelBaseUrl: server.URL, ApiKey: "mock-key", UpstreamModelName: "grok-video-3-pro", ParamOverride: map[string]any{"operations": []any{
					map[string]any{"mode": "delete", "path": "removed"},
					map[string]any{"mode": "set_request_body", "value": map[string]any{"format": format, "field_types": map[string]any{"seed": "integer", "enabled": "boolean"}}},
				}}},
				TaskRelayInfo: &relaycommon.TaskRelayInfo{PublicTaskID: "public-task"}, SupportsRequestBodyOverride: true, OriginModelName: "public-model",
			}
			adaptor := GetTaskAdaptor(constant.TaskPlatform(strconv.Itoa(constant.ChannelTypeLconVideo)))
			require.NotNil(t, adaptor)
			adaptor.Init(info)
			require.Nil(t, adaptor.ValidateRequestAndSetAction(c, info))
			require.Nil(t, applyTaskParamOverrideBeforeBuild(c, info, adaptor))
			body, err := adaptor.BuildRequestBody(c, info)
			require.NoError(t, err)
			body, err = applyTaskParamOverride(c, body, info)
			require.NoError(t, err)
			response, err := adaptor.DoRequest(c, info, body)
			require.NoError(t, err)
			_, taskErr := adaptor.ParseResponse(c, response, info)
			require.Nil(t, taskErr)
			assert.Equal(t, "grok-video-3-pro", actual["model"])
			assert.NotContains(t, actual, "removed")
			assert.Contains(t, actual, "nullable")
			if format == "json" {
				assert.Equal(t, "application/json", contentType)
				assert.Equal(t, float64(0), actual["seed"])
				assert.Equal(t, false, actual["enabled"])
				assert.Nil(t, actual["nullable"])
			} else {
				assert.True(t, strings.HasPrefix(contentType, "multipart/form-data; boundary="))
				assert.Equal(t, "0", actual["seed"])
				assert.Equal(t, "false", actual["enabled"])
				assert.Equal(t, "null", actual["nullable"])
			}
		})
	}
}

func TestCustomTaskPollingDoesNotTreatUnknownStatusAsProgress(t *testing.T) {
	for _, channelType := range []int{
		constant.ChannelTypeGrokVideo, constant.ChannelTypeJunliai,
		constant.ChannelTypeKieVideo, constant.ChannelTypeRunningHub,
		constant.ChannelTypeLconVideo, constant.ChannelTypeLeonardoAdmin,
	} {
		t.Run(constant.GetChannelTypeName(channelType), func(t *testing.T) {
			adaptor := GetTaskAdaptor(constant.TaskPlatform(strconv.Itoa(channelType)))
			require.NotNil(t, adaptor)
			for _, body := range []string{`{}`, `{"status":"unrecognized","code":0}`} {
				result, err := adaptor.ParseTaskResult(&model.Task{}, nil, []byte(body))
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, model.TaskStatusUnknown, result.Status, body)
			}
		})
	}
}

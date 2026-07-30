package runninghub

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestBuildRequestURLUsesDefaultImageToVideoPath(t *testing.T) {
	adaptor := &TaskAdaptor{}
	adaptor.Init(&relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelBaseUrl: "https://www.runninghub.ai/",
		},
	})

	url, err := adaptor.BuildRequestURL(&relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "runninghub-rhart-video-g",
		},
	})

	require.NoError(t, err)
	require.Equal(t, "https://www.runninghub.ai/openapi/v2/rhart-video-g/image-to-video", url)
}

func TestBuildRequestURLUsesMappedEndpointPath(t *testing.T) {
	adaptor := &TaskAdaptor{}
	adaptor.Init(&relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelBaseUrl: "https://www.runninghub.ai",
		},
	})

	url, err := adaptor.BuildRequestURL(&relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "path:/openapi/v2/rhart-video/ltx-2.3/image-to-video-lora",
		},
		OriginModelName: "runninghub-ltx-2.3",
	})

	require.NoError(t, err)
	require.Equal(t, "https://www.runninghub.ai/openapi/v2/rhart-video/ltx-2.3/image-to-video-lora", url)
}

func TestMappedEndpointPathRejectsInvalidPaths(t *testing.T) {
	testCases := []string{
		"path:https://www.runninghub.ai/openapi/v2/rhart-video/ltx-2.3/image-to-video-lora",
		"path:/openapi/v2/../image-to-video",
		"path:/v2/rhart-video/ltx-2.3/image-to-video-lora",
		"path:/openapi/v2/rhart-video/ltx-2.3/image-to-video-lora?debug=1",
		"path:/openapi/v2/rhart-video/ltx-2.3/image to video",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			_, _, err := mappedEndpointPath(tc)
			require.Error(t, err)
		})
	}
}

func TestBuildRequestBodyPreservesOfficialLoraFieldsAndImageURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	storage, err := common.CreateBodyStorage([]byte(`{
		"imageUrl": "https://www.runninghub.ai/view?filename=input.png",
		"prompt": "make a video",
		"resolution": "720p",
		"aspectRatio": "9:16",
		"duration": 7,
		"lora1": "framee_4000.safetensors",
		"lora1_strength_model": 0.5,
		"lora2": "other.safetensors",
		"lora2_strength_model": 0,
		"unknown": "drop me"
	}`))
	require.NoError(t, err)
	c.Set(common.KeyBodyStorage, storage)

	adaptor := &TaskAdaptor{}
	body, err := adaptor.BuildRequestBody(c, nil)
	require.NoError(t, err)
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, common.Unmarshal(data, &payload))
	require.Equal(t, "https://www.runninghub.ai/view?filename=input.png", payload["imageUrl"])
	require.Equal(t, []any{"https://www.runninghub.ai/view?filename=input.png"}, payload["imageUrls"])
	require.Equal(t, "framee_4000.safetensors", payload["lora1"])
	require.Equal(t, 0.5, payload["lora1_strength_model"])
	require.Equal(t, "other.safetensors", payload["lora2"])
	require.Equal(t, float64(0), payload["lora2_strength_model"])
	require.NotContains(t, payload, "unknown")
}

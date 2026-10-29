package plugins_test

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/pkg/jsplugin"
	builtinplugins "github.com/QuantumNous/new-api/plugins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSoraResponsesProtocol(t *testing.T) {
	testVideoResponsesProtocol(t, videoResponsesTestCase{
		pluginKey: "sora",
		model:     "sora-2-pro",
		requestBody: map[string]any{
			"model":   "sora-2-pro",
			"input":   "waves at sunset",
			"seconds": 8,
			"size":    "1792x1024",
		},
		wantAction: "text_to_video",
		wantRequest: map[string]any{
			"model":   "sora-2-pro",
			"prompt":  "waves at sunset",
			"seconds": float64(8),
			"size":    "1792x1024",
		},
		wantUsageKeys:  []string{"seconds", "size"},
		wantVendorName: "sora",
	})
}

func TestSoraContentPrefersCompletedDataURL(t *testing.T) {
	source, err := builtinplugins.Source("sora")
	require.NoError(t, err)
	plugin, err := jsplugin.NewRegistry().RegisterFactory(source, jsplugin.Options{Key: "sora"})
	require.NoError(t, err)
	value, err := plugin.Engine.Call(t.Context(), "buildContentRequest", map[string]any{
		"artifactKey":    "video",
		"baseUrl":        "https://upstream.example",
		"apiKey":         "channel-secret",
		"upstreamTaskId": "provider-id",
		"clientRequest":  map[string]any{"method": "GET"},
		"data": map[string]any{
			"video_url": "/v1/videos/provider-id/content",
			"data":      []any{map[string]any{"url": "https://cdn.example/result.mp4"}},
		},
	})
	require.NoError(t, err)
	encoded, err := common.Marshal(value)
	require.NoError(t, err)
	var descriptor map[string]any
	require.NoError(t, common.Unmarshal(encoded, &descriptor))
	assert.Equal(t, "https://cdn.example/result.mp4", descriptor["url"])
	assert.Equal(t, true, descriptor["credentialless"])
	assert.NotContains(t, descriptor, "headers")
}

func TestSoraLegacyBillingKeepsPerRequestPrice(t *testing.T) {
	source, err := builtinplugins.Source("sora")
	require.NoError(t, err)
	plugin, err := jsplugin.NewRegistry().RegisterFactory(source, jsplugin.Options{Key: "sora"})
	require.NoError(t, err)
	value, err := plugin.Engine.Call(t.Context(), "extractUsage", map[string]any{
		"usagePurpose": "billing_ratios",
		"requestBody":  map[string]any{"seconds": 8, "size": "1792x1024"},
	})
	require.NoError(t, err)
	encoded, err := common.Marshal(value)
	require.NoError(t, err)
	var ratios map[string]any
	require.NoError(t, common.Unmarshal(encoded, &ratios))
	assert.Empty(t, ratios)
}

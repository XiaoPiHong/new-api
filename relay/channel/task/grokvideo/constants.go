package grokvideo

var ModelList = []string{
	"grok-imagine-video",
}

var ChannelName = "grok-video"

const (
	VideoGenerationEndpoint = "/v1/videos/generations"
	VideoEditEndpoint       = "/v1/videos/edits"
	QueryTaskEndpoint       = "/v1/videos"
)

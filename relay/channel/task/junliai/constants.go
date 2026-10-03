package junliai

var ModelList = []string{
	"junliai-grok-imagine-video",
}

var ChannelName = "junliai"

const (
	// Junliai 渠道统一使用通用视频接口，不影响独立的 Grok 视频渠道。
	VideoGenerationEndpoint = "/v1/videos"
	QueryTaskEndpoint       = "/v1/videos"
)

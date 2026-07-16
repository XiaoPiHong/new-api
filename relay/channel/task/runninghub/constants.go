package runninghub

const ModelNamePrefix = "runninghub-"

var ModelList = []string{
	"runninghub-rhart-video-g",
	"runninghub-rhart-video-v3.1-fast",
	"runninghub-rhart-video/sparkvideo-2.0-fast",
}

var ChannelName = "runninghub"

const (
	ImageToVideoEndpointFormat = "/openapi/v2/%s/image-to-video"
	QueryTaskEndpoint          = "/openapi/v2/query"
)

package leonardo

// These are compatibility suggestions shown by the channel model catalog.
// They are not a whitelist: the Leonardo adapter forwards any model name
// configured on a channel (after model mapping) to Leonardo Admin. This keeps
// the channel usable when Leonardo publishes a new image model before the
// static suggestions are updated.
var ImageModelList = []string{
	"leonardo-nano-banana-2",
	"leonardo-gpt-image-2",
}

var VideoModelList = []string{
	"leonardo-motion-2.0",
	"leonardo-hailuo-2_3",
}

var ModelList = append(append([]string{}, ImageModelList...), VideoModelList...)

const ChannelName = "Leonardo Admin"

package leonardo

// Model names exposed by new-api are management aliases. Configure a channel
// model mapping to translate these names to the IDs understood by Leonardo
// Admin (for example leonardo-nano-banana-2 -> nano-banana-2).
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

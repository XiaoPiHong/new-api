package sora

import (
	"testing"

	"github.com/QuantumNous/new-api/model"
)

func TestParseTaskResultPrefersCompletedDataItemURL(t *testing.T) {
	adaptor := &TaskAdaptor{}
	taskInfo, err := adaptor.ParseTaskResult([]byte(`{
		"status": "completed",
		"video_url": "/v1/videos/vid-391fd2436fc8/content",
		"data": [
			{
				"url": "https://download.oaibox.xyz/v1/videos/task_TLuVDxz8dHfzn7Nl0Q2ACjHIl6vWLf5K/content"
			}
		]
	}`))
	if err != nil {
		t.Fatalf("ParseTaskResult returned error: %v", err)
	}
	if taskInfo.Status != model.TaskStatusSuccess {
		t.Fatalf("status = %q, want %q", taskInfo.Status, model.TaskStatusSuccess)
	}
	want := "https://download.oaibox.xyz/v1/videos/task_TLuVDxz8dHfzn7Nl0Q2ACjHIl6vWLf5K/content"
	if taskInfo.Url != want {
		t.Fatalf("url = %q, want %q", taskInfo.Url, want)
	}
}

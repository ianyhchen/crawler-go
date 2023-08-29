package client

import (
	"crawler_go/pkg/ptt"
	"testing"
)

func TestGetLatestBoardData(t *testing.T) {
	url := ptt.GetBoardURL("gossip")
	data, err := GetLatestBoardData(url)
	if err != nil {
		t.Fatalf("Error retrieving board data: %v", err)
	}
	if len(data) == 0 {
		t.Logf("No new data found")
	}
	t.Logf("Retrieved %d topics", len(data))
}
func TestGetTopicContent(t *testing.T) {
	url := "https://www.ptt.cc/bbs/Gossiping/M.1692775209.A.DE7.html"
	data, err := GetTopicContent(url)
	if err != nil {
		t.Fatalf("Error retrieving topic content: %v", err)
	}
	t.Logf("Retrieved topic content:%+v\n", data)
}

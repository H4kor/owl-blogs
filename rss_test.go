package owl_test

import (
	"h4kor/owl-blogs"
	"strings"
	"testing"
)

func TestRenderRSSFeedMeta(t *testing.T) {
	user := getTestUser()
	user.SetConfig(owl.UserConfig{
		Title:    "Test Title",
		SubTitle: "Test SubTitle",
	})
	res, err := owl.RenderRSSFeed(user)
	if err != nil {
		t.Error("Error rendering RSS feed: " + err.Error())
		return
	}
	if !strings.Contains(res, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("xml version not rendered. Got: " + res)
	}
	if !strings.Contains(res, "<rss version=\"2.0\">") {
		t.Error("rss version not rendered. Got: " + res)
	}

}

func TestRenderRSSFeedUserData(t *testing.T) {
	user := getTestUser()
	user.SetConfig(owl.UserConfig{
		Title:    "Test Title",
		SubTitle: "Test SubTitle",
	})
	res, err := owl.RenderRSSFeed(user)
	if err != nil {
		t.Error("Error rendering RSS feed: " + err.Error())
		return
	}
	if !strings.Contains(res, "Test Title") {
		t.Error("Title not rendered. Got: " + res)
	}
	if !strings.Contains(res, "Test SubTitle") {
		t.Error("SubTitle not rendered. Got: " + res)
	}
	if !strings.Contains(res, "http://localhost:8080/user/") {
		t.Error("SubTitle not rendered. Got: " + res)
	}
}

package owl_test

import (
	"h4kor/owl-blogs"
	"os"
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

func TestRenderRSSFeedPostData(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)

	content := "---\n"
	content += "title: Test Post\n"
	content += "date: 2015-01-01\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	res, err := owl.RenderRSSFeed(user)
	if err != nil {
		t.Error("Error rendering RSS feed: " + err.Error())
		return
	}
	if !strings.Contains(res, "Test Post") {
		t.Error("Title not rendered. Got: " + res)
	}
	if !strings.Contains(res, post.FullUrl()) {
		t.Error("SubTitle not rendered. Got: " + res)
	}
	if !strings.Contains(res, "Thu, 01 Jan 2015 00:00:00 +0000") {
		t.Error("Date not rendered. Got: " + res)
	}
}

func TestRenderRSSFeedPostDataWithoutDate(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)

	content := "---\n"
	content += "title: Test Post\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	res, err := owl.RenderRSSFeed(user)
	if err != nil {
		t.Error("Error rendering RSS feed: " + err.Error())
		return
	}
	if !strings.Contains(res, "Test Post") {
		t.Error("Title not rendered. Got: " + res)
	}
	if !strings.Contains(res, post.FullUrl()) {
		t.Error("SubTitle not rendered. Got: " + res)
	}
}

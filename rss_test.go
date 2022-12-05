package owl_test

import (
	"h4kor/owl-blogs"
	"h4kor/owl-blogs/test/assertions"
	"os"
	"testing"
)

func TestRenderRSSFeedMeta(t *testing.T) {
	user := getTestUser()
	user.SetConfig(owl.UserConfig{
		Title:    "Test Title",
		SubTitle: "Test SubTitle",
	})
	res, err := owl.RenderRSSFeed(user)
	assertions.AssertNoError(t, err, "Error rendering RSS feed")
	assertions.AssertContains(t, res, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	assertions.AssertContains(t, res, "<rss version=\"2.0\">")

}

func TestRenderRSSFeedUserData(t *testing.T) {
	user := getTestUser()
	user.SetConfig(owl.UserConfig{
		Title:    "Test Title",
		SubTitle: "Test SubTitle",
	})
	res, err := owl.RenderRSSFeed(user)
	assertions.AssertNoError(t, err, "Error rendering RSS feed")
	assertions.AssertContains(t, res, "Test Title")
	assertions.AssertContains(t, res, "Test SubTitle")
	assertions.AssertContains(t, res, "http://localhost:8080/user/")
}

func TestRenderRSSFeedPostData(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "testpost"}, "")

	content := "---\n"
	content += "title: Test Post\n"
	content += "date: 2015-01-01\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	res, err := owl.RenderRSSFeed(user)
	assertions.AssertNoError(t, err, "Error rendering RSS feed")
	assertions.AssertContains(t, res, "Test Post")
	assertions.AssertContains(t, res, post.FullUrl())
	assertions.AssertContains(t, res, "Thu, 01 Jan 2015 00:00:00 +0000")
}

func TestRenderRSSFeedPostDataWithoutDate(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "testpost"}, "")

	content := "---\n"
	content += "title: Test Post\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	res, err := owl.RenderRSSFeed(user)
	assertions.AssertNoError(t, err, "Error rendering RSS feed")
	assertions.AssertContains(t, res, "Test Post")
	assertions.AssertContains(t, res, post.FullUrl())
}

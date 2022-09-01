package owl_test

import (
	"h4kor/owl-blogs"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestCanRenderPost(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result, err := owl.RenderPost(&post)

	if err != nil {
		t.Error("Error rendering post: " + err.Error())
		return
	}

	if !strings.Contains(result, "<h1 class=\"p-name\">testpost</h1>") {
		t.Error("Post title not rendered as h1. Got: " + result)
	}

}

func TestRenderPostHEntry(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result, _ := owl.RenderPost(&post)
	if !strings.Contains(result, "class=\"h-entry\"") {
		t.Error("h-entry container not rendered. Got: " + result)
	}
	if !strings.Contains(result, "class=\"p-name\"") {
		t.Error("p-name not rendered. Got: " + result)
	}
	if !strings.Contains(result, "class=\"e-content\"") {
		t.Error("e-content not rendered. Got: " + result)
	}

}

func TestRendererUsesBaseTemplate(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result, err := owl.RenderPost(&post)

	if err != nil {
		t.Error("Error rendering post: " + err.Error())
		return
	}

	if !strings.Contains(result, "<html") {
		t.Error("Base template not used. Got: " + result)
	}
}

func TestCanRenderIndexPage(t *testing.T) {
	user := getTestUser()
	user.CreateNewPost("testpost1")
	user.CreateNewPost("testpost2")
	result, _ := owl.RenderIndexPage(user)
	if !strings.Contains(result, "testpost1") {
		t.Error("Post title not rendered. Got: " + result)
	}
	if !strings.Contains(result, "testpost2") {
		t.Error("Post title not rendered. Got: " + result)
	}
}

func TestIndexPageContainsHFeedContainer(t *testing.T) {
	user := getTestUser()
	user.CreateNewPost("testpost1")

	result, _ := owl.RenderIndexPage(user)
	if !strings.Contains(result, "<div class=\"h-feed\">") {
		t.Error("h-feed container not rendered. Got: " + result)
	}
}

func TestIndexPageContainsHEntryAndUUrl(t *testing.T) {
	user := getTestUser()
	user.CreateNewPost("testpost1")

	result, _ := owl.RenderIndexPage(user)
	if !strings.Contains(result, "class=\"h-entry\"") {
		t.Error("h-entry container not rendered. Got: " + result)
	}
	if !strings.Contains(result, "class=\"u-url\"") {
		t.Error("u-url not rendered. Got: " + result)
	}

}

func TestRenderIndexPageWithBrokenBaseTemplate(t *testing.T) {
	user := getTestUser()
	user.CreateNewPost("testpost1")
	user.CreateNewPost("testpost2")

	os.WriteFile(path.Join(user.Dir(), "meta/base.html"), []byte("{{content}}"), 0644)

	_, err := owl.RenderIndexPage(user)
	if err == nil {
		t.Error("Expected error rendering index page, got nil")
	}
}

func TestRenderUserList(t *testing.T) {
	repo := getTestRepo()
	repo.CreateUser("user1")
	repo.CreateUser("user2")

	result, err := owl.RenderUserList(repo)
	if err != nil {
		t.Error("Error rendering user list: " + err.Error())
	}

	if !strings.Contains(result, "user1") {
		t.Error("Post title not rendered. Got: " + result)
	}
	if !strings.Contains(result, "user2") {
		t.Error("Post title not rendered. Got: " + result)
	}
}

func TestRendersHeaderTitle(t *testing.T) {
	user := getTestUser()
	user.SetConfig(owl.UserConfig{
		Title:       "Test Title",
		SubTitle:    "Test SubTitle",
		HeaderColor: "#ff1337",
	})
	post, _ := user.CreateNewPost("testpost")

	result, _ := owl.RenderPost(&post)
	if !strings.Contains(result, "Test Title") {
		t.Error("Header title not rendered. Got: " + result)
	}
	if !strings.Contains(result, "Test SubTitle") {
		t.Error("Header subtitle not rendered. Got: " + result)
	}
	if !strings.Contains(result, "#ff1337") {
		t.Error("Header color not rendered. Got: " + result)
	}
}

func TestRenderPostIncludesRelToWebMention(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")

	result, _ := owl.RenderPost(&post)
	if !strings.Contains(result, "rel=\"webmention\"") {
		t.Error("webmention rel not rendered. Got: " + result)
	}

	if !strings.Contains(result, "href=\""+user.WebmentionUrl()+"\"") {
		t.Error("webmention href not rendered. Got: " + result)
	}
}

func TestRenderPostAddsLinksToApprovedWebmention(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	webmention := owl.Webmention{
		Source:         "http://example.com/source3",
		Title:          "Test Title",
		ApprovalStatus: "approved",
		RetrievedAt:    time.Now().Add(time.Hour * -2),
	}
	post.PersistWebmention(webmention)
	webmention = owl.Webmention{
		Source:         "http://example.com/source4",
		ApprovalStatus: "rejected",
		RetrievedAt:    time.Now().Add(time.Hour * -3),
	}
	post.PersistWebmention(webmention)

	result, _ := owl.RenderPost(&post)
	if !strings.Contains(result, "http://example.com/source3") {
		t.Error("webmention not rendered. Got: " + result)
	}
	if !strings.Contains(result, "Test Title") {
		t.Error("webmention title not rendered. Got: " + result)
	}
	if strings.Contains(result, "http://example.com/source4") {
		t.Error("unapproved webmention rendered. Got: " + result)
	}

}

func TestRenderPostNotMentioningWebmentionsIfNoAvail(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result, _ := owl.RenderPost(&post)

	if strings.Contains(result, "Webmention") {
		t.Error("Webmention mentioned. Got: " + result)
	}

}

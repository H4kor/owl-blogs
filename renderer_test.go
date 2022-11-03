package owl_test

import (
	"h4kor/owl-blogs"
	"h4kor/owl-blogs/priv/assertions"
	"os"
	"path"
	"testing"
	"time"
)

func TestCanRenderPost(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	result, err := owl.RenderPost(post)

	assertions.AssertNoError(t, err, "Error rendering post")
	assertions.AssertContains(t, result, "testpost")

}

func TestRenderOneMe(t *testing.T) {
	user := getTestUser()
	config := user.Config()
	config.Me = append(config.Me, owl.UserMe{
		Name: "Twitter",
		Url:  "https://twitter.com/testhandle",
	})

	user.SetConfig(config)
	post, _ := user.CreateNewPost("testpost", false)
	result, err := owl.RenderPost(post)

	assertions.AssertNoError(t, err, "Error rendering post")
	assertions.AssertContains(t, result, "href=\"https://twitter.com/testhandle\" rel=\"me\"")

}

func TestRenderTwoMe(t *testing.T) {
	user := getTestUser()
	config := user.Config()
	config.Me = append(config.Me, owl.UserMe{
		Name: "Twitter",
		Url:  "https://twitter.com/testhandle",
	})
	config.Me = append(config.Me, owl.UserMe{
		Name: "Github",
		Url:  "https://github.com/testhandle",
	})

	user.SetConfig(config)
	post, _ := user.CreateNewPost("testpost", false)
	result, err := owl.RenderPost(post)

	assertions.AssertNoError(t, err, "Error rendering post")
	assertions.AssertContains(t, result, "href=\"https://twitter.com/testhandle\" rel=\"me\"")
	assertions.AssertContains(t, result, "href=\"https://github.com/testhandle\" rel=\"me\"")

}

func TestRenderPostHEntry(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	result, _ := owl.RenderPost(post)
	assertions.AssertContains(t, result, "class=\"h-entry\"")
	assertions.AssertContains(t, result, "class=\"p-name\"")
	assertions.AssertContains(t, result, "class=\"e-content\"")

}

func TestRendererUsesBaseTemplate(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	result, err := owl.RenderPost(post)

	assertions.AssertNoError(t, err, "Error rendering post")
	assertions.AssertContains(t, result, "<html")
}

func TestCanRenderIndexPage(t *testing.T) {
	user := getTestUser()
	user.CreateNewPost("testpost1", false)
	user.CreateNewPost("testpost2", false)
	result, _ := owl.RenderIndexPage(user)
	assertions.AssertContains(t, result, "testpost1")
	assertions.AssertContains(t, result, "testpost2")
}

func TestIndexPageContainsHFeedContainer(t *testing.T) {
	user := getTestUser()
	user.CreateNewPost("testpost1", false)

	result, _ := owl.RenderIndexPage(user)
	assertions.AssertContains(t, result, "<div class=\"h-feed\">")
}

func TestIndexPageContainsHEntryAndUUrl(t *testing.T) {
	user := getTestUser()
	user.CreateNewPost("testpost1", false)

	result, _ := owl.RenderIndexPage(user)
	assertions.AssertContains(t, result, "class=\"h-entry\"")
	assertions.AssertContains(t, result, "class=\"u-url\"")

}

func TestRenderIndexPageWithBrokenBaseTemplate(t *testing.T) {
	user := getTestUser()
	user.CreateNewPost("testpost1", false)
	user.CreateNewPost("testpost2", false)

	os.WriteFile(path.Join(user.Dir(), "meta/base.html"), []byte("{{content}}"), 0644)

	_, err := owl.RenderIndexPage(user)
	assertions.AssertError(t, err, "Expected error rendering index page")
}

func TestRenderUserList(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.CreateUser("user1")
	repo.CreateUser("user2")

	result, err := owl.RenderUserList(repo)
	assertions.AssertNoError(t, err, "Error rendering user list")
	assertions.AssertContains(t, result, "user1")
	assertions.AssertContains(t, result, "user2")
}

func TestRendersHeaderTitle(t *testing.T) {
	user := getTestUser()
	user.SetConfig(owl.UserConfig{
		Title:       "Test Title",
		SubTitle:    "Test SubTitle",
		HeaderColor: "#ff1337",
	})
	post, _ := user.CreateNewPost("testpost", false)

	result, _ := owl.RenderPost(post)
	assertions.AssertContains(t, result, "Test Title")
	assertions.AssertContains(t, result, "Test SubTitle")
	assertions.AssertContains(t, result, "#ff1337")
}

func TestRenderPostIncludesRelToWebMention(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)

	result, _ := owl.RenderPost(post)
	assertions.AssertContains(t, result, "rel=\"webmention\"")

	assertions.AssertContains(t, result, "href=\""+user.WebmentionUrl()+"\"")
}

func TestRenderPostAddsLinksToApprovedWebmention(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	webmention := owl.WebmentionIn{
		Source:         "http://example.com/source3",
		Title:          "Test Title",
		ApprovalStatus: "approved",
		RetrievedAt:    time.Now().Add(time.Hour * -2),
	}
	post.PersistIncomingWebmention(webmention)
	webmention = owl.WebmentionIn{
		Source:         "http://example.com/source4",
		ApprovalStatus: "rejected",
		RetrievedAt:    time.Now().Add(time.Hour * -3),
	}
	post.PersistIncomingWebmention(webmention)

	result, _ := owl.RenderPost(post)
	assertions.AssertContains(t, result, "http://example.com/source3")
	assertions.AssertContains(t, result, "Test Title")
	assertions.AssertNotContains(t, result, "http://example.com/source4")

}

func TestRenderPostNotMentioningWebmentionsIfNoAvail(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	result, _ := owl.RenderPost(post)

	assertions.AssertNotContains(t, result, "Webmention")

}

func TestRenderIncludesFullUrl(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	result, _ := owl.RenderPost(post)

	assertions.AssertContains(t, result, "class=\"u-url\"")
	assertions.AssertContains(t, result, post.FullUrl())
}

func TestAddAvatarIfExist(t *testing.T) {
	user := getTestUser()
	os.WriteFile(path.Join(user.MediaDir(), "avatar.png"), []byte("test"), 0644)

	result, _ := owl.RenderIndexPage(user)
	assertions.AssertContains(t, result, "avatar.png")
}

func TestAuthorNameInPost(t *testing.T) {
	user := getTestUser()
	user.SetConfig(owl.UserConfig{
		Title:       "Test Title",
		SubTitle:    "Test SubTitle",
		HeaderColor: "#ff1337",
		AuthorName:  "Test Author",
	})
	post, _ := user.CreateNewPost("testpost", false)

	result, _ := owl.RenderPost(post)
	assertions.AssertContains(t, result, "Test Author")
}

func TestRenderReplyWithoutText(t *testing.T) {

	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)

	content := "---\n"
	content += "title: test\n"
	content += "date: Wed, 17 Aug 2022 10:50:02 +0000\n"
	content += "reply: \n"
	content += "  url: https://example.com/post\n"
	content += "---\n"
	content += "\n"
	content += "Hi \n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	result, _ := owl.RenderPost(post)
	assertions.AssertContains(t, result, "https://example.com/post")
}

func TestRenderReplyWithText(t *testing.T) {

	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)

	content := "---\n"
	content += "title: test\n"
	content += "date: Wed, 17 Aug 2022 10:50:02 +0000\n"
	content += "reply: \n"
	content += "  url: https://example.com/post\n"
	content += "  text: \"This is a reply\"\n"
	content += "---\n"
	content += "Hi \n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	result, _ := owl.RenderPost(post)
	assertions.AssertContains(t, result, "https://example.com/post")

	assertions.AssertContains(t, result, "This is a reply")
}

func TestOpenGraphTags(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)

	content := "---\n"
	content += "title: The Rock\n"
	content += "description: Dwayne Johnson\n"
	content += "date: Wed, 17 Aug 2022 10:50:02 +0000\n"
	content += "---\n"
	content += "\n"
	content += "Hi \n"

	os.WriteFile(post.ContentFile(), []byte(content), 0644)
	post, _ = user.GetPost(post.Id())
	result, _ := owl.RenderPost(post)

	assertions.AssertContains(t, result, "<meta property=\"og:title\" content=\"The Rock\" />")
	assertions.AssertContains(t, result, "<meta property=\"og:description\" content=\"Dwayne Johnson\" />")
	assertions.AssertContains(t, result, "<meta property=\"og:type\" content=\"article\" />")
	assertions.AssertContains(t, result, "<meta property=\"og:url\" content=\""+post.FullUrl()+"\" />")

}

func TestAddFaviconIfExist(t *testing.T) {
	user := getTestUser()
	os.WriteFile(path.Join(user.MediaDir(), "favicon.png"), []byte("test"), 0644)

	result, _ := owl.RenderIndexPage(user)
	assertions.AssertContains(t, result, "favicon.png")
}

func TestRenderUserAuth(t *testing.T) {
	user := getTestUser()
	user.ResetPassword("test")
	result, err := owl.RenderUserAuthPage(user)
	assertions.AssertNoError(t, err, "Error rendering user auth page")
	assertions.AssertContains(t, result, "<form")
}

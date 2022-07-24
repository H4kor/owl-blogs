package kiss_test

import (
	"h4kor/kiss-social"
	"os"
	"path"
	"strings"
	"testing"
)

func TestCanRenderPost(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result, err := kiss.RenderPost(post)

	if err != nil {
		t.Error("Error rendering post: " + err.Error())
		return
	}

	if !strings.Contains(result, "<h1>testpost</h1>") {
		t.Error("Post title not rendered as h1. Got: " + result)
	}

}

func TestRendererUsesBaseTemplate(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result, err := kiss.RenderPost(post)

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
	result, _ := kiss.RenderIndexPage(user)
	if !strings.Contains(result, "testpost1") {
		t.Error("Post title not rendered. Got: " + result)
	}
	if !strings.Contains(result, "testpost2") {
		t.Error("Post title not rendered. Got: " + result)
	}
}

func TestRenderIndexPageWithBrokenBaseTemplate(t *testing.T) {
	user := getTestUser()
	user.CreateNewPost("testpost1")
	user.CreateNewPost("testpost2")

	os.WriteFile(path.Join(user.Dir(), "meta/base.html"), []byte("{{content}}"), 0644)

	_, err := kiss.RenderIndexPage(user)
	if err == nil {
		t.Error("Expected error rendering index page, got nil")
	}
}

func TestRenderUserList(t *testing.T) {
	repo := getTestRepo()
	repo.CreateUser("user1")
	repo.CreateUser("user2")

	result, err := kiss.RenderUserList(repo)
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

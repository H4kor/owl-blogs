package kiss_test

import (
	"h4kor/kiss-social"
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
		t.Error("Post title not rendered as h1. Got: " + result)
	}
	if !strings.Contains(result, "testpost2") {
		t.Error("Post title not rendered as h1. Got: " + result)
	}
}

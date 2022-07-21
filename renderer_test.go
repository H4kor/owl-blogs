package kiss_test

import (
	"h4kor/kiss-social"
	"strings"
	"testing"
)

func TestCanRenderPost(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result, _ := kiss.RenderPost(post)
	if !strings.Contains(result, "<h1>testpost</h1>") {
		t.Error("Post title not rendered as h1. Got: " + result)
	}

}

func TestRendererUsesBaseTemplate(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result, _ := kiss.RenderPost(post)
	if !strings.Contains(result, "<html>") {
		t.Error("Base template not used. Got: " + result)
	}
}

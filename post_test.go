package owl_test

import (
	"path"
	"testing"
)

func TestCanGetPostTitle(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result := post.Title()
	if result != "testpost" {
		t.Error("Wrong Title. Got: " + result)
	}
}

func TestMediaDir(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result := post.MediaDir()
	if result != path.Join(post.Dir(), "media") {
		t.Error("Wrong MediaDir. Got: " + result)
	}
}

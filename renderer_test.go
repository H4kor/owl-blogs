package kiss_test

import (
	"h4kor/kiss-social"
	"strings"
	"testing"
)

func getTestUser() kiss.User {
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	return user
}

func TestCanRenderPost(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result := kiss.RenderPost(post)
	if !strings.Contains(result, "<h1>testpost</h1>") {
		t.Error("Post title not rendered as h1. Got: " + result)
	}

}

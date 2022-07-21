package kiss_test

import (
	"fmt"
	"h4kor/kiss-social"
	"io/ioutil"
	"path"
	"testing"
)

func TestCreateNewPostCreatesEntryInPublic(t *testing.T) {
	// Create a new user
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost")
	files, err := ioutil.ReadDir(path.Join(user.Dir(), "public"))
	if err != nil {
		t.Error("Error reading directory")
	}
	if len(files) < 1 {
		t.Error("Post not created")
	}
}

func TestCreateNewPostMultipleCalls(t *testing.T) {
	// Create a new user
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost")
	user.CreateNewPost("testpost")
	user.CreateNewPost("testpost")
	files, err := ioutil.ReadDir(path.Join(user.Dir(), "public"))
	if err != nil {
		t.Error("Error reading directory")
	}
	if len(files) < 3 {
		t.Error(fmt.Sprintf("Only %d posts created", len(files)))
	}
}

func TestCanListUserPosts(t *testing.T) {
	// Create a new user
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost")
	user.CreateNewPost("testpost")
	user.CreateNewPost("testpost")
	posts, err := user.Posts()
	if err != nil {
		t.Error("Error reading posts")
	}
	if len(posts) != 3 {
		t.Error("No posts found")
	}
}

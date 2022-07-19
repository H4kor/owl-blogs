package main_test

import (
	"testing"
	"os"
	"fmt"
	"path"
	"time"
	"io/ioutil"
    "math/rand"
	"h4kor/kiss-social"
)


func testRepo() string {
	return "/tmp/test"
}

func randomUserName() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
    b := make([]rune, 8)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func TestCanCreateANewUser(t *testing.T) {
	// Create a new user
	repo := testRepo()
	user := randomUserName()
	main.CreateNewUser(repo, user)
	if _, err := os.Stat(path.Join(repo, user, "")); err != nil {
		t.Error("User directory not created")
	}
}

func TestCreateUserAddsVersionFile(t *testing.T) {
	// Create a new user
	repo := testRepo()
	user := randomUserName()
	main.CreateNewUser(repo, user)
	if _, err := os.Stat(path.Join(repo, user, "/meta/VERSION")); err != nil {
		t.Error("Version file not created")
	}
}

func TestCreateUserAddsBaseHtmlFile(t *testing.T) {
	// Create a new user
	repo := testRepo()
	user := randomUserName()
	main.CreateNewUser(repo, user)
	if _, err := os.Stat(path.Join(repo, user, "/meta/base.html")); err != nil {
		t.Error("Base html file not created")
	}
}

func TestCreateUserAddsPublicFolder(t *testing.T) {
	// Create a new user
	repo := testRepo()
	user := randomUserName()
	main.CreateNewUser(repo, user)
	if _, err := os.Stat(path.Join(repo, user, "/public")); err != nil {
		t.Error("Public folder not created")
	}
}

func TestCreateNewPostCreatesEntryInPublic(t *testing.T) {
	// Create a new user
	repo := testRepo()
	user := randomUserName()
	main.CreateNewUser(repo, user)
	// Create a new post
	main.CreateNewPost(repo, user, "testpost")
	files, err := ioutil.ReadDir(path.Join(repo,  user,  "public"))
    if err != nil {
		t.Error("Error reading directory")
	}
	if len(files) < 1 {
		t.Error("Post not created")
	}
}

func TestCreateNewPostMultipleCalls(t *testing.T) {
	// Create a new user
	repo := testRepo()
	user := randomUserName()
	main.CreateNewUser(repo, user)
	// Create a new post
	main.CreateNewPost(repo, user, "testpost")
	main.CreateNewPost(repo, user, "testpost")
	main.CreateNewPost(repo, user, "testpost")
	files, err := ioutil.ReadDir(path.Join(repo,  user,  "public"))
    if err != nil {
		t.Error("Error reading directory")
	}
	if len(files) < 3 {
		t.Error(fmt.Sprintf("Only %d posts created", len(files)))
	}
}
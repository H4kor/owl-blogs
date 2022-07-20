package kiss_test

import (
	"fmt"
	"h4kor/kiss-social"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"
)

func randomName() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func testRepoName() string {
	return "/tmp/" + randomName()
}

func randomUserName() string {
	return randomName()
}

func TestCanCreateRepository(t *testing.T) {
	repoName := testRepoName()
	_, err := kiss.CreateRepository(repoName)
	if err != nil {
		t.Error("Error creating repository: ", err)
	}

}

func TestCannotCreateExistingRepository(t *testing.T) {
	repoName := testRepoName()
	kiss.CreateRepository(repoName)
	_, err := kiss.CreateRepository(repoName)
	if err == nil {
		t.Error("No error returned when creating existing repository")
	}
}

func TestCanCreateANewUser(t *testing.T) {
	// Create a new user
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := kiss.CreateNewUser(repo, randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "")); err != nil {
		t.Error("User directory not created")
	}
}

func TestCannotRecreateExisitingUser(t *testing.T) {
	// Create a new user
	repo, _ := kiss.CreateRepository(testRepoName())
	userName := randomUserName()
	kiss.CreateNewUser(repo, userName)
	_, err := kiss.CreateNewUser(repo, userName)
	if err == nil {
		t.Error("No error returned when creating existing user")
	}
}

func TestCreateUserAddsVersionFile(t *testing.T) {
	// Create a new user
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := kiss.CreateNewUser(repo, randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/meta/VERSION")); err != nil {
		t.Error("Version file not created")
	}
}

func TestCreateUserAddsBaseHtmlFile(t *testing.T) {
	// Create a new user
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := kiss.CreateNewUser(repo, randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/meta/base.html")); err != nil {
		t.Error("Base html file not created")
	}
}

func TestCreateUserAddsPublicFolder(t *testing.T) {
	// Create a new user
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := kiss.CreateNewUser(repo, randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/public")); err != nil {
		t.Error("Public folder not created")
	}
}

func TestCreateNewPostCreatesEntryInPublic(t *testing.T) {
	// Create a new user
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := kiss.CreateNewUser(repo, randomUserName())
	// Create a new post
	kiss.CreateNewPost(user, "testpost")
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
	user, _ := kiss.CreateNewUser(repo, randomUserName())
	// Create a new post
	kiss.CreateNewPost(user, "testpost")
	kiss.CreateNewPost(user, "testpost")
	kiss.CreateNewPost(user, "testpost")
	files, err := ioutil.ReadDir(path.Join(user.Dir(), "public"))
	if err != nil {
		t.Error("Error reading directory")
	}
	if len(files) < 3 {
		t.Error(fmt.Sprintf("Only %d posts created", len(files)))
	}
}

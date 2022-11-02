package web_test

import (
	"h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/priv/assertions"
	"math/rand"
	"net/http"
	"net/http/httptest"
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

func getTestRepo(config owl.RepoConfig) owl.Repository {
	repo, _ := owl.CreateRepository(testRepoName(), config)
	return repo
}

func TestMultiUserRepoIndexHandler(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.CreateUser("user_1")
	repo.CreateUser("user_2")

	// Create Request and Response
	req, err := http.NewRequest("GET", "/", nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)

	// Check the response body contains names of users
	assertions.AssertContains(t, rr.Body.String(), "user_1")
	assertions.AssertContains(t, rr.Body.String(), "user_2")
}

func TestMultiUserUserIndexHandler(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	user.CreateNewPost("post-1", false)

	// Create Request and Response
	req, err := http.NewRequest("GET", user.UrlPath(), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)

	// Check the response body contains names of users
	assertions.AssertContains(t, rr.Body.String(), "post-1")
}

func TestMultiUserPostHandler(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1", false)

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlPath(), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)
}

func TestMultiUserPostMediaHandler(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1", false)

	// Create test media file
	path := path.Join(post.MediaDir(), "data.txt")
	err := os.WriteFile(path, []byte("test"), 0644)
	assertions.AssertNoError(t, err, "Error creating request")

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlMediaPath("data.txt"), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)

	// Check the response body contains data of media file
	if !(rr.Body.String() == "test") {
		t.Error("Got wrong media file content. Expected 'test' Got: ")
		t.Error(rr.Body.String())
	}
}

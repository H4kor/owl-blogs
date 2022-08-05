package main_test

import (
	"h4kor/owl-blogs"
	"h4kor/owl-blogs/cmd/owl-web"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
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

func getTestRepo() owl.Repository {
	repo, _ := owl.CreateRepository(testRepoName())
	return repo
}

func TestMultiUserRepoIndexHandler(t *testing.T) {
	repo := getTestRepo()
	repo.CreateUser("user_1")
	repo.CreateUser("user_2")

	// Create Request and Response
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body contains names of users
	if !strings.Contains(rr.Body.String(), "user_1") {
		t.Error("user_1 not listed on index page. Got: ")
		t.Error(rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "user_2") {
		t.Error("user_2 not listed on index page. Got: ")
		t.Error(rr.Body.String())
	}
}

func TestMultiUserUserIndexHandler(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	user.CreateNewPost("post-1")

	// Create Request and Response
	req, err := http.NewRequest("GET", user.UrlPath(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body contains names of users
	if !strings.Contains(rr.Body.String(), "post-1") {
		t.Error("post-1 not listed on index page. Got: ")
		t.Error(rr.Body.String())
	}
}

func TestMultiUserPostHandler(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlPath(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestMultiUserPostMediaHandler(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	// Create test media file
	path := path.Join(post.MediaDir(), "data.txt")
	err := os.WriteFile(path, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlMediaPath("data.txt"), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body contains data of media file
	if !(rr.Body.String() == "test") {
		t.Error("Got wrong media file content. Expected 'test' Got: ")
		t.Error(rr.Body.String())
	}
}

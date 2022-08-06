package main_test

import (
	owl "h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl-web"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

func getSingleUserTestRepo() (owl.Repository, owl.User) {
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser("test-1")
	repo.SetSingleUser(user)
	return repo, user
}

func TestSingleUserUserIndexHandler(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.CreateNewPost("post-1")

	// Create Request and Response
	req, err := http.NewRequest("GET", user.UrlPath(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
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

func TestSingleUserPostHandler(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	post, _ := user.CreateNewPost("post-1")

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlPath(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSingleUserPostMediaHandler(t *testing.T) {
	repo, user := getSingleUserTestRepo()
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
	router := main.SingleUserRouter(&repo)
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

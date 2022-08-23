package main_test

import (
	main "h4kor/owl-blogs/cmd/owl-web"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestWebmentionHandleAccepts(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	target := post.FullUrl()
	source := "https://example.com"
	data := url.Values{}
	data.Set("target", target)
	data.Set("source", source)

	// Create Request and Response
	req, err := http.NewRequest("POST", user.UrlPath()+"webmention/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
		t.Errorf("Body: %v", rr.Body)
	}

}

func TestWebmentionWrittenToPost(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	target := post.FullUrl()
	source := "https://example.com"
	data := url.Values{}
	data.Set("target", target)
	data.Set("source", source)

	// Create Request and Response
	req, err := http.NewRequest("POST", user.UrlPath()+"webmention/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
		return
	}

	if len(post.Webmentions()) != 1 {
		t.Errorf("no webmention written to post")
	}

}

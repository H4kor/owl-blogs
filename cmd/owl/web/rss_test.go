package web_test

import (
	"h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl/web"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMultiUserUserRssIndexHandler(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	user.CreateNewPost("post-1", false)

	// Create Request and Response
	req, err := http.NewRequest("GET", user.UrlPath()+"index.xml", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response Content-Type is what we expect.
	if !strings.Contains(rr.Header().Get("Content-Type"), "application/rss+xml") {
		t.Errorf("handler returned wrong Content-Type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/rss+xml")
	}

	// Check the response body contains names of users
	if !strings.Contains(rr.Body.String(), "post-1") {
		t.Error("post-1 not listed on index page. Got: ")
		t.Error(rr.Body.String())
	}
}

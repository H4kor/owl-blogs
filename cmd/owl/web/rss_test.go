package web_test

import (
	"h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/test/assertions"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMultiUserUserRssIndexHandler(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	user.CreateNewPost(owl.PostMeta{Type: "article", Title: "post-1"}, "")

	// Create Request and Response
	req, err := http.NewRequest("GET", user.UrlPath()+"index.xml", nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)

	// Check the response Content-Type is what we expect.
	assertions.AssertContains(t, rr.Header().Get("Content-Type"), "application/rss+xml")

	// Check the response body contains names of users
	assertions.AssertContains(t, rr.Body.String(), "post-1")
}

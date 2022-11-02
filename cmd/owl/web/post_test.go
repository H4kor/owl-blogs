package web_test

import (
	"h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/priv/assertions"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPostHandlerReturns404OnDrafts(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1", false)

	content := "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "---\n"
	content += "\n"
	content += "Write your post here.\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlPath(), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusNotFound)
}

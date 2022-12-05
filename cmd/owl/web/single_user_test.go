package web_test

import (
	owl "h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/test/assertions"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func getSingleUserTestRepo() (owl.Repository, owl.User) {
	repo, _ := owl.CreateRepository(testRepoName(), owl.RepoConfig{SingleUser: "test-1"})
	user, _ := repo.CreateUser("test-1")
	return repo, user
}

func TestSingleUserUserIndexHandler(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.CreateNewPost(owl.PostMeta{Type: "article", Title: "post-1"}, "")

	// Create Request and Response
	req, err := http.NewRequest("GET", user.UrlPath(), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)

	// Check the response body contains names of users
	assertions.AssertContains(t, rr.Body.String(), "post-1")
}

func TestSingleUserPostHandler(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "post-1"}, "")

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlPath(), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)
}

func TestSingleUserPostMediaHandler(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "post-1"}, "")

	// Create test media file
	path := path.Join(post.MediaDir(), "data.txt")
	err := os.WriteFile(path, []byte("test"), 0644)
	assertions.AssertNoError(t, err, "Error creating request")

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlMediaPath("data.txt"), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)

	// Check the response body contains data of media file
	assertions.Assert(t, rr.Body.String() == "test", "Media file data not returned")
}

func TestHasNoDraftsInList(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "post-1"}, "")
	content := ""
	content += "---\n"
	content += "title: Articles September 2019\n"
	content += "author: h4kor\n"
	content += "type: post\n"
	content += "date: -001-11-30T00:00:00+00:00\n"
	content += "draft: true\n"
	content += "url: /?p=426\n"
	content += "categories:\n"
	content += "  - Uncategorised\n"
	content += "\n"
	content += "---\n"
	content += "<https://nesslabs.com/time-anxiety>\n"

	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", "/", nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	// Check if title is in the response body
	assertions.AssertNotContains(t, rr.Body.String(), "Articles September 2019")
}

func TestSingleUserUserPostListHandler(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.CreateNewPost(owl.PostMeta{
		Title: "post-1",
		Type:  "article",
	}, "hi")
	user.CreateNewPost(owl.PostMeta{
		Title: "post-2",
		Type:  "note",
	}, "hi")
	list := owl.PostList{
		Title:   "list-1",
		Id:      "list-1",
		Include: []string{"article"},
	}
	user.AddPostList(list)

	// Create Request and Response
	req, err := http.NewRequest("GET", user.ListUrl(list), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)

	// Check the response body contains names of users
	assertions.AssertContains(t, rr.Body.String(), "post-1")
	assertions.AssertNotContains(t, rr.Body.String(), "post-2")
}

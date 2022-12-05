package web_test

import (
	"h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/test/assertions"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRedirectOnAliases(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "post-1"}, "")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - /foo/bar\n"
	content += "  - /foo/baz\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", "/foo/bar", nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusMovedPermanently)
	// Check that Location header is set correctly
	assertions.AssertEqual(t, rr.Header().Get("Location"), post.UrlPath())
}

func TestNoRedirectOnNonExistingAliases(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "post-1"}, "")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - /foo/bar\n"
	content += "  - /foo/baz\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", "/foo/bar2", nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusNotFound)

}

func TestNoRedirectIfValidPostUrl(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "post-1"}, "")
	post2, _ := user.CreateNewPost(owl.PostMeta{Title: "post-2"}, "")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - " + post2.UrlPath() + "\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", post2.UrlPath(), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)

}

func TestRedirectIfInvalidPostUrl(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "post-1"}, "")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - " + user.UrlPath() + "posts/not-a-real-post/" + "\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", user.UrlPath()+"posts/not-a-real-post/", nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusMovedPermanently)

}

func TestRedirectIfInvalidUserUrl(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "post-1"}, "")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - /user/not-real/ \n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", "/user/not-real/", nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusMovedPermanently)

}

func TestRedirectIfInvalidMediaUrl(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "post-1"}, "")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - " + post.UrlMediaPath("not-real") + "\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlMediaPath("not-real"), nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusMovedPermanently)

}

func TestDeepAliasInSingleUserMode(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{SingleUser: "test-1"})
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "post-1"}, "")

	content := "---\n"
	content += "title: Create tileable textures with GIMP\n"
	content += "author: h4kor\n"
	content += "type: post\n"
	content += "date: Tue, 13 Sep 2016 16:19:09 +0000\n"
	content += "aliases:\n"
	content += "  - /2016/09/13/create-tileable-textures-with-gimp/\n"
	content += "categories:\n"
	content += "  - GameDev\n"
	content += "tags:\n"
	content += "  - gamedev\n"
	content += "  - textures\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", "/2016/09/13/create-tileable-textures-with-gimp/", nil)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusMovedPermanently)

}

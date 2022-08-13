package main_test

import (
	main "h4kor/owl-blogs/cmd/owl-web"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRedirectOnAliases(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

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
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}

	// Check that Location header is set correctly
	if rr.Header().Get("Location") != post.UrlPath() {
		t.Errorf("Location header is not set correctly, expected: %v Got: %v",
			post.UrlPath(),
			rr.Header().Get("Location"),
		)
	}
}

func TestNoRedirectOnNonExistingAliases(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

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
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

}

func TestNoRedirectIfValidPostUrl(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")
	post2, _ := user.CreateNewPost("post-2")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - " + post2.UrlPath() + "\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", post2.UrlPath(), nil)
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

}

func TestRedirectIfInvalidPostUrl(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - " + user.UrlPath() + "posts/not-a-real-post/" + "\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", user.UrlPath()+"posts/not-a-real-post/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}

}

func TestRedirectIfInvalidUserUrl(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - /user/not-real/ \n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", "/user/not-real/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}

}

func TestRedirectIfInvalidMediaUrl(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - " + post.UrlMediaPath("not-real") + "\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", post.UrlMediaPath("not-real"), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}

}

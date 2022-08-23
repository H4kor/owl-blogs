package main_test

import (
	"h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl-web"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func setupWebmentionTest(repo owl.Repository, user owl.User, target string, source string) (*httptest.ResponseRecorder, error) {

	data := url.Values{}
	data.Set("target", target)
	data.Set("source", source)

	// Create Request and Response
	req, err := http.NewRequest("POST", user.UrlPath()+"webmention/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	router := main.Router(&repo)
	router.ServeHTTP(rr, req)

	return rr, nil
}

func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, expStatus int) {
	if status := rr.Code; status != expStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expStatus)
		return
	}
}

func TestWebmentionHandleAccepts(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	target := post.FullUrl()
	source := "https://example.com"

	rr, err := setupWebmentionTest(repo, user, target, source)
	if err != nil {
		t.Fatal(err)
	}

	assertStatus(t, rr, http.StatusAccepted)

}

func TestWebmentionWrittenToPost(t *testing.T) {

	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	target := post.FullUrl()
	source := "https://example.com"

	rr, err := setupWebmentionTest(repo, user, target, source)
	if err != nil {
		t.Fatal(err)
	}

	// Check the status code is what we expect.
	assertStatus(t, rr, http.StatusAccepted)

	if len(post.Webmentions()) != 1 {
		t.Errorf("no webmention written to post")
	}
}

//
// https://www.w3.org/TR/webmention/#h-request-verification
//

// The receiver MUST check that source and target are valid URLs [URL]
// and are of schemes that are supported by the receiver.
// (Most commonly this means checking that the source and target schemes are http or https).
func TestWebmentionSourceValidation(t *testing.T) {

	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	target := post.FullUrl()
	source := "ftp://example.com"

	rr, err := setupWebmentionTest(repo, user, target, source)
	if err != nil {
		t.Fatal(err)
	}

	assertStatus(t, rr, http.StatusBadRequest)
}

func TestWebmentionTargetValidation(t *testing.T) {

	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	target := "ftp://example.com"
	source := post.FullUrl()

	rr, err := setupWebmentionTest(repo, user, target, source)
	if err != nil {
		t.Fatal(err)
	}

	assertStatus(t, rr, http.StatusBadRequest)
}

// The receiver MUST reject the request if the source URL is the same as the target URL.

func TestWebmentionSameTargetAndSource(t *testing.T) {

	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	target := post.FullUrl()
	source := post.FullUrl()

	rr, err := setupWebmentionTest(repo, user, target, source)
	if err != nil {
		t.Fatal(err)
	}

	assertStatus(t, rr, http.StatusBadRequest)
}

// The receiver SHOULD check that target is a valid resource for which it can accept Webmentions.
// This check SHOULD happen synchronously to reject invalid Webmentions before more in-depth verification begins.
// What a "valid resource" means is up to the receiver.
func TestValidationOfTarget(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("test-1")
	post, _ := user.CreateNewPost("post-1")

	target := post.FullUrl()
	target = target[:len(target)-1] + "invalid"
	source := post.FullUrl()

	rr, err := setupWebmentionTest(repo, user, target, source)
	if err != nil {
		t.Fatal(err)
	}

	assertStatus(t, rr, http.StatusBadRequest)
}

//
// https://www.w3.org/TR/webmention/#h-webmention-verification
//

package web_test

import (
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/test/assertions"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestLoginWrongPassword(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("password", "wrongpassword")
	form.Add("csrf_token", csrfToken)

	req, err := http.NewRequest("POST", user.EditorLoginUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	// check redirect to login page

	assertions.AssertNotEqual(t, rr.Header().Get("Location"), user.EditorUrl())
}

func TestLoginCorrectPassword(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("password", "testpassword")
	form.Add("csrf_token", csrfToken)

	req, err := http.NewRequest("POST", user.EditorLoginUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	// check redirect to login page
	assertions.AssertStatus(t, rr, http.StatusFound)
	assertions.AssertEqual(t, rr.Header().Get("Location"), user.EditorUrl())
}

func TestEditorWithoutSession(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	user.CreateNewSession()

	req, err := http.NewRequest("GET", user.EditorUrl(), nil)
	// req.AddCookie(&http.Cookie{Name: "session", Value: sessionId})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusFound)
	assertions.AssertEqual(t, rr.Header().Get("Location"), user.EditorLoginUrl())

}

func TestEditorWithSession(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	sessionId := user.CreateNewSession()

	req, err := http.NewRequest("GET", user.EditorUrl(), nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: sessionId})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)
}

func TestEditorPostWithoutSession(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	user.CreateNewSession()

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("type", "article")
	form.Add("title", "testtitle")
	form.Add("content", "testcontent")
	form.Add("csrf_token", csrfToken)

	req, err := http.NewRequest("POST", user.EditorUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusFound)
	assertions.AssertEqual(t, rr.Header().Get("Location"), user.EditorLoginUrl())
}

func TestEditorPostWithSession(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	sessionId := user.CreateNewSession()

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("type", "article")
	form.Add("title", "testtitle")
	form.Add("content", "testcontent")
	form.Add("csrf_token", csrfToken)

	req, err := http.NewRequest("POST", user.EditorUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	req.AddCookie(&http.Cookie{Name: "session", Value: sessionId})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	posts, _ := user.AllPosts()
	assertions.AssertEqual(t, len(posts), 1)
	post := posts[0]

	assertions.AssertStatus(t, rr, http.StatusFound)
	assertions.AssertEqual(t, rr.Header().Get("Location"), post.FullUrl())
}

func TestEditorPostWithSessionNote(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	sessionId := user.CreateNewSession()

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("type", "note")
	form.Add("content", "testcontent")
	form.Add("csrf_token", csrfToken)

	req, err := http.NewRequest("POST", user.EditorUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	req.AddCookie(&http.Cookie{Name: "session", Value: sessionId})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	posts, _ := user.AllPosts()
	assertions.AssertEqual(t, len(posts), 1)
	post := posts[0]

	assertions.AssertStatus(t, rr, http.StatusFound)
	assertions.AssertEqual(t, rr.Header().Get("Location"), post.FullUrl())
}

package web_test

import (
	"h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/test/assertions"
	"h4kor/owl-blogs/test/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

type CountMockHttpClient struct {
	InvokedGet      int
	InvokedPost     int
	InvokedPostForm int
}

func (c *CountMockHttpClient) Get(url string) (resp *http.Response, err error) {
	c.InvokedGet++
	return &http.Response{}, nil
}

func (c *CountMockHttpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	c.InvokedPost++
	return &http.Response{}, nil
}

func (c *CountMockHttpClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	c.InvokedPostForm++
	return &http.Response{}, nil
}

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

	assertions.AssertEqual(t, rr.Header().Get("Location"), user.EditorLoginUrl()+"?error=wrong_password")
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

func TestEditorSendsWebmentions(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	repo.HttpClient = &CountMockHttpClient{}
	repo.Parser = &mocks.MockHtmlParser{}
	user.ResetPassword("testpassword")

	mentioned_post, _ := user.CreateNewPost(owl.PostMeta{Title: "test"}, "")

	sessionId := user.CreateNewSession()

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("type", "note")
	form.Add("content", "[test]("+mentioned_post.FullUrl()+")")
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
	assertions.AssertEqual(t, len(posts), 2)
	post := posts[0]
	assertions.AssertLen(t, post.OutgoingWebmentions(), 1)
	assertions.AssertStatus(t, rr, http.StatusFound)
	assertions.AssertEqual(t, repo.HttpClient.(*CountMockHttpClient).InvokedPostForm, 1)

}

func TestEditorPostWithSessionRecipe(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	sessionId := user.CreateNewSession()

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("type", "recipe")
	form.Add("title", "testtitle")
	form.Add("yield", "2")
	form.Add("duration", "1 hour")
	form.Add("ingredients", "water\nwheat")
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

	assertions.AssertLen(t, post.Meta().Recipe.Ingredients, 2)

	assertions.AssertStatus(t, rr, http.StatusFound)
	assertions.AssertEqual(t, rr.Header().Get("Location"), post.FullUrl())
}

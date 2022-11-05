package web_test

import (
	"encoding/json"
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/priv/assertions"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestAuthPostWrongPassword(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("password", "wrongpassword")
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("response_type", "code")
	form.Add("state", "test_state")
	form.Add("csrf_token", csrfToken)

	req, err := http.NewRequest("POST", user.AuthUrl()+"verify/", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusFound)
	assertions.AssertContains(t, rr.Header().Get("Location"), "error=invalid_password")
}

func TestAuthPostCorrectPassword(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("password", "testpassword")
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("response_type", "code")
	form.Add("state", "test_state")
	form.Add("csrf_token", csrfToken)
	req, err := http.NewRequest("POST", user.AuthUrl()+"verify/", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusFound)
	assertions.AssertContains(t, rr.Header().Get("Location"), "code=")
	assertions.AssertContains(t, rr.Header().Get("Location"), "state=test_state")
	assertions.AssertContains(t, rr.Header().Get("Location"), "http://example.com/response")
}

func TestAuthPostWithIncorrectCode(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	user.GenerateAuthCode("http://example.com", "http://example.com/response")

	// Create Request and Response
	form := url.Values{}
	form.Add("code", "wrongcode")
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("grant_type", "authorization_code")
	req, err := http.NewRequest("POST", user.AuthUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusUnauthorized)
}

func TestAuthPostWithCorrectCode(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	code, _ := user.GenerateAuthCode("http://example.com", "http://example.com/response")

	// Create Request and Response
	form := url.Values{}
	form.Add("code", code)
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("grant_type", "authorization_code")
	req, err := http.NewRequest("POST", user.AuthUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Accept", "application/json")
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)
	// parse response as json
	type responseType struct {
		Me string `json:"me"`
	}
	var response responseType
	json.Unmarshal(rr.Body.Bytes(), &response)
	assertions.AssertEqual(t, response.Me, user.FullUrl())

}

package web_test

import (
	"h4kor/owl-blogs"
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/test/assertions"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestMicropubMinimalArticle(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")

	code, _ := user.GenerateAuthCode(
		"test", "test", "test", "test", "test",
	)
	token, _, _ := user.GenerateAccessToken(owl.AuthCode{
		Code:                code,
		ClientId:            "test",
		RedirectUri:         "test",
		CodeChallenge:       "test",
		CodeChallengeMethod: "test",
		Scope:               "test",
	})

	// Create Request and Response
	form := url.Values{}
	form.Add("h", "entry")
	form.Add("name", "Test Article")
	form.Add("content", "Test Content")

	req, err := http.NewRequest("POST", user.MicropubUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Authorization", "Bearer "+token)
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusCreated)
}

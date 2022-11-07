package web_test

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	main "h4kor/owl-blogs/cmd/owl/web"
	"h4kor/owl-blogs/test/assertions"
	"h4kor/owl-blogs/test/mocks"
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
	assertions.AssertContains(t, rr.Header().Get("Location"), "iss="+user.FullUrl())
	assertions.AssertContains(t, rr.Header().Get("Location"), "http://example.com/response")
}

func TestAuthPostWithIncorrectCode(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	user.GenerateAuthCode("http://example.com", "http://example.com/response", "", "", "profile")

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
	code, _ := user.GenerateAuthCode("http://example.com", "http://example.com/response", "", "", "profile")

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

func TestAuthPostWithCorrectCodeAndPKCE(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")

	// Create Request and Response
	code_verifier := "test_code_verifier"
	// create code challenge
	h := sha256.New()
	h.Write([]byte(code_verifier))
	code_challenge := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	code, _ := user.GenerateAuthCode("http://example.com", "http://example.com/response", code_challenge, "S256", "profile")

	form := url.Values{}
	form.Add("code", code)
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("grant_type", "authorization_code")
	form.Add("code_verifier", code_verifier)
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

func TestAuthPostWithCorrectCodeAndWrongPKCE(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")

	// Create Request and Response
	code_verifier := "test_code_verifier"
	// create code challenge
	h := sha256.New()
	h.Write([]byte(code_verifier + "wrong"))
	code_challenge := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	code, _ := user.GenerateAuthCode("http://example.com", "http://example.com/response", code_challenge, "S256", "profile")

	form := url.Values{}
	form.Add("code", code)
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("grant_type", "authorization_code")
	form.Add("code_verifier", code_verifier)
	req, err := http.NewRequest("POST", user.AuthUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Accept", "application/json")
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusUnauthorized)
}

func TestAuthPostWithCorrectCodePKCEPlain(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")

	// Create Request and Response
	code_verifier := "test_code_verifier"
	code_challenge := code_verifier
	code, _ := user.GenerateAuthCode("http://example.com", "http://example.com/response", code_challenge, "plain", "profile")

	form := url.Values{}
	form.Add("code", code)
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("grant_type", "authorization_code")
	form.Add("code_verifier", code_verifier)
	req, err := http.NewRequest("POST", user.AuthUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Accept", "application/json")
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)
}

func TestAuthPostWithCorrectCodePKCEPlainWrong(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")

	// Create Request and Response
	code_verifier := "test_code_verifier"
	code_challenge := code_verifier + "wrong"
	code, _ := user.GenerateAuthCode("http://example.com", "http://example.com/response", code_challenge, "plain", "profile")

	form := url.Values{}
	form.Add("code", code)
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("grant_type", "authorization_code")
	form.Add("code_verifier", code_verifier)
	req, err := http.NewRequest("POST", user.AuthUrl(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Accept", "application/json")
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusUnauthorized)
}

func TestAuthRedirectUriNotSet(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockParseLinksHtmlParser{
		Links: []string{"http://example2.com/response"},
	}
	user.ResetPassword("testpassword")

	csrfToken := "test_csrf_token"

	// Create Request and Response
	form := url.Values{}
	form.Add("password", "wrongpassword")
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example2.com/response_not_set")
	form.Add("response_type", "code")
	form.Add("state", "test_state")
	form.Add("csrf_token", csrfToken)

	req, err := http.NewRequest("GET", user.AuthUrl()+"?"+form.Encode(), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusBadRequest)
}

func TestAuthRedirectUriSet(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockParseLinksHtmlParser{
		Links: []string{"http://example.com/response"},
	}
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

	req, err := http.NewRequest("GET", user.AuthUrl()+"?"+form.Encode(), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)
}

func TestAuthRedirectUriSameHost(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockParseLinksHtmlParser{
		Links: []string{},
	}
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

	req, err := http.NewRequest("GET", user.AuthUrl()+"?"+form.Encode(), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)
}

func TestAccessTokenCorrectPassword(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	code, _ := user.GenerateAuthCode("http://example.com", "http://example.com/response", "", "", "profile create")

	// Create Request and Response
	form := url.Values{}
	form.Add("code", code)
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("grant_type", "authorization_code")
	req, err := http.NewRequest("POST", user.AuthUrl()+"token/", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)
	// parse response as json
	type responseType struct {
		Me           string `json:"me"`
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}
	var response responseType
	json.Unmarshal(rr.Body.Bytes(), &response)
	assertions.AssertEqual(t, response.Me, user.FullUrl())
	assertions.AssertEqual(t, response.TokenType, "Bearer")
	assertions.AssertEqual(t, response.Scope, "profile create")
	assertions.Assert(t, response.ExpiresIn > 0, "ExpiresIn should be greater than 0")
	assertions.Assert(t, len(response.AccessToken) > 0, "AccessToken should be greater than 0")
}

func TestAccessTokenWithIncorrectCode(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	user.GenerateAuthCode("http://example.com", "http://example.com/response", "", "", "profile")

	// Create Request and Response
	form := url.Values{}
	form.Add("code", "wrongcode")
	form.Add("client_id", "http://example.com")
	form.Add("redirect_uri", "http://example.com/response")
	form.Add("grant_type", "authorization_code")
	req, err := http.NewRequest("POST", user.AuthUrl()+"token/", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	assertions.AssertNoError(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusUnauthorized)
}

func TestIndieauthMetadata(t *testing.T) {
	repo, user := getSingleUserTestRepo()
	user.ResetPassword("testpassword")
	req, _ := http.NewRequest("GET", user.IndieauthMetadataUrl(), nil)
	rr := httptest.NewRecorder()
	router := main.SingleUserRouter(&repo)
	router.ServeHTTP(rr, req)

	assertions.AssertStatus(t, rr, http.StatusOK)
	// parse response as json
	type responseType struct {
		Issuer                        string   `json:"issuer"`
		AuthorizationEndpoint         string   `json:"authorization_endpoint"`
		TokenEndpoint                 string   `json:"token_endpoint"`
		CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
		ScopesSupported               []string `json:"scopes_supported"`
		ResponseTypesSupported        []string `json:"response_types_supported"`
		GrantTypesSupported           []string `json:"grant_types_supported"`
	}
	var response responseType
	json.Unmarshal(rr.Body.Bytes(), &response)
	assertions.AssertEqual(t, response.Issuer, user.FullUrl())
	assertions.AssertEqual(t, response.AuthorizationEndpoint, user.AuthUrl())
	assertions.AssertEqual(t, response.TokenEndpoint, user.TokenUrl())
}

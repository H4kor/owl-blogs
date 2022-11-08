package web

import (
	"encoding/json"
	"fmt"
	"h4kor/owl-blogs"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func userAuthMetadataHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		type Response struct {
			Issuer                        string   `json:"issuer"`
			AuthorizationEndpoint         string   `json:"authorization_endpoint"`
			TokenEndpoint                 string   `json:"token_endpoint"`
			CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
			ScopesSupported               []string `json:"scopes_supported"`
			ResponseTypesSupported        []string `json:"response_types_supported"`
			GrantTypesSupported           []string `json:"grant_types_supported"`
		}
		response := Response{
			Issuer:                        user.FullUrl(),
			AuthorizationEndpoint:         user.AuthUrl(),
			TokenEndpoint:                 user.TokenUrl(),
			CodeChallengeMethodsSupported: []string{"S256", "plain"},
			ScopesSupported:               []string{"profile"},
			ResponseTypesSupported:        []string{"code"},
			GrantTypesSupported:           []string{"authorization_code"},
		}
		jsonData, err := json.Marshal(response)
		if err != nil {
			println("Error marshalling json: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
		}
		w.Write(jsonData)
	}
}

func userAuthHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		// get me, cleint_id, redirect_uri, state and response_type from query
		me := r.URL.Query().Get("me")
		clientId := r.URL.Query().Get("client_id")
		redirectUri := r.URL.Query().Get("redirect_uri")
		state := r.URL.Query().Get("state")
		responseType := r.URL.Query().Get("response_type")
		codeChallenge := r.URL.Query().Get("code_challenge")
		codeChallengeMethod := r.URL.Query().Get("code_challenge_method")
		scope := r.URL.Query().Get("scope")

		// check if request is valid
		missing_params := []string{}
		if clientId == "" {
			missing_params = append(missing_params, "client_id")
		}
		if redirectUri == "" {
			missing_params = append(missing_params, "redirect_uri")
		}
		if responseType == "" {
			missing_params = append(missing_params, "response_type")
		}
		if len(missing_params) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Missing parameters: %s", strings.Join(missing_params, ", "))))
			return
		}
		if responseType != "id" {
			responseType = "code"
		}
		if responseType != "code" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid response_type. Must be 'code' ('id' converted to 'code' for legacy support)."))
			return
		}
		if codeChallengeMethod != "" && (codeChallengeMethod != "S256" && codeChallengeMethod != "plain") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid code_challenge_method. Must be 'S256' or 'plain'."))
			return
		}

		client_id_url, err := url.Parse(clientId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid client_id."))
			return
		}
		redirect_uri_url, err := url.Parse(redirectUri)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid redirect_uri."))
			return
		}
		if client_id_url.Host != redirect_uri_url.Host || client_id_url.Scheme != redirect_uri_url.Scheme {
			// check if redirect_uri is registered
			resp, _ := repo.HttpClient.Get(clientId)
			registered_redirects, _ := repo.Parser.GetRedirctUris(resp)
			is_registered := false
			for _, registered_redirect := range registered_redirects {
				if registered_redirect == redirectUri {
					// redirect_uri is registered
					is_registered = true
					break
				}
			}
			if !is_registered {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid redirect_uri. Must be registered with client_id."))
				return
			}
		}

		// Double Submit Cookie Pattern
		// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#double-submit-cookie
		csrfToken := owl.GenerateRandomString(32)
		cookie := http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)

		reqData := owl.AuthRequestData{
			Me:                  me,
			ClientId:            clientId,
			RedirectUri:         redirectUri,
			State:               state,
			Scope:               scope,
			ResponseType:        responseType,
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: codeChallengeMethod,
			User:                user,
			CsrfToken:           csrfToken,
		}

		html, err := owl.RenderUserAuthPage(reqData)
		if err != nil {
			println("Error rendering auth page: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		println("Rendering auth page for user", user.Name())
		w.Write([]byte(html))
	}
}

func verifyAuthCodeRequest(user owl.User, w http.ResponseWriter, r *http.Request) (bool, owl.AuthCode) {
	// get form data from post request
	err := r.ParseForm()
	if err != nil {
		println("Error parsing form: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error parsing form"))
		return false, owl.AuthCode{}
	}
	code := r.Form.Get("code")
	client_id := r.Form.Get("client_id")
	redirect_uri := r.Form.Get("redirect_uri")
	code_verifier := r.Form.Get("code_verifier")

	// check if request is valid
	valid, authCode := user.VerifyAuthCode(code, client_id, redirect_uri, code_verifier)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid code"))
	}
	return valid, authCode
}

func userAuthProfileHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		valid, _ := verifyAuthCodeRequest(user, w, r)
		if valid {
			w.WriteHeader(http.StatusOK)
			type ResponseProfile struct {
				Name  string `json:"name"`
				Url   string `json:"url"`
				Photo string `json:"photo"`
			}
			type Response struct {
				Me      string          `json:"me"`
				Profile ResponseProfile `json:"profile"`
			}
			response := Response{
				Me: user.FullUrl(),
				Profile: ResponseProfile{
					Name:  user.Name(),
					Url:   user.FullUrl(),
					Photo: user.AvatarUrl(),
				},
			}
			jsonData, err := json.Marshal(response)
			if err != nil {
				println("Error marshalling json: ", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
			}
			w.Write(jsonData)
			return
		}
	}
}

func userAuthTokenHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		valid, authCode := verifyAuthCodeRequest(user, w, r)
		if valid {
			if authCode.Scope == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Empty scope, no token issued"))
				return
			}

			type Response struct {
				Me           string `json:"me"`
				TokenType    string `json:"token_type"`
				AccessToken  string `json:"access_token"`
				Scope        string `json:"scope"`
				ExpiresIn    int    `json:"expires_in"`
				RefreshToken string `json:"refresh_token"`
			}
			accessToken, duration, err := user.GenerateAccessToken(authCode)
			if err != nil {
				println("Error generating access token: ", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
				return
			}
			response := Response{
				Me:          user.FullUrl(),
				TokenType:   "Bearer",
				AccessToken: accessToken,
				Scope:       authCode.Scope,
				ExpiresIn:   duration,
			}
			jsonData, err := json.Marshal(response)
			if err != nil {
				println("Error marshalling json: ", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
			}
			w.Write(jsonData)
			return
		}
	}
}

func userAuthVerifyHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		// get form data from post request
		err = r.ParseForm()
		if err != nil {
			println("Error parsing form: ", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error parsing form"))
			return
		}
		password := r.FormValue("password")
		client_id := r.FormValue("client_id")
		redirect_uri := r.FormValue("redirect_uri")
		response_type := r.FormValue("response_type")
		state := r.FormValue("state")
		code_challenge := r.FormValue("code_challenge")
		code_challenge_method := r.FormValue("code_challenge_method")
		scope := r.FormValue("scope")

		// CSRF check
		formCsrfToken := r.FormValue("csrf_token")
		cookieCsrfToken, err := r.Cookie("csrf_token")

		if err != nil {
			println("Error getting csrf token from cookie: ", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error getting csrf token from cookie"))
			return
		}
		if formCsrfToken != cookieCsrfToken.Value {
			println("Invalid csrf token")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid csrf token"))
			return
		}

		password_valid := user.VerifyPassword(password)
		if !password_valid {
			redirect := fmt.Sprintf(
				"%s?error=invalid_password&client_id=%s&redirect_uri=%s&response_type=%s&state=%s",
				user.AuthUrl(), client_id, redirect_uri, response_type, state,
			)
			if code_challenge != "" {
				redirect += fmt.Sprintf("&code_challenge=%s&code_challenge_method=%s", code_challenge, code_challenge_method)
			}
			http.Redirect(w, r,
				redirect,
				http.StatusFound,
			)
			return
		} else {
			// password is valid, generate code
			code, err := user.GenerateAuthCode(
				client_id, redirect_uri, code_challenge, code_challenge_method, scope)
			if err != nil {
				println("Error generating code: ", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
				return
			}
			http.Redirect(w, r,
				fmt.Sprintf(
					"%s?code=%s&state=%s&iss=%s",
					redirect_uri, code, state,
					user.FullUrl(),
				),
				http.StatusFound,
			)
			return
		}

	}
}

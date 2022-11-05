package web

import (
	"encoding/json"
	"fmt"
	"h4kor/owl-blogs"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func getUserFromRepo(repo *owl.Repository, ps httprouter.Params) (owl.User, error) {
	if config, _ := repo.Config(); config.SingleUser != "" {
		return repo.GetUser(config.SingleUser)
	}
	userName := ps.ByName("user")
	user, err := repo.GetUser(userName)
	if err != nil {
		return owl.User{}, err
	}
	return user, nil
}

func repoIndexHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		html, err := owl.RenderUserList(*repo)

		if err != nil {
			println("Error rendering index: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		println("Rendering index")
		w.Write([]byte(html))
	}
}

func userIndexHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		html, err := owl.RenderIndexPage(user)
		if err != nil {
			println("Error rendering index page: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		println("Rendering index page for user", user.Name())
		w.Write([]byte(html))
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

		reqData := owl.AuthRequestData{
			Me:           me,
			ClientId:     clientId,
			RedirectUri:  redirectUri,
			State:        state,
			ResponseType: responseType,
			User:         user,
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

func userAuthProfileHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
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
		code := r.Form.Get("code")
		client_id := r.Form.Get("client_id")
		redirect_uri := r.Form.Get("redirect_uri")

		// check if request is valid
		valid := user.VerifyAuthCode(code, client_id, redirect_uri)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid code"))
			return
		} else {
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
		println("Password: ", password)
		client_id := r.FormValue("client_id")
		redirect_uri := r.FormValue("redirect_uri")
		response_type := r.FormValue("response_type")
		state := r.FormValue("state")

		password_valid := user.VerifyPassword(password)
		if !password_valid {
			http.Redirect(w, r,
				fmt.Sprintf(
					"%s?error=invalid_password&client_id=%s&redirect_uri=%s&response_type=%s&state=%s",
					user.AuthUrl(), client_id, redirect_uri, response_type, state,
				),
				http.StatusFound,
			)
			return
		} else {
			// password is valid, generate code
			code, err := user.GenerateAuthCode(client_id, redirect_uri)
			if err != nil {
				println("Error generating code: ", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
				return
			}
			http.Redirect(w, r,
				fmt.Sprintf(
					"%s?code=%s&state=%s",
					redirect_uri, code, state,
				),
				http.StatusFound,
			)
			return
		}

	}
}

func userWebmentionHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("User not found"))
			return
		}
		err = r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to parse form data"))
			return
		}
		params := r.PostForm
		target := params["target"]
		source := params["source"]
		if len(target) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No target provided"))
			return
		}
		if len(source) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No source provided"))
			return
		}

		if len(target[0]) < 7 || (target[0][:7] != "http://" && target[0][:8] != "https://") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Not a valid target"))
			return
		}

		if len(source[0]) < 7 || (source[0][:7] != "http://" && source[0][:8] != "https://") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Not a valid source"))
			return
		}

		if source[0] == target[0] {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("target and source are equal"))
			return
		}

		tryAlias := func(target string) *owl.Post {
			parsedTarget, _ := url.Parse(target)
			aliases, _ := repo.PostAliases()
			fmt.Printf("aliases %v", aliases)
			fmt.Printf("parsedTarget %v", parsedTarget)
			if _, ok := aliases[parsedTarget.Path]; ok {
				return aliases[parsedTarget.Path]
			}
			return nil
		}

		var aliasPost *owl.Post
		parts := strings.Split(target[0], "/")
		if len(parts) < 2 {
			aliasPost = tryAlias(target[0])
			if aliasPost == nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Not found"))
				return
			}
		}
		postId := parts[len(parts)-2]
		foundPost, err := user.GetPost(postId)
		if err != nil && aliasPost == nil {
			aliasPost = tryAlias(target[0])
			if aliasPost == nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Post not found"))
				return
			}
		}
		if aliasPost != nil {
			foundPost = aliasPost
		}
		err = foundPost.AddIncomingWebmention(source[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to process webmention"))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(""))
	}
}

func userRSSHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		xml, err := owl.RenderRSSFeed(user)
		if err != nil {
			println("Error rendering index page: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		println("Rendering index page for user", user.Name())
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(xml))
	}
}

func postHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		postId := ps.ByName("post")

		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		post, err := user.GetPost(postId)

		if err != nil {
			println("Error getting post: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		meta := post.Meta()
		if meta.Draft {
			println("Post is a draft")
			notFoundHandler(repo)(w, r)
			return
		}

		html, err := owl.RenderPost(post)
		if err != nil {
			println("Error rendering post: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		println("Rendering post", postId)
		w.Write([]byte(html))

	}
}

func postMediaHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		postId := ps.ByName("post")
		filepath := ps.ByName("filepath")

		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		post, err := user.GetPost(postId)
		if err != nil {
			println("Error getting post: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		filepath = path.Join(post.MediaDir(), filepath)
		if _, err := os.Stat(filepath); err != nil {
			println("Error getting file: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		http.ServeFile(w, r, filepath)
	}
}

func userMediaHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filepath := ps.ByName("filepath")

		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		filepath = path.Join(user.MediaDir(), filepath)
		if _, err := os.Stat(filepath); err != nil {
			println("Error getting file: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		http.ServeFile(w, r, filepath)
	}
}

func notFoundHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		aliases, _ := repo.PostAliases()
		if _, ok := aliases[path]; ok {
			http.Redirect(w, r, aliases[path].UrlPath(), http.StatusMovedPermanently)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}
}

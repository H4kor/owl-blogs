package web

import (
	"fmt"
	"h4kor/owl-blogs"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

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
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Internal server error",
				Message: "Internal server error",
			})
			w.Write([]byte(html))
			return
		}
		println("Rendering index page for user", user.Name())
		w.Write([]byte(html))
	}
}

func postListHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		listId := ps.ByName("list")
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		list, err := user.GetPostList(listId)

		if err != nil {
			println("Error getting post list: ", err.Error())
			notFoundUserHandler(repo, user)(w, r)
			return
		}

		html, err := owl.RenderPostList(user, list)
		if err != nil {
			println("Error rendering index page: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Internal server error",
				Message: "Internal server error",
			})
			w.Write([]byte(html))
			return
		}
		println("Rendering index page for user", user.Name())
		w.Write([]byte(html))
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

		tryAlias := func(target string) owl.Post {
			parsedTarget, _ := url.Parse(target)
			aliases, _ := repo.PostAliases()
			fmt.Printf("aliases %v", aliases)
			fmt.Printf("parsedTarget %v", parsedTarget)
			if _, ok := aliases[parsedTarget.Path]; ok {
				return aliases[parsedTarget.Path]
			}
			return nil
		}

		var aliasPost owl.Post
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
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Internal server error",
				Message: "Internal server error",
			})
			w.Write([]byte(html))
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
			notFoundUserHandler(repo, user)(w, r)
			return
		}

		meta := post.Meta()
		if meta.Draft {
			println("Post is a draft")
			notFoundUserHandler(repo, user)(w, r)
			return
		}

		html, err := owl.RenderPost(post)
		if err != nil {
			println("Error rendering post: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Internal server error",
				Message: "Internal server error",
			})
			w.Write([]byte(html))
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
			notFoundUserHandler(repo, user)(w, r)
			return
		}
		filepath = path.Join(post.MediaDir(), filepath)
		if _, err := os.Stat(filepath); err != nil {
			println("Error getting file: ", err.Error())
			notFoundUserHandler(repo, user)(w, r)
			return
		}
		http.ServeFile(w, r, filepath)
	}
}

func userMicropubHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		// parse request form
		err = r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad request"))
			return
		}

		// verify access token
		token := r.Header.Get("Authorization")
		if token == "" {
			token = r.Form.Get("access_token")
		} else {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		valid, _ := user.ValidateAccessToken(token)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		h := r.Form.Get("h")
		content := r.Form.Get("content")
		name := r.Form.Get("name")
		inReplyTo := r.Form.Get("in-reply-to")

		if h != "entry" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad request. h must be entry. Got " + h))
			return
		}
		if content == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad request. content is required"))
			return
		}

		// create post
		post, err := user.CreateNewPost(
			owl.PostMeta{
				Title: name,
				Reply: owl.ReplyData{
					Url: inReplyTo,
				},
				Date: time.Now(),
			},
			content,
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}

		w.Header().Add("Location", post.FullUrl())
		w.WriteHeader(http.StatusCreated)

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
			notFoundUserHandler(repo, user)(w, r)
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

func notFoundUserHandler(repo *owl.Repository, user owl.User) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		aliases, _ := repo.PostAliases()
		if _, ok := aliases[path]; ok {
			http.Redirect(w, r, aliases[path].UrlPath(), http.StatusMovedPermanently)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		html, _ := owl.RenderUserError(user, owl.ErrorMessage{
			Error:   "Not found",
			Message: "The page you requested could not be found",
		})
		w.Write([]byte(html))
	}
}

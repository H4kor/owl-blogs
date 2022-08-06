package main

import (
	"h4kor/owl-blogs"
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
)

func getUserFromRepo(repo *owl.Repository, ps httprouter.Params) (owl.User, error) {
	if repo.SingleUserName() != "" {
		return repo.GetUser(repo.SingleUserName())
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
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("User not found"))
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

func postHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		postId := ps.ByName("post")

		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("User not found"))
			return
		}
		post, err := user.GetPost(postId)
		if err != nil {
			println("Error getting post: ", err.Error())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Post not found"))
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
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("User not found"))
			return
		}
		post, err := user.GetPost(postId)
		if err != nil {
			println("Error getting post: ", err.Error())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Post not found"))
			return
		}
		filepath = path.Join(post.MediaDir(), filepath)
		if _, err := os.Stat(filepath); err != nil {
			println("Error getting file: ", err.Error())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("File not found"))
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

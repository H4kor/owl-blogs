package main

import (
	"h4kor/owl-blogs"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func indexHandler(repo owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		html, err := owl.RenderUserList(repo)

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

func userHandler(repo owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userName := ps.ByName("user")
		user, err := repo.GetUser(userName)
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
		println("Rendering index page for user", userName)
		w.Write([]byte(html))
	}
}

func postHandler(repo owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userName := ps.ByName("user")
		postId := ps.ByName("post")
		user, err := repo.GetUser(userName)
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

func postMediaHandler(repo owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userName := ps.ByName("user")
		postId := ps.ByName("post")
		filepath := ps.ByName("filepath")
		user, err := repo.GetUser(userName)

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
		http.ServeFile(w, r, filepath)
	}
}
func main() {
	println("owl web server")
	println("Parameters")
	println("-repo <repo> - Specify the repository to use. Defaults to '.'")
	println("-port <port> - Specify the port to use, Default is '8080'")
	var repoName string
	var port int
	for i, arg := range os.Args[0 : len(os.Args)-1] {
		if arg == "-port" {
			port, _ = strconv.Atoi(os.Args[i+1])
		}
		if arg == "-repo" {
			repoName = os.Args[i+1]
		}
	}
	if repoName == "" {
		repoName = "."
	}
	if port == 0 {
		port = 8080
	}

	repo, err := owl.OpenRepository(repoName)
	if err != nil {
		println("Error opening repository: ", err.Error())
		os.Exit(1)
	}

	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(repo.StaticDir()))
	router.GET("/", indexHandler(repo))
	router.GET("/user/:user/", userHandler(repo))
	router.GET("/user/:user/posts/:post/", postHandler(repo))
	router.GET("/user/:user/posts/:post/media/*filepath", postMediaHandler(repo))

	println("Listening on port", port)
	http.ListenAndServe(":"+strconv.Itoa(port), router)

}

package main

import (
	"h4kor/owl-blogs"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func Router(repo owl.Repository) http.Handler {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(repo.StaticDir()))
	router.GET("/", repoIndexHandler(repo))
	router.GET("/user/:user/", userIndexHandler(repo))
	router.GET("/user/:user/posts/:post/", postHandler(repo))
	router.GET("/user/:user/posts/:post/media/*filepath", postMediaHandler(repo))
	return router
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

	router := Router(repo)
	println("Listening on port", port)
	http.ListenAndServe(":"+strconv.Itoa(port), router)

}

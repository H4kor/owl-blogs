package main

import (
	"h4kor/owl-blogs"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func Router(repo *owl.Repository) http.Handler {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(repo.StaticDir()))
	router.GET("/", repoIndexHandler(repo))
	router.GET("/user/:user/", userIndexHandler(repo))
	router.POST("/user/:user/webmention/", userWebmentionHandler(repo))
	router.GET("/user/:user/index.xml", userRSSHandler(repo))
	router.GET("/user/:user/posts/:post/", postHandler(repo))
	router.GET("/user/:user/posts/:post/media/*filepath", postMediaHandler(repo))
	router.NotFound = http.HandlerFunc(notFoundHandler(repo))
	return router
}

func SingleUserRouter(repo *owl.Repository) http.Handler {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(repo.StaticDir()))
	router.GET("/", userIndexHandler(repo))
	router.POST("/webmention/", userWebmentionHandler(repo))
	router.GET("/index.xml", userRSSHandler(repo))
	router.GET("/posts/:post/", postHandler(repo))
	router.GET("/posts/:post/media/*filepath", postMediaHandler(repo))
	router.NotFound = http.HandlerFunc(notFoundHandler(repo))
	return router
}

func main() {
	println("owl web server")
	println("Parameters")
	println("-repo <repo> - Specify the repository to use. Defaults to '.'")
	println("-port <port> - Specify the port to use, Default is '8080'")
	println("-user <name> - Start server in single user mode.")
	println("-unsafe - Allow unsafe html.")
	var repoName string
	var port int
	var singleUserName string
	var allowRawHTML bool = false
	for i, arg := range os.Args[0:len(os.Args)] {
		if arg == "-port" {
			if i+1 >= len(os.Args) {
				println("-port requires a port number")
				os.Exit(1)
			}
			port, _ = strconv.Atoi(os.Args[i+1])
		}
		if arg == "-repo" {
			if i+1 >= len(os.Args) {
				println("-repo requires a repopath")
				os.Exit(1)
			}
			repoName = os.Args[i+1]
		}
		if arg == "-user" {
			if i+1 >= len(os.Args) {
				println("-user requires a username")
				os.Exit(1)
			}
			singleUserName = os.Args[i+1]
		}
		if arg == "-unsafe" {
			allowRawHTML = true
		}
	}
	if repoName == "" {
		repoName = "."
	}
	if port == 0 {
		port = 8080
	}

	var repo owl.Repository
	var err error
	if singleUserName != "" {
		println("Single user mode")
		println("Repository:", repoName)
		println("User:", singleUserName)
		repo, err = owl.OpenSingleUserRepo(repoName, singleUserName)
	} else {
		println("Multi user mode")
		println("Repository:", repoName)
		repo, err = owl.OpenRepository(repoName)
	}
	repo.SetAllowRawHtml(allowRawHTML)

	if err != nil {
		println("Error opening repository: ", err.Error())
		os.Exit(1)
	}

	var router http.Handler
	if singleUserName == "" {
		println("Multi user mode Router used")
		router = Router(&repo)
	} else {
		println("Single user mode Router used")
		router = SingleUserRouter(&repo)
	}
	println("Listening on port", port)
	http.ListenAndServe(":"+strconv.Itoa(port), router)

}

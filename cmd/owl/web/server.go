package web

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

func StartServer(repoPath string, port int) {
	var repo owl.Repository
	var err error
	repo, err = owl.OpenRepository(repoPath)

	if err != nil {
		println("Error opening repository: ", err.Error())
		os.Exit(1)
	}

	var router http.Handler
	if config, _ := repo.Config(); config.SingleUser != "" {
		router = SingleUserRouter(&repo)
	} else {
		router = Router(&repo)
	}

	println("Listening on port", port)
	http.ListenAndServe(":"+strconv.Itoa(port), router)

}

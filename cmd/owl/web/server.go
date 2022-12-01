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
	router.GET("/user/:user/lists/:list/", postListHandler(repo))
	// Editor
	router.GET("/user/:user/editor/auth/", userLoginGetHandler(repo))
	router.POST("/user/:user/editor/auth/", userLoginPostHandler(repo))
	router.GET("/user/:user/editor/", userEditorGetHandler(repo))
	router.POST("/user/:user/editor/", userEditorPostHandler(repo))
	// Media
	router.GET("/user/:user/media/*filepath", userMediaHandler(repo))
	// RSS
	router.GET("/user/:user/index.xml", userRSSHandler(repo))
	// Posts
	router.GET("/user/:user/posts/:post/", postHandler(repo))
	router.GET("/user/:user/posts/:post/media/*filepath", postMediaHandler(repo))
	// Webmention
	router.POST("/user/:user/webmention/", userWebmentionHandler(repo))
	// Micropub
	router.POST("/user/:user/micropub/", userMicropubHandler(repo))
	// IndieAuth
	router.GET("/user/:user/auth/", userAuthHandler(repo))
	router.POST("/user/:user/auth/", userAuthProfileHandler(repo))
	router.POST("/user/:user/auth/verify/", userAuthVerifyHandler(repo))
	router.POST("/user/:user/auth/token/", userAuthTokenHandler(repo))
	router.GET("/user/:user/.well-known/oauth-authorization-server", userAuthMetadataHandler(repo))
	router.NotFound = http.HandlerFunc(notFoundHandler(repo))
	return router
}

func SingleUserRouter(repo *owl.Repository) http.Handler {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(repo.StaticDir()))
	router.GET("/", userIndexHandler(repo))
	router.GET("/lists/:list/", postListHandler(repo))
	// Editor
	router.GET("/editor/auth/", userLoginGetHandler(repo))
	router.POST("/editor/auth/", userLoginPostHandler(repo))
	router.GET("/editor/", userEditorGetHandler(repo))
	router.POST("/editor/", userEditorPostHandler(repo))
	// Media
	router.GET("/media/*filepath", userMediaHandler(repo))
	// RSS
	router.GET("/index.xml", userRSSHandler(repo))
	// Posts
	router.GET("/posts/:post/", postHandler(repo))
	router.GET("/posts/:post/media/*filepath", postMediaHandler(repo))
	// Webmention
	router.POST("/webmention/", userWebmentionHandler(repo))
	// Micropub
	router.POST("/micropub/", userMicropubHandler(repo))
	// IndieAuth
	router.GET("/auth/", userAuthHandler(repo))
	router.POST("/auth/", userAuthProfileHandler(repo))
	router.POST("/auth/verify/", userAuthVerifyHandler(repo))
	router.POST("/auth/token/", userAuthTokenHandler(repo))
	router.GET("/.well-known/oauth-authorization-server", userAuthMetadataHandler(repo))
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

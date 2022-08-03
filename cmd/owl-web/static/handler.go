package static

import (
	"h4kor/owl-blogs"
	"net/http"
)

func StaticHandler(repo owl.Repository) http.Handler {
	return http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir(repo.StaticDir())),
	)
}

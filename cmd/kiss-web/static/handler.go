package static

import (
	"h4kor/kiss-social"
	"net/http"
)

func StaticHandler(repo kiss.Repository) http.Handler {
	return http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir(repo.StaticDir())),
	)
}

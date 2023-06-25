package web

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type WebApp struct {
	app          *fiber.App
	entryService *app.EntryService
}

func NewWebApp(entryService *app.EntryService) *WebApp {
	app := fiber.New()

	indexHandler := NewIndexHandler(entryService)
	listHandler := NewListHandler(entryService)
	entryHandler := NewEntryHandler(entryService)
	mediaHandler := NewMediaHandler(entryService)
	rssHandler := NewRSSHandler(entryService)
	loginHandler := NewLoginHandler(entryService)
	editorHandler := NewEditorHandler(entryService)

	// app.ServeFiles("/static/*filepath", http.Dir(repo.StaticDir()))
	app.Get("/", indexHandler.Handle)
	app.Get("/lists/:list/", listHandler.Handle)
	// Editor
	app.Get("/editor/auth/", loginHandler.HandleGet)
	app.Post("/editor/auth/", loginHandler.HandlePost)
	app.Get("/editor/", editorHandler.HandleGet)
	app.Post("/editor/", editorHandler.HandlePost)
	// Media
	app.Get("/media/*filepath", mediaHandler.Handle)
	// RSS
	app.Get("/index.xml", rssHandler.Handle)
	// Posts
	app.Get("/posts/:post/", entryHandler.Handle)
	app.Get("/posts/:post/media/*filepath", mediaHandler.Handle)
	// Webmention
	// app.Post("/webmention/", userWebmentionHandler(repo))
	// Micropub
	// app.Post("/micropub/", userMicropubHandler(repo))
	// IndieAuth
	// app.Get("/auth/", userAuthHandler(repo))
	// app.Post("/auth/", userAuthProfileHandler(repo))
	// app.Post("/auth/verify/", userAuthVerifyHandler(repo))
	// app.Post("/auth/token/", userAuthTokenHandler(repo))
	// app.Get("/.well-known/oauth-authorization-server", userAuthMetadataHandler(repo))
	// app.NotFound = http.HandlerFunc(notFoundHandler(repo))
	return &WebApp{app: app, entryService: entryService}
}

func (w *WebApp) Run() {
	w.app.Listen(":3000")
}

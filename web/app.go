package web

import (
	"embed"
	"net/http"
	"owl-blogs/app"
	"owl-blogs/web/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed static/*
var embedDirStatic embed.FS

type WebApp struct {
	FiberApp      *fiber.App
	EntryService  *app.EntryService
	BinaryService *app.BinaryService
	Registry      *app.EntryTypeRegistry
	AuthorService *app.AuthorService
}

func NewWebApp(
	entryService *app.EntryService,
	typeRegistry *app.EntryTypeRegistry,
	binService *app.BinaryService,
	authorService *app.AuthorService,
) *WebApp {
	app := fiber.New()

	indexHandler := NewIndexHandler(entryService)
	listHandler := NewListHandler(entryService)
	entryHandler := NewEntryHandler(entryService, typeRegistry, authorService)
	mediaHandler := NewMediaHandler(binService)
	rssHandler := NewRSSHandler(entryService)
	loginHandler := NewLoginHandler(authorService)
	editorListHandler := NewEditorListHandler(typeRegistry)
	editorHandler := NewEditorHandler(entryService, typeRegistry, binService)

	// Login
	app.Get("/auth/login", loginHandler.HandleGet)
	app.Post("/auth/login", loginHandler.HandlePost)

	// Editor
	editor := app.Group("/editor")
	editor.Use(middleware.NewAuthMiddleware(authorService).Handle)
	editor.Get("/", editorListHandler.Handle)
	editor.Get("/:editor/", editorHandler.HandleGet)
	editor.Post("/:editor/", editorHandler.HandlePost)

	// app.Static("/static/*filepath", http.Dir(repo.StaticDir()))
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "static",
		Browse:     false,
	}))
	app.Get("/", indexHandler.Handle)
	app.Get("/lists/:list/", listHandler.Handle)
	// Media
	app.Get("/media/+", mediaHandler.Handle)
	// RSS
	app.Get("/index.xml", rssHandler.Handle)
	// Posts
	app.Get("/posts/:post/", entryHandler.Handle)
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
	return &WebApp{
		FiberApp:      app,
		EntryService:  entryService,
		Registry:      typeRegistry,
		BinaryService: binService,
		AuthorService: authorService,
	}
}

func (w *WebApp) Run() {
	w.FiberApp.Listen(":3000")
}

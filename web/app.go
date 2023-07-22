package web

import (
	"embed"
	"net/http"
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/web/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed static/*
var embedDirStatic embed.FS

type WebApp struct {
	FiberApp       *fiber.App
	EntryService   *app.EntryService
	BinaryService  *app.BinaryService
	Registry       *app.EntryTypeRegistry
	AuthorService  *app.AuthorService
	SiteConfigRepo repository.ConfigRepository
}

func NewWebApp(
	entryService *app.EntryService,
	typeRegistry *app.EntryTypeRegistry,
	binService *app.BinaryService,
	authorService *app.AuthorService,
	configRepo repository.ConfigRepository,
) *WebApp {
	app := fiber.New()

	indexHandler := NewIndexHandler(entryService, configRepo)
	listHandler := NewListHandler(entryService, configRepo)
	entryHandler := NewEntryHandler(entryService, typeRegistry, authorService, configRepo)
	mediaHandler := NewMediaHandler(binService)
	rssHandler := NewRSSHandler(entryService, configRepo)
	loginHandler := NewLoginHandler(authorService, configRepo)
	editorListHandler := NewEditorListHandler(typeRegistry, configRepo)
	editorHandler := NewEditorHandler(entryService, typeRegistry, binService, configRepo)

	// Login
	app.Get("/auth/login", loginHandler.HandleGet)
	app.Post("/auth/login", loginHandler.HandlePost)

	// Editor
	editor := app.Group("/editor")
	editor.Use(middleware.NewAuthMiddleware(authorService).Handle)
	editor.Get("/", editorListHandler.Handle)
	editor.Get("/:editor/", editorHandler.HandleGet)
	editor.Post("/:editor/", editorHandler.HandlePost)

	// SiteConfig
	siteConfig := app.Group("/site-config")
	siteConfig.Use(middleware.NewAuthMiddleware(authorService).Handle)

	siteConfigHandler := NewSiteConfigHandler(configRepo)
	siteConfig.Get("/", siteConfigHandler.HandleGet)
	siteConfig.Post("/", siteConfigHandler.HandlePost)

	siteConfigMeHandler := NewSiteConfigMeHandler(configRepo)
	siteConfig.Get("/me", siteConfigMeHandler.HandleGet)
	siteConfig.Post("/me/create/", siteConfigMeHandler.HandleCreate)
	siteConfig.Post("/me/delete/", siteConfigMeHandler.HandleDelete)

	siteConfigListHandler := NewSiteConfigListHandler(configRepo, typeRegistry)
	siteConfig.Get("/lists", siteConfigListHandler.HandleGet)
	siteConfig.Post("/lists/create/", siteConfigListHandler.HandleCreate)
	siteConfig.Post("/lists/delete/", siteConfigListHandler.HandleDelete)

	siteConfigMenusHandler := NewSiteConfigMenusHandler(configRepo)
	siteConfig.Get("/menus", siteConfigMenusHandler.HandleGet)
	siteConfig.Post("/menus/create/", siteConfigMenusHandler.HandleCreate)
	siteConfig.Post("/menus/delete/", siteConfigMenusHandler.HandleDelete)

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

	// ActivityPub
	// activityPubServer := NewActivityPubServer(configRepo)
	// app.Get("/.well-known/webfinger", activityPubServer.HandleWebfinger)
	// app.Route("/activitypub", activityPubServer.Router)

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
		FiberApp:       app,
		EntryService:   entryService,
		Registry:       typeRegistry,
		BinaryService:  binService,
		AuthorService:  authorService,
		SiteConfigRepo: configRepo,
	}
}

func (w *WebApp) Run() {
	w.FiberApp.Listen(":3000")
}

package web

import (
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/config"
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
	configRegister *app.ConfigRegister,
	siteConfigService *app.SiteConfigService,
	webmentionService *app.WebmentionService,
	interactionRepo repository.InteractionRepository,
	apService *app.ActivityPubService,
) *WebApp {
	fiberApp := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50MB in bytes
	})
	fiberApp.Use(middleware.NewUserMiddleware(authorService).Handle)

	indexHandler := NewIndexHandler(entryService, siteConfigService)
	listHandler := NewListHandler(entryService, siteConfigService)
	entryHandler := NewEntryHandler(entryService, typeRegistry, authorService, configRepo, interactionRepo)
	mediaHandler := NewMediaHandler(binService)
	rssHandler := NewRSSHandler(entryService, siteConfigService)
	loginHandler := NewLoginHandler(authorService, configRepo)
	editorHandler := NewEditorHandler(entryService, typeRegistry, binService, configRepo)
	webmentionHandler := NewWebmentionHandler(webmentionService, configRepo)

	// Login
	fiberApp.Get("/auth/login", loginHandler.HandleGet)
	fiberApp.Post("/auth/login", loginHandler.HandlePost)

	// admin
	adminHandler := NewAdminHandler(configRepo, configRegister, typeRegistry)
	draftHandler := NewDraftHandler(entryService, siteConfigService)
	binaryManageHandler := NewBinaryManageHandler(configRepo, binService)
	adminInteractionHandler := NewAdminInteractionHandler(configRepo, interactionRepo)
	admin := fiberApp.Group("/admin")
	admin.Use(middleware.NewAuthMiddleware(authorService).Handle)
	admin.Get("/", adminHandler.Handle)
	admin.Get("/drafts/", draftHandler.Handle)
	admin.Get("/config/:config/", adminHandler.HandleConfigGet)
	admin.Post("/config/:config/", adminHandler.HandleConfigPost)
	admin.Get("/binaries/", binaryManageHandler.Handle)
	admin.Post("/binaries/new/", binaryManageHandler.HandleUpload)
	admin.Post("/binaries/delete", binaryManageHandler.HandleDelete)
	admin.Post("/interactions/:id/delete/", adminInteractionHandler.HandleDelete)
	admin.Get("/interactions/", adminInteractionHandler.HandleGet)

	adminApi := admin.Group("/api")
	adminApi.Post("/binaries", binaryManageHandler.HandleUploadApi)

	// Editor
	editor := fiberApp.Group("/editor")
	editor.Use(middleware.NewAuthMiddleware(authorService).Handle)
	editor.Get("/new/:editor/", editorHandler.HandleGetNew)
	editor.Post("/new/:editor/", editorHandler.HandlePostNew)
	editor.Get("/edit/:id/", editorHandler.HandleGetEdit)
	editor.Post("/edit/:id/", editorHandler.HandlePostEdit)
	editor.Post("/delete/:id/", editorHandler.HandlePostDelete)
	editor.Post("/unpublish/:id/", editorHandler.HandlePostUnpublish)

	// SiteConfig
	siteConfig := fiberApp.Group("/site-config")
	siteConfig.Use(middleware.NewAuthMiddleware(authorService).Handle)

	siteConfigHandler := NewSiteConfigHandler(siteConfigService)
	siteConfig.Get("/", siteConfigHandler.HandleGet)
	siteConfig.Post("/", siteConfigHandler.HandlePost)

	siteConfigMeHandler := NewSiteConfigMeHandler(siteConfigService)
	siteConfig.Get("/me", siteConfigMeHandler.HandleGet)
	siteConfig.Post("/me/create/", siteConfigMeHandler.HandleCreate)
	siteConfig.Post("/me/delete/", siteConfigMeHandler.HandleDelete)

	siteConfigListHandler := NewSiteConfigListHandler(siteConfigService, typeRegistry)
	siteConfig.Get("/lists", siteConfigListHandler.HandleGet)
	siteConfig.Post("/lists/create/", siteConfigListHandler.HandleCreate)
	siteConfig.Post("/lists/delete/", siteConfigListHandler.HandleDelete)

	siteConfigMenusHandler := NewSiteConfigMenusHandler(siteConfigService)
	siteConfig.Get("/menus", siteConfigMenusHandler.HandleGet)
	siteConfig.Post("/menus/create/", siteConfigMenusHandler.HandleCreate)
	siteConfig.Post("/menus/delete/", siteConfigMenusHandler.HandleDelete)

	fiberApp.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "static",
		Browse:     false,
	}))
	fiberApp.Get("/", indexHandler.Handle)
	fiberApp.Get("/lists/:list/", listHandler.Handle)
	// Media
	fiberApp.Get("/media/+", mediaHandler.Handle)
	// RSS
	fiberApp.Get("/index.xml", rssHandler.Handle)
	// Posts
	fiberApp.Get("/posts/:post/", entryHandler.Handle)
	// Webmention
	fiberApp.Post("/webmention/", webmentionHandler.Handle)
	// robots.txt
	fiberApp.Get("/robots.txt", func(c *fiber.Ctx) error {
		siteConfig, _ := siteConfigService.GetSiteConfig()
		sitemapUrl, _ := url.JoinPath(siteConfig.FullUrl, "/sitemap.xml")
		c.Set("Content-Type", "text/plain")
		return c.SendString(fmt.Sprintf("User-agent: GPTBot\nDisallow: /\n\nUser-agent: *\nAllow: /\n\nSitemap: %s\n", sitemapUrl))
	})
	// sitemap.xml
	fiberApp.Get("/sitemap.xml", NewSiteMapHandler(entryService, siteConfigService).Handle)

	// ActivityPub
	activityPubServer := NewActivityPubServer(configRepo, entryService, apService)
	configRegister.Register(config.ACT_PUB_CONF_NAME, &app.ActivityPubConfig{})
	fiberApp.Get("/.well-known/webfinger", activityPubServer.HandleWebfinger)
	fiberApp.Route("/activitypub", activityPubServer.Router)

	return &WebApp{
		FiberApp:       fiberApp,
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

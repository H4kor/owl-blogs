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
	thumbnailService *app.ThumbnailService,
	authorService *app.AuthorService,
	configRepo repository.ConfigRepository,
	configRegister *app.ConfigRegister,
	siteConfigService *app.SiteConfigService,
	webmentionService *app.WebmentionService,
	interactionRepo repository.InteractionRepository,
	apService *app.ActivityPubService,
) *WebApp {
	fiberApp := fiber.New(fiber.Config{
		BodyLimit:             50 * 1024 * 1024, // 50MB in bytes
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if werr, ok := err.(app.WebError); ok {
				c.SendStatus(werr.Status())
				return c.Send([]byte(werr.Error()))
			}
			return fiber.DefaultErrorHandler(c, err)
		},
	})
	fiberApp.Use(middleware.NewUserMiddleware(authorService).Handle)

	indexHandler := NewIndexHandler(entryService, siteConfigService)
	listHandler := NewListHandler(entryService, siteConfigService)
	tagHandler := NewTagHandler(entryService, siteConfigService)
	entryHandler := NewEntryHandler(entryService, typeRegistry, authorService, configRepo, interactionRepo)
	mediaHandler := NewMediaHandler(binService)
	thumbnailHandler := NewThumbnailHandler(binService, thumbnailService)
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
	admin.Post("/interactions/delete/", adminInteractionHandler.HandleDelete)
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

	siteConfigHandler := NewSiteConfigHandler(siteConfigService, typeRegistry)
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

	activityPubServer := NewActivityPubServer(siteConfigService, entryService, apService)
	configRegister.Register(config.ACT_PUB_CONF_NAME, &app.ActivityPubConfig{})

	fiberApp.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "static",
		Browse:     false,
	}))
	fiberApp.Get("/", activityPubServer.HandleActor, indexHandler.Handle)
	// Posts
	fiberApp.Get("/posts/:post/", activityPubServer.HandleEntry, entryHandler.Handle)
	// Tags
	fiberApp.Get("/tags/", tagHandler.HandleList)
	fiberApp.Get("/tags/:tag/", tagHandler.Handle)
	// Lists
	fiberApp.Get("/lists/:list/", listHandler.Handle)
	// Media
	fiberApp.Get("/media/+", mediaHandler.Handle)
	fiberApp.Get("/thumbnail/+", thumbnailHandler.Handle)
	// RSS
	fiberApp.Get("/index.xml", rssHandler.Handle)
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

func (w *WebApp) Run(bindAddr string) {
	fmt.Printf("owl-blogs starts listning on: %s\n", bindAddr)
	err := w.FiberApp.Listen(bindAddr)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

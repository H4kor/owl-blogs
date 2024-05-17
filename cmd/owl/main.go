package main

import (
	"fmt"
	"os"
	"owl-blogs/app"
	entrytypes "owl-blogs/entry_types"
	"owl-blogs/infra"
	"owl-blogs/interactions"
	"owl-blogs/plugings"
	"owl-blogs/render"
	"owl-blogs/web"

	"github.com/spf13/cobra"
)

const DbPath = "owlblogs.db"

var rootCmd = &cobra.Command{
	Use:   "owl",
	Short: "Owl Blogs is a not so static blog generator",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func App(db infra.Database) *web.WebApp {
	// Register Types
	entryRegister := app.NewEntryTypeRegistry()
	entryRegister.Register(&entrytypes.Image{})
	entryRegister.Register(&entrytypes.Article{})
	entryRegister.Register(&entrytypes.Page{})
	entryRegister.Register(&entrytypes.Recipe{})
	entryRegister.Register(&entrytypes.Note{})
	entryRegister.Register(&entrytypes.Bookmark{})
	entryRegister.Register(&entrytypes.Reply{})

	interactionRegister := app.NewInteractionTypeRegistry()
	interactionRegister.Register(&interactions.Webmention{})

	configRegister := app.NewConfigRegister()

	// Create Repositories
	entryRepo := infra.NewEntryRepository(db, entryRegister)
	binRepo := infra.NewBinaryFileRepo(db)
	authorRepo := infra.NewDefaultAuthorRepo(db)
	configRepo := infra.NewConfigRepo(db)
	interactionRepo := infra.NewInteractionRepo(db, interactionRegister)
	followersRepo := infra.NewFollowerRepository(db)

	// Create External Services
	httpClient := &infra.OwlHttpClient{}

	// busses
	eventBus := app.NewEventBus()

	// Create Services
	siteConfigService := app.NewSiteConfigService(configRepo)
	entryService := app.NewEntryService(entryRepo, eventBus)
	binaryService := app.NewBinaryFileService(binRepo)
	authorService := app.NewAuthorService(authorRepo, siteConfigService)
	webmentionService := app.NewWebmentionService(
		siteConfigService, interactionRepo, entryRepo, httpClient, eventBus,
	)
	apService := app.NewActivityPubService(followersRepo, configRepo, siteConfigService)

	// setup render functions
	render.SiteConfigService = siteConfigService

	// plugins
	plugings.NewEcho(eventBus)
	plugings.RegisterInstagram(
		configRepo, configRegister, binaryService, eventBus,
	)

	// Create WebApp
	return web.NewWebApp(
		entryService, entryRegister, binaryService,
		authorService, configRepo, configRegister,
		siteConfigService, webmentionService, interactionRepo,
		apService,
	)

}

func main() {
	Execute()
}

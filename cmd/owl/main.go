package main

import (
	"fmt"
	"os"
	"owl-blogs/app"
	entrytypes "owl-blogs/entry_types"
	"owl-blogs/infra"
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
	registry := app.NewEntryTypeRegistry()
	registry.Register(&entrytypes.Image{})
	registry.Register(&entrytypes.Article{})
	registry.Register(&entrytypes.Page{})
	registry.Register(&entrytypes.Recipe{})
	registry.Register(&entrytypes.Note{})

	entryRepo := infra.NewEntryRepository(db, registry)
	binRepo := infra.NewBinaryFileRepo(db)
	authorRepo := infra.NewDefaultAuthorRepo(db)
	siteConfigRepo := infra.NewSiteConfigRepo(db)

	entryService := app.NewEntryService(entryRepo)
	binaryService := app.NewBinaryFileService(binRepo)
	authorService := app.NewAuthorService(authorRepo, siteConfigRepo)

	return web.NewWebApp(entryService, registry, binaryService, authorService)

}

func main() {
	Execute()
}

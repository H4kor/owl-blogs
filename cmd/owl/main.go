package main

import (
	"fmt"
	"os"
	"owl-blogs/app"
	"owl-blogs/config"
	"owl-blogs/domain/model"
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
	config := config.NewConfig()

	registry := app.NewEntryTypeRegistry()
	registry.Register(&model.Image{})
	registry.Register(&model.Article{})
	registry.Register(&model.Page{})
	registry.Register(&model.Recipe{})
	registry.Register(&model.Note{})

	entryRepo := infra.NewEntryRepository(db, registry)
	binRepo := infra.NewBinaryFileRepo(db)
	authorRepo := infra.NewDefaultAuthorRepo(db)

	entryService := app.NewEntryService(entryRepo)
	binaryService := app.NewBinaryFileService(binRepo)
	authorService := app.NewAuthorService(authorRepo, config)

	return web.NewWebApp(entryService, registry, binaryService, authorService)

}

func main() {
	Execute()
}

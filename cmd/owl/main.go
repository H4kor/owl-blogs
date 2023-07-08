package main

import (
	"owl-blogs/app"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"owl-blogs/web"
)

const DbPath = "owlblogs.db"

func App(db infra.Database) *web.WebApp {
	config := config.NewConfig()

	registry := app.NewEntryTypeRegistry()
	registry.Register(&model.ImageEntry{})

	entryRepo := infra.NewEntryRepository(db, registry)
	binRepo := infra.NewBinaryFileRepo(db)
	authorRepo := infra.NewDefaultAuthorRepo(db)

	entryService := app.NewEntryService(entryRepo)
	binaryService := app.NewBinaryFileService(binRepo)
	authorService := app.NewAuthorService(authorRepo, config)

	return web.NewWebApp(entryService, registry, binaryService, authorService)

}

func main() {
	db := infra.NewSqliteDB(DbPath)
	App(db).Run()
}

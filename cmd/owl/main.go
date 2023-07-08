package main

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"owl-blogs/web"
)

const DbPath = "owlblogs.db"

func App(db infra.Database) *web.WebApp {
	registry := app.NewEntryTypeRegistry()

	registry.Register(&model.ImageEntry{})

	entryRepo := infra.NewEntryRepository(db, registry)
	binRepo := infra.NewBinaryFileRepo(db)

	entryService := app.NewEntryService(entryRepo)
	binaryService := app.NewBinaryFileService(binRepo)
	return web.NewWebApp(entryService, registry, binaryService)

}

func main() {
	db := infra.NewSqliteDB(DbPath)
	App(db).Run()
}

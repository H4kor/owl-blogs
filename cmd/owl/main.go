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

	repo := infra.NewEntryRepository(db, registry)
	entryService := app.NewEntryService(repo)
	return web.NewWebApp(entryService, registry)

}

func main() {
	db := infra.NewSqliteDB(DbPath)
	App(db).Run()
}

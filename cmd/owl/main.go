package main

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"owl-blogs/web"
)

func App() *web.WebApp {
	db := infra.NewSqliteDB("owlblogs.db")
	registry := app.NewEntryTypeRegistry()

	registry.Register(&model.ImageEntry{})

	repo := infra.NewEntryRepository(db, registry)
	entryService := app.NewEntryService(repo)
	return web.NewWebApp(entryService, registry)

}

func main() {
	App().Run()
}

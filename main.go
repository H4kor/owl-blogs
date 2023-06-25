package main

import (
	"owl-blogs/app"
	"owl-blogs/infra"
	"owl-blogs/web"
)

func main() {
	db := infra.NewSqliteDB("owlblogs.db")
	registry := app.NewEntryTypeRegistry()
	repo := infra.NewEntryRepository(db, registry)
	entryService := app.NewEntryService(repo)
	webApp := web.NewWebApp(entryService)
	webApp.Run()
}

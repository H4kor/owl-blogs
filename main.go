package main

import (
	"owl-blogs/app"
	"owl-blogs/infra"
	"owl-blogs/web"
)

func main() {
	db := infra.NewSqliteDB("owlblogs.db")
	repo := infra.NewEntryRepository(db)
	entryService := app.NewEntryService(repo)
	webApp := web.NewWebApp(entryService)
	webApp.Run()
}

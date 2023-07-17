package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
)

func getConfig(repo repository.SiteConfigRepository) model.SiteConfig {
	config, err := repo.Get()
	if err != nil {
		panic(err)
	}
	return config
}

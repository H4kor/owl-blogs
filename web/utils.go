package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"
)

func getSiteConfig(repo repository.ConfigRepository) model.SiteConfig {
	siteConfig := model.SiteConfig{}
	err := repo.Get(config.SITE_CONFIG, &siteConfig)
	if err != nil {
		panic(err)
	}
	return siteConfig
}

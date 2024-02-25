package app

import (
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"reflect"
)

// SiteConfigService is a service to retrieve and store the site config
// Even though the site config is a standard config, it is handle by an extra service
// as it is used in many places.
// The SiteConfig contains global settings require by multiple parts of the app
type SiteConfigService struct {
	repo repository.ConfigRepository
}

func NewSiteConfigService(repo repository.ConfigRepository) *SiteConfigService {
	return &SiteConfigService{
		repo: repo,
	}
}

func (svc *SiteConfigService) defaultConfig() model.SiteConfig {
	return model.SiteConfig{
		Title:              "My Owl-Blog",
		SubTitle:           "A freshly created blog",
		HeaderColor:        "#efc48c",
		PrimaryColor:       "#d37f12",
		AuthorName:         "",
		Me:                 []model.MeLinks{},
		Lists:              []model.EntryList{},
		PrimaryListInclude: []string{},
		HeaderMenu:         []model.MenuItem{},
		FooterMenu:         []model.MenuItem{},
		Secret:             "",
		AvatarUrl:          "",
		FullUrl:            "http://localhost:3000",
		HtmlHeadExtra:      "",
		FooterExtra:        "",
	}
}

func (svc *SiteConfigService) GetSiteConfig() (model.SiteConfig, error) {
	siteConfig := model.SiteConfig{}
	err := svc.repo.Get(config.SITE_CONFIG, &siteConfig)
	if err != nil {
		println("ERROR IN SITE CONFIG")
		return model.SiteConfig{}, err
	}
	if reflect.ValueOf(siteConfig).IsZero() {
		return svc.defaultConfig(), nil
	}
	return siteConfig, nil
}

func (svc *SiteConfigService) UpdateSiteConfig(cfg model.SiteConfig) error {
	return svc.repo.Update(config.SITE_CONFIG, cfg)
}

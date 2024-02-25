package model

type SiteConfigInterface interface {
	GetSiteConfig() (SiteConfig, error)
	UpdateSiteConfig(cfg SiteConfig) error
}

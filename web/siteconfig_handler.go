package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigHandler struct {
	siteConfigRepo repository.ConfigRepository
}

func NewSiteConfigHandler(siteConfigRepo repository.ConfigRepository) *SiteConfigHandler {
	return &SiteConfigHandler{
		siteConfigRepo: siteConfigRepo,
	}
}

func (h *SiteConfigHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := model.SiteConfig{}
	err := h.siteConfigRepo.Get(config.SITE_CONFIG, &siteConfig)
	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(c, getSiteConfig(h.siteConfigRepo), "views/site_config", siteConfig)
}

func (h *SiteConfigHandler) HandlePost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := model.SiteConfig{}
	err := h.siteConfigRepo.Get(config.SITE_CONFIG, &siteConfig)

	if err != nil {
		return err
	}

	siteConfig.Title = c.FormValue("Title")
	siteConfig.SubTitle = c.FormValue("SubTitle")
	siteConfig.HeaderColor = c.FormValue("HeaderColor")
	siteConfig.AuthorName = c.FormValue("AuthorName")
	siteConfig.AvatarUrl = c.FormValue("AvatarUrl")
	siteConfig.FullUrl = c.FormValue("FullUrl")

	err = h.siteConfigRepo.Update(config.SITE_CONFIG, siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/")
}

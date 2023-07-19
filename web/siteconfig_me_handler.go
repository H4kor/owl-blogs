package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigMeHandler struct {
	siteConfigRepo repository.ConfigRepository
}

func NewSiteConfigMeHandler(siteConfigRepo repository.ConfigRepository) *SiteConfigMeHandler {
	return &SiteConfigMeHandler{
		siteConfigRepo: siteConfigRepo,
	}
}

func (h *SiteConfigMeHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := model.SiteConfig{}
	err := h.siteConfigRepo.Get(config.SITE_CONFIG, &siteConfig)

	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(
		c, getSiteConfig(h.siteConfigRepo), "views/site_config_me", siteConfig.Me)
}

func (h *SiteConfigMeHandler) HandleCreate(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := model.SiteConfig{}
	err := h.siteConfigRepo.Get(config.SITE_CONFIG, &siteConfig)

	if err != nil {
		return err
	}

	siteConfig.Me = append(siteConfig.Me, model.MeLinks{
		Name: c.FormValue("Name"),
		Url:  c.FormValue("Url"),
	})

	err = h.siteConfigRepo.Update(config.SITE_CONFIG, siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/me")
}

func (h *SiteConfigMeHandler) HandleDelete(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := model.SiteConfig{}
	err := h.siteConfigRepo.Get(config.SITE_CONFIG, &siteConfig)

	if err != nil {
		return err
	}

	idx, err := strconv.Atoi(c.FormValue("idx"))
	if err != nil {
		return err
	}
	siteConfig.Me = append(siteConfig.Me[:idx], siteConfig.Me[idx+1:]...)

	err = h.siteConfigRepo.Update(config.SITE_CONFIG, siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/me")
}

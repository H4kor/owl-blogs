package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigMeHandler struct {
	siteConfigService *app.SiteConfigService
}

func NewSiteConfigMeHandler(siteConfigService *app.SiteConfigService) *SiteConfigMeHandler {
	return &SiteConfigMeHandler{
		siteConfigService: siteConfigService,
	}
}

func (h *SiteConfigMeHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()
	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(
		c, "views/site_config_me", siteConfig.Me)
}

func (h *SiteConfigMeHandler) HandleCreate(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()

	if err != nil {
		return err
	}

	siteConfig.Me = append(siteConfig.Me, model.MeLinks{
		Name: c.FormValue("Name"),
		Url:  c.FormValue("Url"),
	})

	err = h.siteConfigService.UpdateSiteConfig(siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/me")
}

func (h *SiteConfigMeHandler) HandleDelete(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()

	if err != nil {
		return err
	}

	idx, err := strconv.Atoi(c.FormValue("idx"))
	if err != nil {
		return err
	}
	siteConfig.Me = append(siteConfig.Me[:idx], siteConfig.Me[idx+1:]...)

	err = h.siteConfigService.UpdateSiteConfig(siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/me")
}

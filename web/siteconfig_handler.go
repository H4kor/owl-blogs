package web

import (
	"owl-blogs/app"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigHandler struct {
	svc *app.SiteConfigService
}

func NewSiteConfigHandler(svc *app.SiteConfigService) *SiteConfigHandler {
	return &SiteConfigHandler{
		svc: svc,
	}
}

func (h *SiteConfigHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.svc.GetSiteConfig()
	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(c, "views/site_config", siteConfig)
}

func (h *SiteConfigHandler) HandlePost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.svc.GetSiteConfig()
	if err != nil {
		return err
	}

	siteConfig.Title = c.FormValue("Title")
	siteConfig.SubTitle = c.FormValue("SubTitle")
	siteConfig.PrimaryColor = c.FormValue("PrimaryColor")
	siteConfig.AuthorName = c.FormValue("AuthorName")
	siteConfig.AvatarUrl = c.FormValue("AvatarUrl")
	siteConfig.FullUrl = c.FormValue("FullUrl")
	siteConfig.HtmlHeadExtra = c.FormValue("HtmlHeadExtra")
	siteConfig.FooterExtra = c.FormValue("FooterExtra")

	err = h.svc.UpdateSiteConfig(siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/")
}

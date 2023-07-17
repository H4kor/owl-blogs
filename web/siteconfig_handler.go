package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigHandler struct {
	siteConfigRepo repository.SiteConfigRepository
}

func NewSiteConfigHandler(siteConfigRepo repository.SiteConfigRepository) *SiteConfigHandler {
	return &SiteConfigHandler{
		siteConfigRepo: siteConfigRepo,
	}
}

func (h *SiteConfigHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(c, getConfig(h.siteConfigRepo), "views/site_config", config)
}

func (h *SiteConfigHandler) HandlePost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
	if err != nil {
		return err
	}

	config.Title = c.FormValue("Title")
	config.SubTitle = c.FormValue("SubTitle")
	config.HeaderColor = c.FormValue("HeaderColor")
	config.AuthorName = c.FormValue("AuthorName")
	config.AvatarUrl = c.FormValue("AvatarUrl")

	err = h.siteConfigRepo.Update(config)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/")
}

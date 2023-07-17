package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigMeHandler struct {
	siteConfigRepo repository.SiteConfigRepository
}

func NewSiteConfigMeHandler(siteConfigRepo repository.SiteConfigRepository) *SiteConfigMeHandler {
	return &SiteConfigMeHandler{
		siteConfigRepo: siteConfigRepo,
	}
}

func (h *SiteConfigMeHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(
		c, getConfig(h.siteConfigRepo), "views/site_config_me", config.Me)
}

func (h *SiteConfigMeHandler) HandleCreate(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
	if err != nil {
		return err
	}

	config.Me = append(config.Me, model.MeLinks{
		Name: c.FormValue("Name"),
		Url:  c.FormValue("Url"),
	})

	err = h.siteConfigRepo.Update(config)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/me")
}

func (h *SiteConfigMeHandler) HandleDelete(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
	if err != nil {
		return err
	}

	idx, err := strconv.Atoi(c.FormValue("idx"))
	if err != nil {
		return err
	}
	config.Me = append(config.Me[:idx], config.Me[idx+1:]...)

	err = h.siteConfigRepo.Update(config)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/me")
}

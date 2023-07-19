package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigMenusHandler struct {
	siteConfigRepo repository.ConfigRepository
}

type siteConfigMenusTemplateData struct {
	HeaderMenu []model.MenuItem
	FooterMenu []model.MenuItem
}

func NewSiteConfigMenusHandler(siteConfigRepo repository.ConfigRepository) *SiteConfigMenusHandler {
	return &SiteConfigMenusHandler{
		siteConfigRepo: siteConfigRepo,
	}
}

func (h *SiteConfigMenusHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := model.SiteConfig{}
	err := h.siteConfigRepo.Get(config.SITE_CONFIG, &siteConfig)

	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(
		c, getSiteConfig(h.siteConfigRepo), "views/site_config_menus", siteConfigMenusTemplateData{
			HeaderMenu: siteConfig.HeaderMenu,
			FooterMenu: siteConfig.FooterMenu,
		})
}

func (h *SiteConfigMenusHandler) HandleCreate(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := model.SiteConfig{}
	err := h.siteConfigRepo.Get(config.SITE_CONFIG, &siteConfig)

	if err != nil {
		return err
	}

	menuItem := model.MenuItem{
		Title: c.FormValue("Title"),
		List:  c.FormValue("List"),
		Url:   c.FormValue("Url"),
		Post:  c.FormValue("Post"),
	}

	if c.FormValue("menu") == "header" {
		siteConfig.HeaderMenu = append(siteConfig.HeaderMenu, menuItem)
	} else if c.FormValue("menu") == "footer" {
		siteConfig.FooterMenu = append(siteConfig.FooterMenu, menuItem)
	}

	err = h.siteConfigRepo.Update(config.SITE_CONFIG, siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/menus")
}

func (h *SiteConfigMenusHandler) HandleDelete(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := model.SiteConfig{}
	err := h.siteConfigRepo.Get(config.SITE_CONFIG, &siteConfig)

	if err != nil {
		return err
	}

	menu := c.FormValue("menu")
	idx, err := strconv.Atoi(c.FormValue("idx"))
	if err != nil {
		return err
	}

	if menu == "header" {
		siteConfig.HeaderMenu = append(siteConfig.HeaderMenu[:idx], siteConfig.HeaderMenu[idx+1:]...)
	} else if menu == "footer" {
		siteConfig.FooterMenu = append(siteConfig.FooterMenu[:idx], siteConfig.FooterMenu[idx+1:]...)
	}

	err = h.siteConfigRepo.Update(config.SITE_CONFIG, siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/menus")
}

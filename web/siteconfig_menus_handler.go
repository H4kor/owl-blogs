package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigMenusHandler struct {
	siteConfigRepo repository.SiteConfigRepository
}

type siteConfigMenusTemplateData struct {
	HeaderMenu []model.MenuItem
	FooterMenu []model.MenuItem
}

func NewSiteConfigMenusHandler(siteConfigRepo repository.SiteConfigRepository) *SiteConfigMenusHandler {
	return &SiteConfigMenusHandler{
		siteConfigRepo: siteConfigRepo,
	}
}

func (h *SiteConfigMenusHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(
		c, getConfig(h.siteConfigRepo), "views/site_config_menus", siteConfigMenusTemplateData{
			HeaderMenu: config.HeaderMenu,
			FooterMenu: config.FooterMenu,
		})
}

func (h *SiteConfigMenusHandler) HandleCreate(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
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
		config.HeaderMenu = append(config.HeaderMenu, menuItem)
	} else if c.FormValue("menu") == "footer" {
		config.FooterMenu = append(config.FooterMenu, menuItem)
	}

	err = h.siteConfigRepo.Update(config)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/menus")
}

func (h *SiteConfigMenusHandler) HandleDelete(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
	if err != nil {
		return err
	}

	menu := c.FormValue("menu")
	idx, err := strconv.Atoi(c.FormValue("idx"))
	if err != nil {
		return err
	}

	if menu == "header" {
		config.HeaderMenu = append(config.HeaderMenu[:idx], config.HeaderMenu[idx+1:]...)
	} else if menu == "footer" {
		config.FooterMenu = append(config.FooterMenu[:idx], config.FooterMenu[idx+1:]...)
	}

	err = h.siteConfigRepo.Update(config)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/menus")
}

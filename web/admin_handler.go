package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/render"
	"owl-blogs/web/forms"

	"github.com/gofiber/fiber/v2"
)

type adminHandler struct {
	configRepo     repository.ConfigRepository
	configRegister *app.ConfigRegister
}

func NewAdminHandler(configRepo repository.ConfigRepository, configRegister *app.ConfigRegister) *adminHandler {
	return &adminHandler{
		configRepo:     configRepo,
		configRegister: configRegister,
	}
}

func (h *adminHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := getSiteConfig(h.configRepo)
	configs := h.configRegister.Configs()
	return render.RenderTemplateWithBase(c, siteConfig, "views/admin", configs)
}

func (h *adminHandler) HandleConfigGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	configName := c.Params("config")
	config := h.configRegister.GetConfig(configName)
	if config == nil {
		return c.SendStatus(404)
	}
	err := h.configRepo.Get(configName, config)
	if err != nil {
		return err
	}
	siteConfig := getSiteConfig(h.configRepo)

	form := forms.NewForm(config, nil)
	htmlForm, err := form.HtmlForm()
	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(c, siteConfig, "views/admin_config", htmlForm)
}

func (h *adminHandler) HandleConfigPost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	configName := c.Params("config")
	config := h.configRegister.GetConfig(configName)
	if config == nil {
		return c.SendStatus(404)
	}

	form := forms.NewForm(config, nil)

	newConfig, err := form.Parse(c)
	if err != nil {
		return err
	}

	h.configRepo.Update(configName, newConfig)

	return c.Redirect("")

}

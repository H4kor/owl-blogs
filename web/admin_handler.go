package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/render"
	"owl-blogs/web/forms"
	"sort"

	"github.com/gofiber/fiber/v2"
)

type adminHandler struct {
	configRepo     repository.ConfigRepository
	configRegister *app.ConfigRegister
	typeRegistry   *app.EntryTypeRegistry
}

type adminContet struct {
	Configs []app.RegisteredConfig
	Types   []string
}

func NewAdminHandler(
	configRepo repository.ConfigRepository,
	configRegister *app.ConfigRegister,
	typeRegistry *app.EntryTypeRegistry,
) *adminHandler {
	return &adminHandler{
		configRepo:     configRepo,
		configRegister: configRegister,
		typeRegistry:   typeRegistry,
	}
}

func (h *adminHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig := getSiteConfig(h.configRepo)
	configs := h.configRegister.Configs()

	types := h.typeRegistry.Types()
	typeNames := []string{}

	for _, t := range types {
		name, _ := h.typeRegistry.TypeName(t)
		typeNames = append(typeNames, name)
	}

	// sort names to have a consistent order
	sort.Slice(typeNames, func(i, j int) bool {
		return typeNames[i] < typeNames[j]
	})

	return render.RenderTemplateWithBase(
		c, siteConfig,
		"views/admin", &adminContet{
			Configs: configs,
			Types:   typeNames,
		},
	)
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

package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigListHandler struct {
	siteConfigService *app.SiteConfigService
	typeRegistry      *app.EntryTypeRegistry
}

type siteConfigListTemplateData struct {
	Lists []model.EntryList
	Types []string
}

func NewSiteConfigListHandler(
	siteConfigService *app.SiteConfigService,
	typeRegistry *app.EntryTypeRegistry,
) *SiteConfigListHandler {
	return &SiteConfigListHandler{
		siteConfigService: siteConfigService,
		typeRegistry:      typeRegistry,
	}
}

func (h *SiteConfigListHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()
	if err != nil {
		return err
	}

	types := make([]string, 0)
	for _, t := range h.typeRegistry.Types() {
		typeName, err := h.typeRegistry.TypeName(t)
		if err != nil {
			continue
		}
		types = append(types, typeName)
	}

	return render.RenderTemplateWithBase(
		c, "views/site_config_list", siteConfigListTemplateData{
			Lists: siteConfig.Lists,
			Types: types,
		})
}

func (h *SiteConfigListHandler) HandleCreate(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()

	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	siteConfig.Lists = append(siteConfig.Lists, model.EntryList{
		Id:       c.FormValue("Id"),
		Title:    c.FormValue("Title"),
		Include:  form.Value["Include"],
		ListType: c.FormValue("ListType"),
	})

	err = h.siteConfigService.UpdateSiteConfig(siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/lists")
}

func (h *SiteConfigListHandler) HandleDelete(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()

	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.FormValue("idx"))
	if err != nil {
		return err
	}

	siteConfig.Lists = append(siteConfig.Lists[:id], siteConfig.Lists[id+1:]...)

	err = h.siteConfigService.UpdateSiteConfig(siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/lists")
}

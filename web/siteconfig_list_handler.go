package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigListHandler struct {
	siteConfigRepo repository.SiteConfigRepository
	typeRegistry   *app.EntryTypeRegistry
}

type siteConfigListTemplateData struct {
	Lists []model.EntryList
	Types []string
}

func NewSiteConfigListHandler(
	siteConfigRepo repository.SiteConfigRepository,
	typeRegistry *app.EntryTypeRegistry,
) *SiteConfigListHandler {
	return &SiteConfigListHandler{
		siteConfigRepo: siteConfigRepo,
		typeRegistry:   typeRegistry,
	}
}

func (h *SiteConfigListHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
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
		c, getConfig(h.siteConfigRepo), "views/site_config_list", siteConfigListTemplateData{
			Lists: config.Lists,
			Types: types,
		})
}

func (h *SiteConfigListHandler) HandleCreate(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	config.Lists = append(config.Lists, model.EntryList{
		Id:       c.FormValue("Id"),
		Title:    c.FormValue("Title"),
		Include:  form.Value["Include"],
		ListType: c.FormValue("ListType"),
	})

	err = h.siteConfigRepo.Update(config)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/lists")
}

func (h *SiteConfigListHandler) HandleDelete(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	config, err := h.siteConfigRepo.Get()
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.FormValue("idx"))
	if err != nil {
		return err
	}

	config.Lists = append(config.Lists[:id], config.Lists[id+1:]...)

	err = h.siteConfigRepo.Update(config)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/lists")
}

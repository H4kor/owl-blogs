package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/internal"
	"owl-blogs/render"
	"slices"
	"sort"
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

	sort.Slice(types, func(i, j int) bool {
		return types[i] < types[j]
	})

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

	id := c.FormValue("Id")
	title := c.FormValue("Title")
	if id == "" {
		id = internal.TurnIntoId(title, func(s string) bool {
			return !slices.ContainsFunc(siteConfig.Lists, func(l model.EntryList) bool {
				return s == l.Id
			})
		})
	}

	siteConfig.Lists = append(siteConfig.Lists, model.EntryList{
		Id:       id,
		Title:    title,
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

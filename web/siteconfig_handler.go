package web

import (
	"html/template"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"sort"

	"github.com/gofiber/fiber/v2"
)

type SiteConfigHandler struct {
	svc          *app.SiteConfigService
	typeRegistry *app.EntryTypeRegistry
}

type typeInMain struct {
	Type     string
	Included bool
}
type siteConfigTemplateData struct {
	Config model.SiteConfig
	Types  []typeInMain
}

func NewSiteConfigHandler(svc *app.SiteConfigService, typeRegistry *app.EntryTypeRegistry) *SiteConfigHandler {
	return &SiteConfigHandler{
		svc:          svc,
		typeRegistry: typeRegistry,
	}
}

func (h *SiteConfigHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.svc.GetSiteConfig()
	if err != nil {
		return err
	}

	types := make([]typeInMain, 0)
	for _, t := range h.typeRegistry.Types() {
		typeName, err := h.typeRegistry.TypeName(t)
		included := false
		for _, t := range siteConfig.PrimaryListInclude {
			if typeName == t {
				included = true
				break
			}
		}

		if err != nil {
			continue
		}
		types = append(types, typeInMain{
			Type:     typeName,
			Included: included,
		})
	}

	sort.Slice(types, func(i, j int) bool {
		return types[i].Type < types[j].Type
	})

	return render.RenderTemplateWithBase(c, "views/site_config", siteConfigTemplateData{
		Config: siteConfig,
		Types:  types,
	})
}

func (h *SiteConfigHandler) HandlePost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.svc.GetSiteConfig()
	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	siteConfig.Title = c.FormValue("Title")
	siteConfig.SubTitle = c.FormValue("SubTitle")
	siteConfig.PrimaryColor = c.FormValue("PrimaryColor")
	siteConfig.AuthorName = c.FormValue("AuthorName")
	siteConfig.AvatarUrl = c.FormValue("AvatarUrl")
	siteConfig.FullUrl = c.FormValue("FullUrl")
	siteConfig.HtmlHeadExtra = template.HTML(c.FormValue("HtmlHeadExtra"))
	siteConfig.FooterExtra = template.HTML(c.FormValue("FooterExtra"))
	siteConfig.PrimaryListInclude = form.Value["Include"]

	err = h.svc.UpdateSiteConfig(siteConfig)
	if err != nil {
		return err
	}

	return c.Redirect("/site-config/")
}

package web

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type SiteMapHandler struct {
	entryService      *app.EntryService
	siteConfigService *app.SiteConfigService
}

type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	Url     []SitemapUrl `xml:"url"`
}

type SitemapUrl struct {
	Loc string `xml:"loc"`
}

func NewSiteMapHandler(entryService *app.EntryService, siteConfigService *app.SiteConfigService) *SiteMapHandler {
	return &SiteMapHandler{entryService: entryService, siteConfigService: siteConfigService}
}

// Handle handles GET /sitemap.xml
func (h *SiteMapHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationXML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()
	if err != nil {
		return err
	}
	entries, err := h.entryService.FindAllByType(nil, true, false)
	if err != nil {
		return err
	}

	sitemap := Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Url:   make([]SitemapUrl, 0),
	}

	for _, entry := range entries {
		entryUrl, _ := url.JoinPath(siteConfig.FullUrl, "/posts/", url.PathEscape(entry.ID()), "/")
		sitemap.Url = append(sitemap.Url, SitemapUrl{
			Loc: entryUrl,
		})
	}

	buf := new(bytes.Buffer)
	encoder := xml.NewEncoder(buf)
	encoder.Indent("", "  ")
	err = encoder.Encode(sitemap)
	if err != nil {
		return err
	}

	return c.SendString(xml.Header + buf.String())
}

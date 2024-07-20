package web

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Atom    string     `xml:"xmlns:atom,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSDescription struct {
	XMLName xml.Name `xml:"description"`
	Text    string   `xml:",cdata"`
}

type AtomLink struct {
	Href  string `xml:"href,attr"`
	Rel   string `xml:"rel,attr,omitempty"`
	Type  string `xml:"type,attr,omitempty"`
	Title string `xml:"title,attr,omitempty"`
}

type RSSChannel struct {
	Title       string         `xml:"title"`
	Link        string         `xml:"link"`
	AtomLinks   []AtomLink     `xml:"atom:link"`
	Description RSSDescription `xml:"description"`
	PubDate     string         `xml:"pubDate"`
	LastBuild   string         `xml:"lastBuildDate"`
	Generator   string         `xml:"generator"`
	Items       []RSSItem      `xml:"item"`
}

type RSSItem struct {
	Guid        string         `xml:"guid"`
	Title       string         `xml:"title"`
	Link        string         `xml:"link"`
	PubDate     string         `xml:"pubDate"`
	Description RSSDescription `xml:"description"`
}

func RenderRSSFeed(config model.SiteConfig, entries []model.Entry) (string, error) {

	rss := RSS{
		Version: "2.0",
		Atom:    "http://www.w3.org/2005/Atom",
		Channel: RSSChannel{
			Title: config.Title,
			Link:  config.FullUrl,
			AtomLinks: []AtomLink{
				{
					Href: config.FullUrl + "/index.xml",
					Rel:  "self",
					Type: "application/rss+xml",
				},
			},
			Description: RSSDescription{
				Text: config.SubTitle,
			},
			PubDate:   time.Now().Format(time.RFC1123Z),
			LastBuild: time.Now().Format(time.RFC1123Z),
			Generator: "owl-blogs",
			Items:     make([]RSSItem, 0),
		},
	}

	for _, entry := range entries {
		content := entry.Content()
		entryUrl, _ := url.JoinPath(config.FullUrl, "/posts/", url.PathEscape(entry.ID()), "/")

		rss.Channel.Items = append(rss.Channel.Items, RSSItem{
			Guid:    entryUrl,
			Title:   entry.Title(),
			Link:    entryUrl,
			PubDate: entry.PublishedAt().Format(time.RFC1123Z),
			Description: RSSDescription{
				Text: string(content),
			},
		})
	}

	buf := new(bytes.Buffer)
	encoder := xml.NewEncoder(buf)
	encoder.Indent("", "  ")
	err := encoder.Encode(rss)
	if err != nil {
		return "", err
	}

	return xml.Header + buf.String(), nil

}

type RSSHandler struct {
	siteConfigService *app.SiteConfigService
	entrySvc          *app.EntryService
}

func NewRSSHandler(entryService *app.EntryService, siteConfigService *app.SiteConfigService) *RSSHandler {
	return &RSSHandler{entrySvc: entryService, siteConfigService: siteConfigService}
}

func (h *RSSHandler) HandleMainFeed(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationXML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()
	if err != nil {
		return err
	}

	entries, err := h.entrySvc.FindAllByType(&siteConfig.PrimaryListInclude, true, false)
	if err != nil {
		return err
	}

	// sort entries by date descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PublishedAt().After(*entries[j].PublishedAt())
	})

	rss, err := RenderRSSFeed(siteConfig, entries)
	if err != nil {
		return err
	}

	return c.SendString(rss)
}

func (h *RSSHandler) HandleListFeed(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationXML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()
	if err != nil {
		return err
	}
	listId := c.Params("list")
	list := model.EntryList{}
	for _, l := range siteConfig.Lists {
		if l.Id == listId {
			list = l
		}
	}
	if list.Id == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	entries, err := h.entrySvc.FindAllByType(&list.Include, true, false)
	if err != nil {
		return err
	}

	// sort entries by date descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PublishedAt().After(*entries[j].PublishedAt())
	})

	rss, err := RenderRSSFeed(siteConfig, entries)
	if err != nil {
		return err
	}

	return c.SendString(rss)
}

func (h *RSSHandler) HandleTagFeed(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationXML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()
	if err != nil {
		return err
	}

	tag := c.Params("tag")
	entries, err := h.entrySvc.FindAllByTag(tag, true, false)
	if err != nil {
		return err
	}

	// sort entries by date descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PublishedAt().After(*entries[j].PublishedAt())
	})

	rss, err := RenderRSSFeed(siteConfig, entries)
	if err != nil {
		return err
	}

	return c.SendString(rss)
}

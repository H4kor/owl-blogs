package web

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"owl-blogs/app"
	"owl-blogs/app/repository"
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
			Items: make([]RSSItem, 0),
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
	configRepo repository.ConfigRepository
	entrySvc   *app.EntryService
}

func NewRSSHandler(entryService *app.EntryService, configRepo repository.ConfigRepository) *RSSHandler {
	return &RSSHandler{entrySvc: entryService, configRepo: configRepo}
}

func (h *RSSHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationXML)

	siteConfig := getSiteConfig(h.configRepo)

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

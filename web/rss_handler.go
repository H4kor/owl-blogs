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
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Guid        string `xml:"guid"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func RenderRSSFeed(config model.SiteConfig, entries []model.Entry) (string, error) {

	rss := RSS{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       config.Title,
			Link:        config.FullUrl,
			Description: config.SubTitle,
			Items:       make([]RSSItem, 0),
		},
	}

	for _, entry := range entries {
		content := entry.Content()
		url, _ := url.JoinPath(config.FullUrl, "/posts/", entry.ID())
		rss.Channel.Items = append(rss.Channel.Items, RSSItem{
			Guid:        url,
			Title:       entry.Title(),
			Link:        url,
			PubDate:     entry.PublishedAt().Format(time.RFC1123Z),
			Description: string(content),
		})
	}

	buf := new(bytes.Buffer)
	err := xml.NewEncoder(buf).Encode(rss)
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

package owl

import (
	"bytes"
	"encoding/xml"
	"time"
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

func RenderRSSFeed(user User) (string, error) {

	config, _ := user.Config()

	rss := RSS{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       config.Title,
			Link:        user.FullUrl(),
			Description: config.SubTitle,
			Items:       make([]RSSItem, 0),
		},
	}

	posts, _ := user.Posts()
	for _, post := range posts {
		meta := post.Meta()
		content, _ := renderPostContent(post)
		rss.Channel.Items = append(rss.Channel.Items, RSSItem{
			Guid:        post.FullUrl(),
			Title:       post.Title(),
			Link:        post.FullUrl(),
			PubDate:     meta.Date.Format(time.RFC1123Z),
			Description: content,
		})
	}

	buf := new(bytes.Buffer)
	err := xml.NewEncoder(buf).Encode(rss)
	if err != nil {
		return "", err
	}

	return xml.Header + buf.String(), nil

}

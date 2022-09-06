package owl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type WebmentionIn struct {
	Source         string    `yaml:"source"`
	Title          string    `yaml:"title"`
	ApprovalStatus string    `yaml:"approval_status"`
	RetrievedAt    time.Time `yaml:"retrieved_at"`
}

type WebmentionOut struct {
	Target     string    `yaml:"target"`
	Supported  bool      `yaml:"supported"`
	ScannedAt  time.Time `yaml:"scanned_at"`
	LastSentAt time.Time `yaml:"last_sent_at"`
}

type HttpClient interface {
	Get(url string) ([]byte, error)
	Post(url string, data url.Values) ([]byte, error)
}

type HtmlParser interface {
	ParseHEntry(data []byte) (ParsedHEntry, error)
	ParseLinks(data []byte) ([]string, error)
	GetWebmentionEndpoint(data []byte) (string, error)
}

type OwlHttpClient struct{}

type OwlHtmlParser struct{}

type ParsedHEntry struct {
	Title string
}

func (OwlHttpClient) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return make([]byte, 0), errors.New("Failed to get url. Status code: " + fmt.Sprint(resp.StatusCode))
	}

	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (OwlHttpClient) Post(url string, data url.Values) ([]byte, error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

func (OwlHtmlParser) ParseHEntry(data []byte) (ParsedHEntry, error) {
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return ParsedHEntry{}, err
	}

	var interpretHFeed func(*html.Node, *ParsedHEntry, bool) (ParsedHEntry, error)
	interpretHFeed = func(n *html.Node, curr *ParsedHEntry, parent bool) (ParsedHEntry, error) {
		attrs := n.Attr
		for _, attr := range attrs {
			if attr.Key == "class" && strings.Contains(attr.Val, "p-name") {
				buf := &bytes.Buffer{}
				collectText(n, buf)
				curr.Title = buf.String()
				return *curr, nil
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			interpretHFeed(c, curr, false)
		}
		return *curr, nil
	}

	var findHFeed func(*html.Node) (ParsedHEntry, error)
	findHFeed = func(n *html.Node) (ParsedHEntry, error) {
		attrs := n.Attr
		for _, attr := range attrs {
			if attr.Key == "class" && strings.Contains(attr.Val, "h-entry") {
				return interpretHFeed(n, &ParsedHEntry{}, true)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			entry, err := findHFeed(c)
			if err == nil {
				return entry, nil
			}
		}
		return ParsedHEntry{}, errors.New("no h-entry found")
	}
	return findHFeed(doc)
}

func (OwlHtmlParser) ParseLinks(data []byte) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return make([]string, 0), err
	}

	var findLinks func(*html.Node) ([]string, error)
	findLinks = func(n *html.Node) ([]string, error) {
		links := make([]string, 0)
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			childLinks, _ := findLinks(c)
			links = append(links, childLinks...)
		}
		return links, nil
	}
	return findLinks(doc)

}

func (OwlHtmlParser) GetWebmentionEndpoint(data []byte) (string, error) {
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return "", err
	}

	var findEndpoint func(*html.Node) (string, error)
	findEndpoint = func(n *html.Node) (string, error) {
		if n.Type == html.ElementNode && (n.Data == "link" || n.Data == "a") {
			for _, attr := range n.Attr {
				if attr.Key == "rel" && attr.Val == "webmention" {
					for _, attr := range n.Attr {
						if attr.Key == "href" {
							return attr.Val, nil
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			endpoint, err := findEndpoint(c)
			if err == nil {
				return endpoint, nil
			}
		}
		return "", errors.New("no webmention endpoint found")
	}
	return findEndpoint(doc)
}

package owl

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Webmention struct {
	Source   string `yaml:"source"`
	Title    string `yaml:"title"`
	Approved bool   `yaml:"approved"`
}

type HttpRetriever interface {
	Get(url string) ([]byte, error)
}

type MicroformatParser interface {
	ParseHEntry(data []byte) (ParsedHEntry, error)
}

type OwlHttpRetriever struct{}

type OwlMicroformatParser struct{}

type ParsedHEntry struct {
	Title string
}

func (OwlHttpRetriever) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	var data []byte
	_, err = resp.Body.Read(data)
	// TODO: encoding
	return data, err
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

func (OwlMicroformatParser) ParseHEntry(data []byte) (ParsedHEntry, error) {
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

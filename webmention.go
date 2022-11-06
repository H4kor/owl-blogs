package owl

import (
	"bytes"
	"errors"
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

func (webmention *WebmentionIn) UpdateWith(update WebmentionIn) {
	if update.Title != "" {
		webmention.Title = update.Title
	}
	if update.ApprovalStatus != "" {
		webmention.ApprovalStatus = update.ApprovalStatus
	}
	if !update.RetrievedAt.IsZero() {
		webmention.RetrievedAt = update.RetrievedAt
	}
}

type WebmentionOut struct {
	Target     string    `yaml:"target"`
	Supported  bool      `yaml:"supported"`
	ScannedAt  time.Time `yaml:"scanned_at"`
	LastSentAt time.Time `yaml:"last_sent_at"`
}

func (webmention *WebmentionOut) UpdateWith(update WebmentionOut) {
	if update.Supported {
		webmention.Supported = update.Supported
	}
	if !update.ScannedAt.IsZero() {
		webmention.ScannedAt = update.ScannedAt
	}
	if !update.LastSentAt.IsZero() {
		webmention.LastSentAt = update.LastSentAt
	}
}

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type HtmlParser interface {
	ParseHEntry(resp *http.Response) (ParsedHEntry, error)
	ParseLinks(resp *http.Response) ([]string, error)
	ParseLinksFromString(string) ([]string, error)
	GetWebmentionEndpoint(resp *http.Response) (string, error)
	GetRedirctUris(resp *http.Response) ([]string, error)
}

type OwlHttpClient = http.Client

type OwlHtmlParser struct{}

type ParsedHEntry struct {
	Title string
}

func collectText(n *html.Node, buf *bytes.Buffer) {

	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

func readResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func (OwlHtmlParser) ParseHEntry(resp *http.Response) (ParsedHEntry, error) {
	htmlStr, err := readResponseBody(resp)
	if err != nil {
		return ParsedHEntry{}, err
	}
	doc, err := html.Parse(strings.NewReader(htmlStr))
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

func (OwlHtmlParser) ParseLinks(resp *http.Response) ([]string, error) {
	htmlStr, err := readResponseBody(resp)
	if err != nil {
		return []string{}, err
	}
	return OwlHtmlParser{}.ParseLinksFromString(htmlStr)
}

func (OwlHtmlParser) ParseLinksFromString(htmlStr string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
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

func (OwlHtmlParser) GetWebmentionEndpoint(resp *http.Response) (string, error) {
	//request url
	requestUrl := resp.Request.URL

	// Check link headers
	for _, linkHeader := range resp.Header["Link"] {
		linkHeaderParts := strings.Split(linkHeader, ",")
		for _, linkHeaderPart := range linkHeaderParts {
			linkHeaderPart = strings.TrimSpace(linkHeaderPart)
			params := strings.Split(linkHeaderPart, ";")
			if len(params) != 2 {
				continue
			}
			for _, param := range params[1:] {
				param = strings.TrimSpace(param)
				if strings.Contains(param, "webmention") {
					link := strings.Split(params[0], ";")[0]
					link = strings.Trim(link, "<>")
					linkUrl, err := url.Parse(link)
					if err != nil {
						return "", err
					}
					return requestUrl.ResolveReference(linkUrl).String(), nil
				}
			}
		}
	}

	htmlStr, err := readResponseBody(resp)
	if err != nil {
		return "", err
	}
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}

	var findEndpoint func(*html.Node) (string, error)
	findEndpoint = func(n *html.Node) (string, error) {
		if n.Type == html.ElementNode && (n.Data == "link" || n.Data == "a") {
			for _, attr := range n.Attr {
				if attr.Key == "rel" {
					vals := strings.Split(attr.Val, " ")
					for _, val := range vals {
						if val == "webmention" {
							for _, attr := range n.Attr {
								if attr.Key == "href" {
									return attr.Val, nil
								}
							}
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
	linkUrlStr, err := findEndpoint(doc)
	if err != nil {
		return "", err
	}
	linkUrl, err := url.Parse(linkUrlStr)
	if err != nil {
		return "", err
	}
	return requestUrl.ResolveReference(linkUrl).String(), nil
}

func (OwlHtmlParser) GetRedirctUris(resp *http.Response) ([]string, error) {
	//request url
	requestUrl := resp.Request.URL

	htmlStr, err := readResponseBody(resp)
	if err != nil {
		return make([]string, 0), err
	}
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return make([]string, 0), err
	}

	var findLinks func(*html.Node) ([]string, error)
	findLinks = func(n *html.Node) ([]string, error) {
		links := make([]string, 0)
		if n.Type == html.ElementNode && n.Data == "link" {
			// check for rel="redirect_uri"
			rel := ""
			href := ""

			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
				}
				if attr.Key == "rel" {
					rel = attr.Val
				}
			}
			if rel == "redirect_uri" {
				linkUrl, err := url.Parse(href)
				if err == nil {
					links = append(links, requestUrl.ResolveReference(linkUrl).String())
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

package app

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"owl-blogs/app/owlhttp"
	"owl-blogs/app/repository"
	"owl-blogs/interactions"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type WebmentionService struct {
	InteractionRepository repository.InteractionRepository
	EntryRepository       repository.EntryRepository
	Http                  owlhttp.HttpClient
}

type ParsedHEntry struct {
	Title string
}

func NewWebmentionService(
	interactionRepository repository.InteractionRepository,
	entryRepository repository.EntryRepository,
	http owlhttp.HttpClient,
) *WebmentionService {
	return &WebmentionService{
		InteractionRepository: interactionRepository,
		EntryRepository:       entryRepository,
		Http:                  http,
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

func collectText(n *html.Node, buf *bytes.Buffer) {

	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

func (WebmentionService) ParseHEntry(resp *http.Response) (ParsedHEntry, error) {
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

func (s *WebmentionService) GetExistingWebmention(entryId string, source string, target string) (*interactions.Webmention, error) {
	inters, err := s.InteractionRepository.FindAll(entryId)
	if err != nil {
		return nil, err
	}
	for _, interaction := range inters {
		if webm, ok := interaction.(*interactions.Webmention); ok {
			m := webm.MetaData().(*interactions.WebmentionMetaData)
			if m.Source == source && m.Target == target {
				return webm, nil
			}
		}
	}
	return nil, nil
}

func (s *WebmentionService) ProcessWebmention(source string, target string) error {
	resp, err := s.Http.Get(source)
	if err != nil {
		return err
	}

	hEntry, err := s.ParseHEntry(resp)
	if err != nil {
		return err
	}

	entryId := UrlToEntryId(target)
	_, err = s.EntryRepository.FindById(entryId)
	if err != nil {
		return err
	}

	webmention, err := s.GetExistingWebmention(entryId, source, target)
	if err != nil {
		return err
	}
	if webmention != nil {
		data := interactions.WebmentionMetaData{
			Source: source,
			Target: target,
			Title:  hEntry.Title,
		}
		webmention.SetMetaData(&data)
		webmention.SetEntryID(entryId)
		webmention.SetCreatedAt(time.Now())
		err = s.InteractionRepository.Update(webmention)
		return err
	} else {
		webmention = &interactions.Webmention{}
		data := interactions.WebmentionMetaData{
			Source: source,
			Target: target,
			Title:  hEntry.Title,
		}
		webmention.SetMetaData(&data)
		webmention.SetEntryID(entryId)
		webmention.SetCreatedAt(time.Now())
		err = s.InteractionRepository.Create(webmention)
		return err
	}
}

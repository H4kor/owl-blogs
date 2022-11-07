package mocks

import (
	"h4kor/owl-blogs"
	"io"
	"net/http"
	"net/url"
)

type MockHtmlParser struct{}

func (*MockHtmlParser) ParseHEntry(resp *http.Response) (owl.ParsedHEntry, error) {
	return owl.ParsedHEntry{Title: "Mock Title"}, nil

}
func (*MockHtmlParser) ParseLinks(resp *http.Response) ([]string, error) {
	return []string{"http://example.com"}, nil

}
func (*MockHtmlParser) ParseLinksFromString(string) ([]string, error) {
	return []string{"http://example.com"}, nil

}
func (*MockHtmlParser) GetWebmentionEndpoint(resp *http.Response) (string, error) {
	return "http://example.com/webmention", nil

}
func (*MockHtmlParser) GetRedirctUris(resp *http.Response) ([]string, error) {
	return []string{"http://example.com/redirect"}, nil
}

type MockParseLinksHtmlParser struct {
	Links []string
}

func (*MockParseLinksHtmlParser) ParseHEntry(resp *http.Response) (owl.ParsedHEntry, error) {
	return owl.ParsedHEntry{Title: "Mock Title"}, nil
}
func (parser *MockParseLinksHtmlParser) ParseLinks(resp *http.Response) ([]string, error) {
	return parser.Links, nil
}
func (parser *MockParseLinksHtmlParser) ParseLinksFromString(string) ([]string, error) {
	return parser.Links, nil
}
func (*MockParseLinksHtmlParser) GetWebmentionEndpoint(resp *http.Response) (string, error) {
	return "http://example.com/webmention", nil
}
func (parser *MockParseLinksHtmlParser) GetRedirctUris(resp *http.Response) ([]string, error) {
	return parser.Links, nil
}

type MockHttpClient struct{}

func (*MockHttpClient) Get(url string) (resp *http.Response, err error) {
	return &http.Response{}, nil
}
func (*MockHttpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {

	return &http.Response{}, nil
}
func (*MockHttpClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {

	return &http.Response{}, nil
}

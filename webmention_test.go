package owl_test

import (
	"bytes"
	"h4kor/owl-blogs"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func constructResponse(html []byte) *http.Response {
	url, _ := url.Parse("http://example.com/foo/bar")
	return &http.Response{
		Request: &http.Request{
			URL: url,
		},
		Body: io.NopCloser(bytes.NewReader([]byte(html))),
	}
}

//
// https://www.w3.org/TR/webmention/#h-webmention-verification
//

func TestParseValidHEntry(t *testing.T) {
	html := []byte("<div class=\"h-entry\"><div class=\"p-name\">Foo</div></div>")
	parser := &owl.OwlHtmlParser{}
	entry, err := parser.ParseHEntry(&http.Response{Body: io.NopCloser(bytes.NewReader(html))})

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if entry.Title != "Foo" {
		t.Errorf("Wrong Title. Expected %v, got %v", "Foo", entry.Title)
	}
}

func TestParseValidHEntryWithoutTitle(t *testing.T) {
	html := []byte("<div class=\"h-entry\"></div><div class=\"p-name\">Foo</div>")
	parser := &owl.OwlHtmlParser{}
	entry, err := parser.ParseHEntry(&http.Response{Body: io.NopCloser(bytes.NewReader(html))})

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if entry.Title != "" {
		t.Errorf("Wrong Title. Expected %v, got %v", "Foo", entry.Title)
	}
}

func TestGetWebmentionEndpointLink(t *testing.T) {
	html := []byte("<link rel=\"webmention\" href=\"http://example.com/webmention\" />")
	parser := &owl.OwlHtmlParser{}
	endpoint, err := parser.GetWebmentionEndpoint(constructResponse(html))

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if endpoint != "http://example.com/webmention" {
		t.Errorf("Wrong endpoint. Expected %v, got %v", "http://example.com/webmention", endpoint)
	}
}

func TestGetWebmentionEndpointLinkA(t *testing.T) {
	html := []byte("<a rel=\"webmention\" href=\"http://example.com/webmention\" />")
	parser := &owl.OwlHtmlParser{}
	endpoint, err := parser.GetWebmentionEndpoint(constructResponse(html))

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if endpoint != "http://example.com/webmention" {
		t.Errorf("Wrong endpoint. Expected %v, got %v", "http://example.com/webmention", endpoint)
	}
}

func TestGetWebmentionEndpointLinkAFakeWebmention(t *testing.T) {
	html := []byte("<a rel=\"not-webmention\" href=\"http://example.com/foo\" /><a rel=\"webmention\" href=\"http://example.com/webmention\" />")
	parser := &owl.OwlHtmlParser{}
	endpoint, err := parser.GetWebmentionEndpoint(constructResponse(html))

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if endpoint != "http://example.com/webmention" {
		t.Errorf("Wrong endpoint. Expected %v, got %v", "http://example.com/webmention", endpoint)
	}
}

func TestGetWebmentionEndpointLinkHeader(t *testing.T) {
	html := []byte("")
	parser := &owl.OwlHtmlParser{}
	resp := constructResponse(html)
	resp.Header = http.Header{"Link": []string{"<http://example.com/webmention>; rel=\"webmention\""}}
	endpoint, err := parser.GetWebmentionEndpoint(resp)

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if endpoint != "http://example.com/webmention" {
		t.Errorf("Wrong endpoint. Expected %v, got %v", "http://example.com/webmention", endpoint)
	}
}

func TestGetWebmentionEndpointLinkHeaderCommas(t *testing.T) {
	html := []byte("")
	parser := &owl.OwlHtmlParser{}
	resp := constructResponse(html)
	resp.Header = http.Header{
		"Link": []string{"<https://webmention.rocks/test/19/webmention/error>; rel=\"other\", <https://webmention.rocks/test/19/webmention>; rel=\"webmention\""},
	}
	endpoint, err := parser.GetWebmentionEndpoint(resp)

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if endpoint != "https://webmention.rocks/test/19/webmention" {
		t.Errorf("Wrong endpoint. Expected %v, got %v", "https://webmention.rocks/test/19/webmention", endpoint)
	}
}

func TestGetWebmentionEndpointRelativeLink(t *testing.T) {
	html := []byte("<link rel=\"webmention\" href=\"/webmention\" />")
	parser := &owl.OwlHtmlParser{}
	endpoint, err := parser.GetWebmentionEndpoint(constructResponse(html))

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if endpoint != "http://example.com/webmention" {
		t.Errorf("Wrong endpoint. Expected %v, got %v", "http://example.com/webmention", endpoint)
	}
}

func TestGetWebmentionEndpointRelativeLinkInHeader(t *testing.T) {
	html := []byte("<link rel=\"webmention\" href=\"/webmention\" />")
	parser := &owl.OwlHtmlParser{}
	resp := constructResponse(html)
	resp.Header = http.Header{"Link": []string{"</webmention>; rel=\"webmention\""}}
	endpoint, err := parser.GetWebmentionEndpoint(resp)

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if endpoint != "http://example.com/webmention" {
		t.Errorf("Wrong endpoint. Expected %v, got %v", "http://example.com/webmention", endpoint)
	}
}

func TestRealWorldWebmention(t *testing.T) {
	links := []string{
		"https://webmention.rocks/test/1",
		"https://webmention.rocks/test/2",
		"https://webmention.rocks/test/3",
		"https://webmention.rocks/test/4",
		"https://webmention.rocks/test/5",
		"https://webmention.rocks/test/6",
		"https://webmention.rocks/test/7",
		"https://webmention.rocks/test/8",
		"https://webmention.rocks/test/9",
		// "https://webmention.rocks/test/10", // not supported
		"https://webmention.rocks/test/11",
		"https://webmention.rocks/test/12",
		"https://webmention.rocks/test/13",
		"https://webmention.rocks/test/14",
		"https://webmention.rocks/test/15",
		"https://webmention.rocks/test/16",
		"https://webmention.rocks/test/17",
		"https://webmention.rocks/test/18",
		"https://webmention.rocks/test/19",
		"https://webmention.rocks/test/20",
		"https://webmention.rocks/test/21",
		"https://webmention.rocks/test/22",
		"https://webmention.rocks/test/23/page",
	}

	for _, link := range links {
		parser := &owl.OwlHtmlParser{}
		client := &owl.OwlHttpClient{}
		html, _ := client.Get(link)
		_, err := parser.GetWebmentionEndpoint(html)

		if err != nil {
			t.Errorf("Unable to find webmention: %v for link %v", err, link)
		}
	}

}

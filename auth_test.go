package owl_test

import (
	"h4kor/owl-blogs"
	"h4kor/owl-blogs/test/assertions"
	"net/http"
	"testing"
)

func TestGetRedirctUrisLink(t *testing.T) {
	html := []byte("<link rel=\"redirect_uri\" href=\"http://example.com/redirect\" />")
	parser := &owl.OwlHtmlParser{}
	uris, err := parser.GetRedirctUris(constructResponse(html))

	assertions.AssertNoError(t, err, "Unable to parse feed")

	assertions.AssertArrayContains(t, uris, "http://example.com/redirect")
}

func TestGetRedirctUrisLinkMultiple(t *testing.T) {
	html := []byte(`
		<link rel="redirect_uri" href="http://example.com/redirect1" />
		<link rel="redirect_uri" href="http://example.com/redirect2" />
		<link rel="redirect_uri" href="http://example.com/redirect3" />
		<link rel="foo" href="http://example.com/redirect4" />
		<link href="http://example.com/redirect5" />	
	`)
	parser := &owl.OwlHtmlParser{}
	uris, err := parser.GetRedirctUris(constructResponse(html))

	assertions.AssertNoError(t, err, "Unable to parse feed")

	assertions.AssertArrayContains(t, uris, "http://example.com/redirect1")
	assertions.AssertArrayContains(t, uris, "http://example.com/redirect2")
	assertions.AssertArrayContains(t, uris, "http://example.com/redirect3")
	assertions.AssertLen(t, uris, 3)
}

func TestGetRedirectUrisLinkHeader(t *testing.T) {
	html := []byte("")
	parser := &owl.OwlHtmlParser{}
	resp := constructResponse(html)
	resp.Header = http.Header{"Link": []string{"<http://example.com/redirect>; rel=\"redirect_uri\""}}
	uris, err := parser.GetRedirctUris(resp)

	assertions.AssertNoError(t, err, "Unable to parse feed")
	assertions.AssertArrayContains(t, uris, "http://example.com/redirect")
}

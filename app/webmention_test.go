package app_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"owl-blogs/app"
	"owl-blogs/infra"
	"owl-blogs/interactions"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
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

type MockHttpClient struct {
	PageContent string
}

// Post implements owlhttp.HttpClient.
func (MockHttpClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	panic("unimplemented")
}

// PostForm implements owlhttp.HttpClient.
func (MockHttpClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	panic("unimplemented")
}

func (c *MockHttpClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte(c.PageContent))),
	}, nil
}

func getWebmentionService() *app.WebmentionService {
	db := test.NewMockDb()
	entryRegister := app.NewEntryTypeRegistry()
	entryRegister.Register(&test.MockEntry{})
	entryRepo := infra.NewEntryRepository(db, entryRegister)

	interactionRegister := app.NewInteractionTypeRegistry()
	interactionRegister.Register(&interactions.Webmention{})

	interactionRepo := infra.NewInteractionRepo(db, interactionRegister)

	configRepo := infra.NewConfigRepo(db)

	bus := app.NewEventBus()

	http := infra.OwlHttpClient{}
	return app.NewWebmentionService(
		configRepo, interactionRepo, entryRepo, &http, bus,
	)
}

//
// https://www.w3.org/TR/webmention/#h-webmention-verification
//

func TestParseValidHEntry(t *testing.T) {
	html := []byte("<div class=\"h-entry\"><div class=\"p-name\">Foo</div></div>")
	entry, err := app.ParseHEntry(&http.Response{Body: io.NopCloser(bytes.NewReader(html))})

	require.NoError(t, err)
	require.Equal(t, entry.Title, "Foo")
}

func TestParseValidHEntryWithoutTitle(t *testing.T) {
	html := []byte("<div class=\"h-entry\"></div><div class=\"p-name\">Foo</div>")
	entry, err := app.ParseHEntry(&http.Response{Body: io.NopCloser(bytes.NewReader(html))})

	require.NoError(t, err)
	require.Equal(t, entry.Title, "")
}

func TestCreateNewWebmention(t *testing.T) {
	service := getWebmentionService()
	service.Http = &MockHttpClient{
		PageContent: "<div class=\"h-entry\"><div class=\"p-name\">Foo</div></div>",
	}
	entry := test.MockEntry{}
	service.EntryRepository.Create(&entry)

	err := service.ProcessWebmention(
		"http://example.com/foo",
		fmt.Sprintf("https.//example.com/posts/%s/", entry.ID()),
	)
	require.NoError(t, err)

	inters, err := service.InteractionRepository.FindAll(entry.ID())
	require.NoError(t, err)
	require.Equal(t, len(inters), 1)
	webm := inters[0].(*interactions.Webmention)
	meta := webm.MetaData().(*interactions.WebmentionMetaData)
	require.Equal(t, meta.Source, "http://example.com/foo")
	require.Equal(t, meta.Target, fmt.Sprintf("https.//example.com/posts/%s/", entry.ID()))
	require.Equal(t, meta.Title, "Foo")
}

func TestGetWebmentionEndpointLink(t *testing.T) {
	html := []byte("<link rel=\"webmention\" href=\"http://example.com/webmention\" />")
	endpoint, err := app.GetWebmentionEndpoint(constructResponse(html))

	require.NoError(t, err)

	require.Equal(t, endpoint, "http://example.com/webmention")
}

func TestGetWebmentionEndpointLinkA(t *testing.T) {
	html := []byte("<a rel=\"webmention\" href=\"http://example.com/webmention\" />")
	endpoint, err := app.GetWebmentionEndpoint(constructResponse(html))

	require.NoError(t, err)
	require.Equal(t, endpoint, "http://example.com/webmention")
}

func TestGetWebmentionEndpointLinkAFakeWebmention(t *testing.T) {
	html := []byte("<a rel=\"not-webmention\" href=\"http://example.com/foo\" /><a rel=\"webmention\" href=\"http://example.com/webmention\" />")
	endpoint, err := app.GetWebmentionEndpoint(constructResponse(html))

	require.NoError(t, err)
	require.Equal(t, endpoint, "http://example.com/webmention")
}

func TestGetWebmentionEndpointLinkHeader(t *testing.T) {
	html := []byte("")
	resp := constructResponse(html)
	resp.Header = http.Header{"Link": []string{"<http://example.com/webmention>; rel=\"webmention\""}}
	endpoint, err := app.GetWebmentionEndpoint(resp)

	require.NoError(t, err)
	require.Equal(t, endpoint, "http://example.com/webmention")
}

func TestGetWebmentionEndpointLinkHeaderCommas(t *testing.T) {
	html := []byte("")
	resp := constructResponse(html)
	resp.Header = http.Header{
		"Link": []string{"<https://webmention.rocks/test/19/webmention/error>; rel=\"other\", <https://webmention.rocks/test/19/webmention>; rel=\"webmention\""},
	}
	endpoint, err := app.GetWebmentionEndpoint(resp)

	require.NoError(t, err)
	require.Equal(t, endpoint, "https://webmention.rocks/test/19/webmention")
}

func TestGetWebmentionEndpointRelativeLink(t *testing.T) {
	html := []byte("<link rel=\"webmention\" href=\"/webmention\" />")
	endpoint, err := app.GetWebmentionEndpoint(constructResponse(html))

	require.NoError(t, err)
	require.Equal(t, endpoint, "http://example.com/webmention")
}

func TestGetWebmentionEndpointRelativeLinkInHeader(t *testing.T) {
	html := []byte("<link rel=\"webmention\" href=\"/webmention\" />")
	resp := constructResponse(html)
	resp.Header = http.Header{"Link": []string{"</webmention>; rel=\"webmention\""}}
	endpoint, err := app.GetWebmentionEndpoint(resp)

	require.NoError(t, err)
	require.Equal(t, endpoint, "http://example.com/webmention")
}

// func TestRealWorldWebmention(t *testing.T) {
//  service := getWebmentionService()
// 	links := []string{
// 		"https://webmention.rocks/test/1",
// 		"https://webmention.rocks/test/2",
// 		"https://webmention.rocks/test/3",
// 		"https://webmention.rocks/test/4",
// 		"https://webmention.rocks/test/5",
// 		"https://webmention.rocks/test/6",
// 		"https://webmention.rocks/test/7",
// 		"https://webmention.rocks/test/8",
// 		"https://webmention.rocks/test/9",
// 		// "https://webmention.rocks/test/10", // not supported
// 		"https://webmention.rocks/test/11",
// 		"https://webmention.rocks/test/12",
// 		"https://webmention.rocks/test/13",
// 		"https://webmention.rocks/test/14",
// 		"https://webmention.rocks/test/15",
// 		"https://webmention.rocks/test/16",
// 		"https://webmention.rocks/test/17",
// 		"https://webmention.rocks/test/18",
// 		"https://webmention.rocks/test/19",
// 		"https://webmention.rocks/test/20",
// 		"https://webmention.rocks/test/21",
// 		"https://webmention.rocks/test/22",
// 		"https://webmention.rocks/test/23/page",
// 	}

// 	for _, link := range links {
//
// 		client := &owl.OwlHttpClient{}
// 		html, _ := client.Get(link)
// 		_, err := app.GetWebmentionEndpoint(html)

// 		if err != nil {
// 			t.Errorf("Unable to find webmention: %v for link %v", err, link)
// 		}
// 	}

// }

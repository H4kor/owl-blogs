package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func TestRssFeeds(t *testing.T) {
	type give struct {
		path string
	}
	tests := []struct {
		give
	}{
		{
			give: give{path: "/index.xml"},
		},
		{
			give: give{path: "/lists/list_one/index.xml"},
		},
		{
			give: give{path: "/tags/a_tag/index.xml"},
		},
	}
	// test
	app := DefaultTestApp()

	srv := adaptor.FiberApp(app.FiberApp)
	for _, test := range tests {
		// test
		req := httptest.NewRequest("GET", test.give.path, nil)
		resp := httptest.NewRecorder()
		srv.ServeHTTP(resp, req)
		// validation
		require.Equal(t, 200, resp.Result().StatusCode)
		require.Contains(t, resp.Body.String(), "<rss version=\"2.0\"")
	}
}

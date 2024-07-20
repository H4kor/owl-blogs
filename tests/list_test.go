package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func Test404PageForList(t *testing.T) {
	type give struct {
		path string
	}
	type want struct {
		status   int
		contains string
	}
	tests := []struct {
		give
		want
	}{
		{
			give: give{path: "/lists/not-found/"},
			want: want{status: 404, contains: "<h1>List not found</h1>"},
		},
		{
			give: give{path: "/lists/list_one/"},
			want: want{status: 200, contains: ""},
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
		require.Equal(t, test.want.status, resp.Result().StatusCode)
		require.Contains(t, resp.Body.String(), test.want.contains)
	}
}

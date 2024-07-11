package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func TestInvalidPaginationHandling(t *testing.T) {
	type give struct {
		path string
	}
	type want struct {
		status int
		path   string
	}
	tests := []struct {
		give
		want
	}{
		{
			give: give{path: "/"},
			want: want{path: "", status: 200},
		},
		{
			give: give{path: "/?page=1"},
			want: want{path: "", status: 200},
		},
		{
			give: give{path: "/?page=2000"},
			want: want{path: "", status: 200},
		},
		{
			give: give{path: "/?page=aaa"},
			want: want{path: "/", status: 301},
		},
		{
			give: give{path: "/?page=2%25%27%20ORDER%20BY%201%23"},
			want: want{path: "/", status: 301},
		},
		{
			give: give{path: "/lists/list_one/"},
			want: want{path: "", status: 200},
		},
		{
			give: give{path: "/lists/list_one/?page=1"},
			want: want{path: "", status: 200},
		},
		{
			give: give{path: "/lists/list_one/?page=2000"},
			want: want{path: "", status: 200},
		},
		{
			give: give{path: "/lists/list_one/?page=aaa"},
			want: want{path: "/lists/list_one/", status: 301},
		},
		{
			give: give{path: "/lists/list_one/?page=2%25%27%20ORDER%20BY%201%23"},
			want: want{path: "/lists/list_one/", status: 301},
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
		require.Equal(t, test.want.path, resp.Header().Get("Location"))
	}
}

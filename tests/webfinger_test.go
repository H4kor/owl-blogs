package tests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func TestWebfinger(t *testing.T) {
	type give struct {
		path string
	}
	type want struct {
		status int
		equals map[string]interface{}
	}
	tests := []struct {
		give
		want
	}{
		{
			give: give{"/.well-known/webfinger?resource=acct:tester@example.com"},
			want: want{200, map[string]interface{}{
				"subject": "acct:tester@example.com",
			}},
		},
		{
			give: give{"/.well-known/webfinger"},
			want: want{404, map[string]interface{}{}},
		},
		{
			give: give{"/.well-known/webfinger?resource=acct:no-game@example.com"},
			want: want{404, map[string]interface{}{}},
		},
		{
			give: give{"/.well-known/webfinger?resource=acct:rps@example.com"},
			want: want{404, map[string]interface{}{}},
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
		t.Log(resp.Body.String())
		require.Equal(t, test.want.status, resp.Result().StatusCode)
		for key, value := range test.want.equals {
			var data map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			require.NoError(t, err)
			require.Equal(t, value, data[key])
		}
	}

}

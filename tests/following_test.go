package tests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func TestFollowing(t *testing.T) {
	//setup
	app := DefaultTestApp()
	srv := adaptor.FiberApp(app.FiberApp)
	actorUrl := GetActorUrl(srv)
	inbox := GetInboxUrl(srv)
	mock := NewMockAPServer()

	// test
	{
		follow := map[string]interface{}{
			"@context": "https://www.w3.org/ns/activitystreams",
			"id":       mock.MockActivityUrl("1"),
			"type":     "Follow",
			"actor":    mock.MockActorUrl("foo"),
			"object":   actorUrl,
		}
		reqData, _ := json.Marshal(follow)
		req, err := mock.SignedRequest(actorUrl, "POST", Path(inbox), reqData)
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		srv.ServeHTTP(resp, req)
		require.Equal(t, resp.Result().StatusCode, 200)
		time.Sleep(200 * time.Millisecond)
	}
	// verification
	{
		followers := GetFollowersUrl(srv)
		req := httptest.NewRequest("GET", Path(followers), nil)
		resp := httptest.NewRecorder()
		srv.ServeHTTP(resp, req)
		var data map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &data)
		require.NoError(t, err)
		require.Equal(t, []interface{}{mock.MockActorUrl("foo")}, data["items"])
	}

}

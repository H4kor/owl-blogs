package tests

import (
	"crypto/rsa"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	owlblogs "owl-blogs"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"owl-blogs/web"
	"strconv"
	"testing"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Sign(privateKey *rsa.PrivateKey, pubKeyId string, body []byte, r *http.Request) error {
	prefs := []httpsig.Algorithm{httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256
	// The "Date" and "Digest" headers must already be set on r, as well as r.URL.
	headersToSign := []string{httpsig.RequestTarget, "host", "date"}
	if body != nil {
		headersToSign = append(headersToSign, "digest")
	}
	signer, _, err := httpsig.NewSigner(prefs, digestAlgorithm, headersToSign, httpsig.Signature, 0)
	if err != nil {
		return err
	}
	// To sign the digest, we need to give the signer a copy of the body...
	// ...but it is optional, no digest will be signed if given "nil"
	// If r were a http.ResponseWriter, call SignResponse instead.
	err = signer.SignRequest(privateKey, pubKeyId, r, body)

	slog.Info("Signed Request", "req", r.Header)
	return err
}

// DefaultTestApp creates an app with a default configuration for testing
//
// == SiteConfig ==
// FullUrl: https://example.com
// Lists: two lists named "list_one" and "list_two"
//
// == ActivityPub Config ==
// PreferredUsername: tester
func DefaultTestApp() *web.WebApp {
	// setup default test server
	dbId, _ := uuid.NewRandom()
	db := infra.NewSqliteDB("file:" + dbId.String() + "?mode=memory&cache=shared")
	app := owlblogs.App(db)
	cfg, _ := app.SiteConfigService.GetSiteConfig()
	cfg.FullUrl = "https://example.com"
	cfg.Lists = []model.EntryList{
		{
			Id:      "list_one",
			Title:   "List One",
			Include: []string{"Article"},
		},
		{
			Id:      "list_two",
			Title:   "List Two",
			Include: []string{"Note"},
		},
	}
	app.SiteConfigService.UpdateSiteConfig(cfg)
	acPubCfg, _ := app.ActivityPubService.GetApConfig()
	acPubCfg.PreferredUsername = "tester"
	app.SiteConfigRepo.Update(config.ACT_PUB_CONF_NAME, acPubCfg)

	return app

}

func GetActorUrl(srv http.HandlerFunc) string {
	req := httptest.NewRequest(
		"GET", "/.well-known/webfinger?resource=acct:tester@example.com", nil)
	resp := httptest.NewRecorder()
	srv.ServeHTTP(resp, req)
	var data map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &data)
	return data["links"].([]interface{})[0].(map[string]interface{})["href"].(string)
}

func EnsureFollowed(t *testing.T, srv http.HandlerFunc, mock MockApServer, follower string) {
	actorUrl := GetActorUrl(srv)
	inbox := GetInboxUrl(srv)
	follow := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       mock.MockActivityUrl(strconv.Itoa(time.Now().Nanosecond())),
		"type":     "Follow",
		"actor":    follower,
		"object":   actorUrl,
	}
	reqData, _ := json.Marshal(follow)
	req, err := mock.SignedRequest(actorUrl, "POST", Path(inbox), reqData)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	srv.ServeHTTP(resp, req)
	require.Equal(t, resp.Result().StatusCode, 200)

}

func GetActor(srv http.HandlerFunc) map[string]interface{} {
	url := GetActorUrl(srv)
	req := httptest.NewRequest(
		"GET", url, nil)
	req.Header.Set("Accept", "application/ld+json")
	resp := httptest.NewRecorder()
	srv.ServeHTTP(resp, req)
	var data map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &data)
	return data
}

func GetInboxUrl(srv http.HandlerFunc) string {
	actor := GetActor(srv)
	return actor["inbox"].(string)
}

func GetFollowersUrl(srv http.HandlerFunc) string {
	actor := GetActor(srv)
	return actor["followers"].(string)
}

func GetOutboxUrl(srv http.HandlerFunc) string {
	actor := GetActor(srv)
	return actor["outbox"].(string)
}

func Path(u string) string {
	_url, _ := url.Parse(u)
	return _url.Path
}

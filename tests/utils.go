package tests

import (
	"crypto/rsa"
	"log/slog"
	"net/http"
	owlblogs "owl-blogs"
	"owl-blogs/config"
	"owl-blogs/infra"
	"owl-blogs/web"

	"github.com/go-fed/httpsig"
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
//
// == ActivityPub Config ==
// PreferredUsername: tester
func DefaultTestApp() *web.WebApp {
	// setup default test server
	db := infra.NewSqliteDB(":memory:")
	app := owlblogs.App(db)
	cfg, _ := app.SiteConfigService.GetSiteConfig()
	cfg.FullUrl = "https://example.com"
	app.SiteConfigService.UpdateSiteConfig(cfg)
	acPubCfg, _ := app.ActivityPubService.GetApConfig()
	acPubCfg.PreferredUsername = "tester"
	app.SiteConfigRepo.Update(config.ACT_PUB_CONF_NAME, acPubCfg)

	return app

}

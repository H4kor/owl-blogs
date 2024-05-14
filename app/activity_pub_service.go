package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"reflect"
	"time"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-fed/httpsig"
)

type ActivityPubConfig struct {
	PreferredUsername string
	PublicKeyPem      string
	PrivateKeyPem     string
}

// Form implements app.AppConfig.
func (cfg *ActivityPubConfig) Form(binSvc model.BinaryStorageInterface) string {
	f, _ := render.RenderTemplateToString("forms/ActivityPubConfig", cfg)
	return f
}

// ParseFormData implements app.AppConfig.
func (cfg *ActivityPubConfig) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	cfg.PreferredUsername = data.FormValue("PreferredUsername")
	cfg.PublicKeyPem = data.FormValue("PublicKeyPem")
	cfg.PrivateKeyPem = data.FormValue("PrivateKeyPem")
	return nil
}

func (cfg *ActivityPubConfig) PrivateKey() *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(cfg.PrivateKeyPem))
	privKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privKey
}

type ActivityPubService struct {
	followersRepo repository.FollowerRepository
	configRepo    repository.ConfigRepository
}

func NewActivityPubService(followersRepo repository.FollowerRepository, configRepo repository.ConfigRepository) *ActivityPubService {
	return &ActivityPubService{
		followersRepo: followersRepo,
		configRepo:    configRepo,
	}
}

func (svc *ActivityPubService) defaultConfig() ActivityPubConfig {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := privKey.Public().(*rsa.PublicKey)

	pubKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubKey),
		},
	)

	privKeyPrm := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	)

	return ActivityPubConfig{
		PreferredUsername: "blog",
		PublicKeyPem:      string(pubKeyPem),
		PrivateKeyPem:     string(privKeyPrm),
	}
}

func (svc *ActivityPubService) GetApConfig() (ActivityPubConfig, error) {
	apConfig := ActivityPubConfig{}
	err := svc.configRepo.Get(config.ACT_PUB_CONF_NAME, &apConfig)
	if err != nil {
		println("ERROR IN ACTIVITY PUB CONFIG")
		return ActivityPubConfig{}, err
	}
	if reflect.ValueOf(apConfig).IsZero() {
		cfg := svc.defaultConfig()
		svc.configRepo.Update(config.ACT_PUB_CONF_NAME, cfg)
		return cfg, nil
	}
	return apConfig, nil
}

func (s *ActivityPubService) AddFollower(follower string) error {
	return s.followersRepo.Add(follower)
}

func (s *ActivityPubService) RemoveFollower(follower string) error {
	return s.followersRepo.Remove(follower)
}

func (s *ActivityPubService) AllFollowers() ([]string, error) {
	return s.followersRepo.All()
}

func (s *ActivityPubService) sign(privateKey *rsa.PrivateKey, pubKeyId string, body []byte, r *http.Request) error {
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

func (s *ActivityPubService) GetActor(reqUrl string, fromGame string) (vocab.Actor, error) {
	c := http.Client{}

	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		slog.Error("parse error", "err", err)
		return vocab.Actor{}, err
	}

	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Set("Accept", "application/ld+json")
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	req.Header.Set("Host", parsedUrl.Host)

	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(config.ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	err = s.sign(apConfig.PrivateKey(), siteConfig.FullUrl+"/games/"+fromGame+"#main-key", nil, req)
	if err != nil {
		slog.Error("Signing error", "err", err)
		return vocab.Actor{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return vocab.Actor{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return vocab.Actor{}, err
	}

	item, err := vocab.UnmarshalJSON(data)
	if err != nil {
		return vocab.Actor{}, err
	}

	var actor vocab.Actor

	err = vocab.OnActor(item, func(o *vocab.Actor) error {
		actor = *o
		return nil
	})

	return actor, err
}

func (s *ActivityPubService) VerifySignature(r *http.Request, sender string) error {
	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(config.ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	actor, err := s.GetActor(sender, siteConfig.FullUrl+"/activitypub/actor")
	// actor does not have a pub key -> don't verify
	if actor.PublicKey.PublicKeyPem == "" {
		return nil
	}

	if err != nil {
		return err
	}
	block, _ := pem.Decode([]byte(actor.PublicKey.PublicKeyPem))
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	slog.Info("retrieved pub key of sender", "actor", actor, "pubKey", pubKey)

	verifier, err := httpsig.NewVerifier(r)
	if err != nil {
		return err
	}
	return verifier.Verify(pubKey, httpsig.RSA_SHA256)
}

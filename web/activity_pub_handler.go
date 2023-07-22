package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"

	"github.com/gofiber/fiber/v2"
)

const ACT_PUB_CONF_NAME = "activity_pub"

type ActivityPubServer struct {
	configRepo repository.ConfigRepository
}

type ActivityPubConfig struct {
	PreferredUsername string `owl:"inputType=text"`
	PublicKeyPem      string `owl:"inputType=text widget=textarea"`
	PrivateKeyPem     string `owl:"inputType=text widget=textarea"`
}

type WebfingerResponse struct {
	Subject string            `json:"subject"`
	Links   []ActivityPubLink `json:"links"`
}

type ActivityPubLink struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

type ActivityPubActor struct {
	Context []string `json:"@context"`

	ID                string `json:"id"`
	Type              string `json:"type"`
	PreferredUsername string `json:"preferredUsername"`
	Inbox             string `json:"inbox"`
	Oubox             string `json:"outbox"`
	Followers         string `json:"followers"`

	PublicKey ActivityPubPublicKey `json:"publicKey"`
}

type ActivityPubPublicKey struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

type ActivityPubOrderedCollection struct {
	Context []string `json:"@context"`

	ID         string `json:"id"`
	Type       string `json:"type"`
	TotalItems int    `json:"totalItems"`
	First      string `json:"first"`
	Last       string `json:"last"`
}

func NewActivityPubServer(configRepo repository.ConfigRepository) *ActivityPubServer {
	return &ActivityPubServer{
		configRepo: configRepo,
	}
}

func (s *ActivityPubServer) HandleWebfinger(ctx *fiber.Ctx) error {
	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	webfinger := WebfingerResponse{
		Subject: ctx.Query("resource"),

		Links: []ActivityPubLink{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: siteConfig.FullUrl + "/activitypub/actor",
			},
		},
	}

	return ctx.JSON(webfinger)

}

func (s *ActivityPubServer) Router(router fiber.Router) {
	router.Get("/actor", s.HandleActor)
}

func (s *ActivityPubServer) HandleActor(ctx *fiber.Ctx) error {
	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	actor := ActivityPubActor{
		Context: []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},

		ID:                siteConfig.FullUrl + "/activitypub/actor",
		Type:              "Person",
		PreferredUsername: apConfig.PreferredUsername,
		Inbox:             siteConfig.FullUrl + "/activitypub/inbox",
		Oubox:             siteConfig.FullUrl + "/activitypub/outbox",
		Followers:         siteConfig.FullUrl + "/activitypub/followers",

		PublicKey: ActivityPubPublicKey{
			ID:           siteConfig.FullUrl + "/activitypub/actor#main-key",
			Owner:        siteConfig.FullUrl + "/activitypub/actor",
			PublicKeyPem: apConfig.PublicKeyPem,
		},
	}

	return ctx.JSON(actor)
}

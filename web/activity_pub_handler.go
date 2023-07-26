package web

import (
	"net/url"
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"

	vocab "github.com/go-ap/activitypub"

	"github.com/gofiber/fiber/v2"
)

const ACT_PUB_CONF_NAME = "activity_pub"

type ActivityPubServer struct {
	configRepo   repository.ConfigRepository
	entryService *app.EntryService
}

type ActivityPubConfig struct {
	PreferredUsername string `owl:"inputType=text"`
	PublicKeyPem      string `owl:"inputType=text widget=textarea"`
	PrivateKeyPem     string `owl:"inputType=text widget=textarea"`
}

type WebfingerResponse struct {
	Subject string          `json:"subject"`
	Links   []WebfingerLink `json:"links"`
}

type WebfingerLink struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

func NewActivityPubServer(configRepo repository.ConfigRepository, entryService *app.EntryService) *ActivityPubServer {
	return &ActivityPubServer{
		configRepo:   configRepo,
		entryService: entryService,
	}
}

func (s *ActivityPubServer) HandleWebfinger(ctx *fiber.Ctx) error {
	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	webfinger := WebfingerResponse{
		Subject: ctx.Query("resource"),

		Links: []WebfingerLink{
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
	router.Get("/outbox", s.HandleOutbox)
}

func (s *ActivityPubServer) HandleActor(ctx *fiber.Ctx) error {
	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	actor := vocab.PersonNew(vocab.IRI(siteConfig.FullUrl + "/activitypub/actor"))
	actor.PreferredUsername = vocab.NaturalLanguageValues{{Value: vocab.Content(apConfig.PreferredUsername)}}
	actor.Inbox = vocab.IRI(siteConfig.FullUrl + "/activitypub/inbox")
	actor.Outbox = vocab.IRI(siteConfig.FullUrl + "/activitypub/outbox")
	actor.Followers = vocab.IRI(siteConfig.FullUrl + "/activitypub/followers")
	actor.PublicKey = vocab.PublicKey{
		ID:           vocab.ID(siteConfig.FullUrl + "/activitypub/actor#main-key"),
		Owner:        vocab.IRI(siteConfig.FullUrl + "/activitypub/actor"),
		PublicKeyPem: apConfig.PublicKeyPem,
	}

	data, err := actor.MarshalJSON()
	if err != nil {
		return err
	}
	ctx.Set("Content-Type", "application/activity+json")
	return ctx.Send(data)
}

func (s *ActivityPubServer) HandleOutbox(ctx *fiber.Ctx) error {
	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	entries, err := s.entryService.FindAllByType(nil, true, false)
	if err != nil {
		return err
	}

	items := make([]vocab.Item, len(entries))
	for i, entry := range entries {
		url, _ := url.JoinPath(siteConfig.FullUrl, "/posts/"+entry.ID()+"/")
		items[i] = *vocab.ActivityNew(vocab.IRI(url), vocab.CreateType, vocab.Object{
			ID:   vocab.ID(url),
			Type: vocab.ArticleType,
			Content: vocab.NaturalLanguageValues{
				{Value: vocab.Content(entry.Content())},
			},
		})
	}

	outbox := vocab.OrderedCollectionNew(vocab.IRI(siteConfig.FullUrl + "/activitypub/outbox"))
	outbox.TotalItems = uint(len(items))
	outbox.OrderedItems = items

	data, err := outbox.MarshalJSON()
	if err != nil {
		return err
	}
	ctx.Set("Content-Type", "application/activity+json")
	return ctx.Send(data)

}

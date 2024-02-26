package web

import (
	"net/url"
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"owl-blogs/render"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"

	"github.com/gofiber/fiber/v2"
)

const ACT_PUB_CONF_NAME = "activity_pub"

type ActivityPubServer struct {
	configRepo   repository.ConfigRepository
	entryService *app.EntryService
}

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

type WebfingerResponse struct {
	Subject string          `json:"subject"`
	Aliases []string        `json:"aliases"`
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

	domain, err := url.Parse(siteConfig.FullUrl)
	if err != nil {
		return err
	}

	subject := ctx.Query("resource", "")
	if subject != "acct:"+apConfig.PreferredUsername+"@"+domain.Host {
		return ctx.Status(404).JSON(nil)
	}

	webfinger := WebfingerResponse{
		Subject: subject,

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
	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
		jsonld.IRI(vocab.SecurityContextURI),
	).Marshal(actor)
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

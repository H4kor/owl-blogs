package web

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"owl-blogs/app"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type ActivityPubServer struct {
	siteConfigService *app.SiteConfigService
	apService         *app.ActivityPubService
	entryService      *app.EntryService
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

func NewActivityPubServer(siteConfigService *app.SiteConfigService, entryService *app.EntryService, apService *app.ActivityPubService) *ActivityPubServer {
	return &ActivityPubServer{
		siteConfigService: siteConfigService,
		entryService:      entryService,
		apService:         apService,
	}
}

func (s *ActivityPubServer) HandleWebfinger(ctx *fiber.Ctx) error {
	siteConfig, _ := s.siteConfigService.GetSiteConfig()
	apConfig, _ := s.apService.GetApConfig()

	domain, err := url.Parse(siteConfig.FullUrl)
	if err != nil {
		return err
	}

	subject := ctx.Query("resource", "")
	blogSubject := "acct:" + apConfig.PreferredUsername + "@" + domain.Host
	slog.Info("webfinger request", "for", subject, "required", blogSubject)
	if subject != blogSubject {
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
	router.Post("/inbox", s.HandleInbox)
	router.Get("/followers", s.HandleFollowers)
}

func (s *ActivityPubServer) HandleActor(ctx *fiber.Ctx) error {
	siteConfig, _ := s.siteConfigService.GetSiteConfig()
	apConfig, _ := s.apService.GetApConfig()

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
	siteConfig, _ := s.siteConfigService.GetSiteConfig()
	// apConfig, _ := s.apService.GetApConfig()

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

func (s *ActivityPubServer) processFollow(r *http.Request, act *vocab.Activity) error {
	follower := act.Actor.GetID().String()
	err := s.apService.VerifySignature(r, follower)
	if err != nil {
		slog.Error("wrong signature", "err", err)
		return err
	}
	err = s.apService.AddFollower(follower)
	if err != nil {
		return err
	}

	// go acpub.Accept(gameName, act)

	return nil
}

func (s *ActivityPubServer) processUndo(r *http.Request, act *vocab.Activity) error {
	follower := act.Actor.GetID().String()
	err := s.apService.VerifySignature(r, follower)
	if err != nil {
		slog.Error("wrong signature", "err", err)
		return err
	}
	err = s.apService.RemoveFollower(follower)
	if err != nil {
		return err
	}

	// go acpub.Accept(gameName, act)

	return nil
}

func (s *ActivityPubServer) HandleInbox(ctx *fiber.Ctx) error {
	// siteConfig, _ := s.siteConfigService.GetSiteConfig()
	// apConfig, _ := s.apService.GetApConfig()

	body := ctx.Request().Body()
	data, err := vocab.UnmarshalJSON(body)
	if err != nil {
		slog.Error("failed to parse request body", "body", body, "err", err)
		return err
	}

	err = vocab.OnActivity(data, func(act *vocab.Activity) error {
		slog.Info("activity retrieved", "activity", act, "type", act.Type)

		r, err := adaptor.ConvertRequest(ctx, true)
		if err != nil {
			return err
		}

		if act.Type == vocab.FollowType {
			return s.processFollow(r, act)
		}

		if act.Type == vocab.UndoType {
			slog.Info("processing undo")
			return s.processUndo(r, act)
		}
		return errors.New("only follow and undo actions supported")
	})
	return err

}

func (s *ActivityPubServer) HandleFollowers(ctx *fiber.Ctx) error {
	siteConfig, _ := s.siteConfigService.GetSiteConfig()
	// apConfig, _ := s.apService.GetApConfig()

	fs, err := s.apService.AllFollowers()
	if err != nil {
		return err
	}

	followers := vocab.Collection{}
	for _, f := range fs {
		followers.Append(vocab.IRI(f))
	}
	followers.TotalItems = uint(len(fs))
	followers.ID = vocab.IRI(siteConfig.FullUrl + "/activitypub/followers")
	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(followers)

	if err != nil {
		return err
	}
	ctx.Set("Content-Type", "application/activity+json")
	return ctx.Send(data)
}

package web

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"owl-blogs/app"
	"strings"

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
				Href: s.apService.ActorUrl(),
			},
		},
	}

	return ctx.JSON(webfinger)

}

func (s *ActivityPubServer) Router(router fiber.Router) {
	router.Get("/outbox", s.HandleOutbox)
	router.Post("/inbox", s.HandleInbox)
	router.Get("/followers", s.HandleFollowers)
}

func (s *ActivityPubServer) HandleActor(ctx *fiber.Ctx) error {
	accepts := (strings.Contains(string(ctx.Request().Header.Peek("Accept")), "application/activity+json") ||
		strings.Contains(string(ctx.Request().Header.Peek("Accept")), "application/ld+json"))
	req_content := (strings.Contains(string(ctx.Request().Header.Peek("Content-Type")), "application/activity+json") ||
		strings.Contains(string(ctx.Request().Header.Peek("Content-Type")), "application/ld+json"))
	if !accepts && !req_content {
		return ctx.Next()
	}
	apConfig, _ := s.apService.GetApConfig()

	actor := vocab.PersonNew(vocab.IRI(s.apService.ActorUrl()))
	actor.PreferredUsername = vocab.NaturalLanguageValues{{Value: vocab.Content(apConfig.PreferredUsername)}}
	actor.Inbox = vocab.IRI(s.apService.InboxUrl())
	actor.Outbox = vocab.IRI(s.apService.OutboxUrl())
	actor.Followers = vocab.IRI(s.apService.FollowersUrl())
	actor.PublicKey = vocab.PublicKey{
		ID:           vocab.IRI(s.apService.MainKeyUri()),
		Owner:        vocab.IRI(s.apService.ActorUrl()),
		PublicKeyPem: apConfig.PublicKeyPem,
	}

	actor.Name = vocab.NaturalLanguageValues{{Value: vocab.Content(s.apService.ActorName())}}
	actor.Icon = s.apService.ActorIcon()
	actor.Summary = vocab.NaturalLanguageValues{{Value: vocab.Content(s.apService.ActorSummary())}}

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

	outbox := vocab.OrderedCollectionNew(vocab.IRI(s.apService.OutboxUrl()))
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

	go s.apService.Accept(act)

	return nil
}

func (s *ActivityPubServer) processUndo(r *http.Request, act *vocab.Activity) error {
	sender := act.Actor.GetID().String()
	err := s.apService.VerifySignature(r, sender)

	return vocab.OnObject(act.Object, func(o *vocab.Object) error {
		if o.Type == vocab.FollowType {
			if err != nil {
				slog.Error("wrong signature", "err", err)
				return err
			}
			err = s.apService.RemoveFollower(sender)
			if err != nil {
				return err
			}
			go s.apService.Accept(act)
			return nil
		}
		if o.Type == vocab.LikeType {
			return s.apService.RemoveLike(o.ID.String())
		}
		if o.Type == vocab.AnnounceType {
			return s.apService.RemoveRepost(o.ID.String())
		}
		slog.Warn("unsupporeted object type for undo", "object", o)
		return errors.New("unsupporeted object type")
	})

}

func (s *ActivityPubServer) processLike(r *http.Request, act *vocab.Activity) error {
	sender := act.Actor.GetID().String()
	liked := act.Object.GetID().String()
	err := s.apService.VerifySignature(r, sender)
	if err != nil {
		slog.Error("wrong signature", "err", err)
		return err
	}

	err = s.apService.AddLike(sender, liked, act.ID.String())
	if err != nil {
		slog.Error("error saving like", "err", err)
		return err
	}

	go s.apService.Accept(act)
	return nil
}

func (s *ActivityPubServer) processAnnounce(r *http.Request, act *vocab.Activity) error {
	sender := act.Actor.GetID().String()
	liked := act.Object.GetID().String()
	err := s.apService.VerifySignature(r, sender)
	if err != nil {
		slog.Error("wrong signature", "err", err)
		return err
	}

	err = s.apService.AddRepost(sender, liked, act.ID.String())
	if err != nil {
		slog.Error("error saving like", "err", err)
		return err
	}

	go s.apService.Accept(act)
	return nil
}

func (s *ActivityPubServer) processDelete(r *http.Request, act *vocab.Activity) error {
	return vocab.OnObject(act.Object, func(o *vocab.Object) error {
		slog.Warn("Not processing delete", "action", act, "object", o)
		return nil
	})
}

func (s *ActivityPubServer) HandleInbox(ctx *fiber.Ctx) error {
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
			return s.processUndo(r, act)
		}
		if act.Type == vocab.DeleteType {
			return s.processDelete(r, act)
		}
		if act.Type == vocab.LikeType {
			return s.processLike(r, act)
		}
		if act.Type == vocab.AnnounceType {
			return s.processAnnounce(r, act)
		}

		slog.Warn("Unsupported action", "body", body)

		return errors.New("only follow and undo actions supported")
	})
	return err

}

func (s *ActivityPubServer) HandleFollowers(ctx *fiber.Ctx) error {
	fs, err := s.apService.AllFollowers()
	if err != nil {
		return err
	}

	followers := vocab.Collection{}
	for _, f := range fs {
		followers.Append(vocab.IRI(f))
	}
	followers.TotalItems = uint(len(fs))
	followers.ID = vocab.IRI(s.apService.FollowersUrl())
	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(followers)

	if err != nil {
		return err
	}
	ctx.Set("Content-Type", "application/activity+json")
	return ctx.Send(data)
}

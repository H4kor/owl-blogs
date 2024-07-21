package web

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"owl-blogs/app"
	"owl-blogs/config"

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

func (s *ActivityPubServer) HandleNodeInfo(ctx *fiber.Ctx) error {
	siteConfig, _ := s.siteConfigService.GetSiteConfig()
	href, _ := url.JoinPath(siteConfig.FullUrl, "/activitypub/nodeinfo/2.1")
	nodeinfo := map[string]interface{}{
		"links": []map[string]interface{}{
			{
				"rel":  "http://nodeinfo.diaspora.software/ns/schema/2.1",
				"href": href,
			},
		},
	}
	return ctx.JSON(nodeinfo)
}

func (s *ActivityPubServer) HandleNodeInfoDetails(ctx *fiber.Ctx) error {
	nodeinfo := map[string]interface{}{
		"version": "2.1",
		"software": map[string]interface{}{
			"name":       "owl-blogs",
			"version":    config.OWL_VERSION,
			"repository": "https://github.com/H4kor/owl-blogs",
		},
		"protocols": []string{
			"activitypub",
		},
		"services": map[string]interface{}{
			"inbound": []string{},
			"outbound": []string{
				"atom1.0", "rss2.0",
			},
		},
		"openRegistrations": false,
		"usage": map[string]interface{}{
			"users": map[string]interface{}{},
		},
	}
	return ctx.JSON(nodeinfo)
}

func (s *ActivityPubServer) Router(router fiber.Router) {
	router.Get("/outbox", s.HandleOutbox)
	router.Post("/inbox", s.HandleInbox)
	router.Get("/followers", s.HandleFollowers)
	router.Get("/nodeinfo/2.1", s.HandleNodeInfoDetails)
}

func (s *ActivityPubServer) HandleActor(ctx *fiber.Ctx) error {
	if !isActivityPub(ctx) {
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
		return s.handleError(ctx, err)
	}
	ctx.Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	return ctx.Send(data)
}

func (s *ActivityPubServer) HandleEntry(ctx *fiber.Ctx) error {
	if !isActivityPub(ctx) {
		return ctx.Next()
	}

	entryId := ctx.Params("post")
	entry, err := s.entryService.FindById(entryId)
	if err != nil {
		return err
	}

	obj, err := s.apService.EntryToObject(entry)
	if err != nil {
		if errors.Is(err, app.ErrEntryTypeNotSupported) || errors.Is(err, app.ErrEntryNotFound) {
			return ctx.SendStatus(404)
		}
		return err
	}

	data, err := app.ApEncoder.Marshal(obj)
	if err != nil {
		return s.handleError(ctx, err)
	}
	ctx.Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	return ctx.Send(data)
}

func (s *ActivityPubServer) HandleOutbox(ctx *fiber.Ctx) error {
	siteConfig, _ := s.siteConfigService.GetSiteConfig()
	// apConfig, _ := s.apService.GetApConfig()

	entries, err := s.entryService.FindAllByType(nil, true, false)
	if err != nil {
		return s.handleError(ctx, err)

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
		return s.handleError(ctx, err)

	}
	ctx.Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
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
	if err != nil {
		slog.Error("wrong signature", "err", err)
		return err
	}

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
		return app.ErrUnsupportedObjectType
	})
}

func (s *ActivityPubServer) processCreate(r *http.Request, act *vocab.Activity) error {
	sender := act.Actor.GetID().String()
	err := s.apService.VerifySignature(r, sender)
	if err != nil {
		slog.Error("wrong signature", "err", err)
		return err
	}

	return vocab.OnObject(act.Object, func(o *vocab.Object) error {
		if o.Type == vocab.NoteType {
			slog.Info("processing note")
			return s.apService.AddReply(sender, o.InReplyTo.GetID().String(), o.ID.String(), o.Content.String())
		}
		if o.Type == vocab.ArticleType {
			slog.Info("processing article")
			return s.apService.AddReply(sender, o.InReplyTo.GetID().String(), o.ID.String(), o.Name.String())
		}

		slog.Warn("Not processing craete", "action", act, "object", o)
		return nil
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
	sender := act.Actor.GetID().String()
	err := s.apService.VerifySignature(r, sender)
	if err != nil {
		slog.Error("wrong signature", "err", err)
		return err
	}

	err = vocab.OnObject(act.Object, func(o *vocab.Object) error {
		if o.Type == vocab.NoteType || o.Type == vocab.ArticleType {
			return s.apService.RemoveReply(o.ID.String())
		}

		slog.Warn("Not processing delete", "action", act, "object", o)
		return nil
	})
	// error can be because object is an IRI
	// this can be safely ignored -> log as warning and return nil
	slog.Warn("Error on procressDelete", "error", err)
	return nil
}

func (s *ActivityPubServer) HandleInbox(ctx *fiber.Ctx) error {
	body := ctx.Request().Body()
	data, err := vocab.UnmarshalJSON(body)
	if err != nil {
		slog.Error("failed to parse request body", "body", body, "err", err)
		return s.handleError(ctx, err)
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
		if act.Type == vocab.CreateType {
			return s.processCreate(r, act)
		}

		slog.Warn("Unsupported action", "body", body)

		return app.ErrUnsupportedActionType
	})
	return s.handleError(ctx, err)

}

func (s *ActivityPubServer) HandleFollowers(ctx *fiber.Ctx) error {
	fs, err := s.apService.AllFollowers()
	if err != nil {
		return s.handleError(ctx, err)
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
		return s.handleError(ctx, err)
	}
	ctx.Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	return ctx.Send(data)
}

func (s *ActivityPubServer) handleError(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	webErr, ok := err.(app.WebError)
	if ok {
		ctx.Status(webErr.Status())
		return ctx.JSON(map[string]string{
			"error": webErr.Error(),
		}, "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	}
	slog.Error("unhandled error", "error", err)
	return err
}

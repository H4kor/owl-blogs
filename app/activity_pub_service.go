package app

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"owl-blogs/interactions"
	"owl-blogs/render"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"github.com/go-fed/httpsig"
	"github.com/microcosm-cc/bluemonday"
)

var ApEncoder = jsonld.WithContext(
	jsonld.Context{
		{IRI: jsonld.IRI(vocab.ActivityBaseURI)},
		{Term: jsonld.Term("Hashtag"), IRI: jsonld.IRI("https://www.w3.org/ns/activitystreams#Hashtag")},
	},
)

type ActivityPubConfig struct {
	PreferredUsername string
	PublicKeyPem      string
	PrivateKeyPem     string
}

// Form implements app.AppConfig.
func (cfg *ActivityPubConfig) Form(binSvc model.BinaryStorageInterface) template.HTML {
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
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		slog.Error("error x509.ParsePKCS1PrivateKey", "err", err)
	}
	return privKey
}

type ActivityPubService struct {
	followersRepo         repository.FollowerRepository
	configRepo            repository.ConfigRepository
	interactionRepository repository.InteractionRepository
	entryService          *EntryService
	siteConfigServcie     *SiteConfigService
	binService            *BinaryService
}

func NewActivityPubService(
	followersRepo repository.FollowerRepository,
	configRepo repository.ConfigRepository,
	interactionRepository repository.InteractionRepository,
	entryService *EntryService,
	siteConfigServcie *SiteConfigService,
	binService *BinaryService,
	bus *EventBus,
) *ActivityPubService {
	service := &ActivityPubService{
		followersRepo:         followersRepo,
		configRepo:            configRepo,
		interactionRepository: interactionRepository,
		entryService:          entryService,
		binService:            binService,
		siteConfigServcie:     siteConfigServcie,
	}

	bus.Subscribe(service)

	return service
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
		println("ERROR IN ACTIVITY PUB CONFIG", err.Error())
		return ActivityPubConfig{}, err
	}
	if reflect.ValueOf(apConfig).IsZero() {
		cfg := svc.defaultConfig()
		svc.configRepo.Update(config.ACT_PUB_CONF_NAME, cfg)
		return cfg, nil
	}
	return apConfig, nil
}

func (svc *ActivityPubService) ActorUrl() string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.FullUrl
}
func (svc *ActivityPubService) MainKeyUri() string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.FullUrl + "#main-key"
}
func (svc *ActivityPubService) InboxUrl() string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.FullUrl + "/activitypub/inbox"
}
func (svc *ActivityPubService) OutboxUrl() string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.FullUrl + "/activitypub/outbox"
}
func (svc *ActivityPubService) FollowersUrl() string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.FullUrl + "/activitypub/followers"
}
func (svc *ActivityPubService) AcccepId() string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.FullUrl + "#accept-" + strconv.FormatInt(time.Now().UnixNano(), 16)
}

func (svc *ActivityPubService) HashtagId(tag string) string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.FullUrl + "/tags/" + tag + "/"
}

func (svc *ActivityPubService) ActorName() string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.Title
}

func (svc *ActivityPubService) ActorIcon() vocab.Image {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	u := cfg.AvatarUrl
	pUrl, _ := url.Parse(u)
	parts := strings.Split(pUrl.Path, ".")
	fullUrl, _ := url.JoinPath(cfg.FullUrl, u)
	return vocab.Image{
		Type:      vocab.ImageType,
		MediaType: vocab.MimeType("image/" + parts[len(parts)-1]),
		URL:       vocab.IRI(fullUrl),
	}
}

func (svc *ActivityPubService) ActorSummary() string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.SubTitle
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
	return err
}

func (s *ActivityPubService) GetActor(reqUrl string) (vocab.Actor, error) {

	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(config.ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

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

	err = s.sign(apConfig.PrivateKey(), s.MainKeyUri(), nil, req)
	if err != nil {
		slog.Error("Signing error", "err", err)
		return vocab.Actor{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		slog.Error("failed to retrieve sender actor", "err", err, "url", reqUrl)
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

func (s *ActivityPubService) GetObject(objUrl string) (vocab.Object, error) {

	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(config.ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	c := http.Client{}

	parsedUrl, err := url.Parse(objUrl)
	if err != nil {
		slog.Error("parse error", "err", err)
		return vocab.Object{}, err
	}

	req, _ := http.NewRequest("GET", objUrl, nil)
	req.Header.Set("Accept", "application/ld+json")
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	req.Header.Set("Host", parsedUrl.Host)

	err = s.sign(apConfig.PrivateKey(), s.MainKeyUri(), nil, req)
	if err != nil {
		slog.Error("Signing error", "err", err)
		return vocab.Object{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		slog.Error("failed to retrieve object", "err", err, "url", objUrl)
		return vocab.Object{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return vocab.Object{}, err
	}

	item, err := vocab.UnmarshalJSON(data)
	if err != nil {
		return vocab.Object{}, err
	}

	var obj vocab.Object

	err = vocab.OnObject(item, func(o *vocab.Object) error {
		obj = *o
		return nil
	})

	return obj, err
}

func (s *ActivityPubService) VerifySignature(r *http.Request, sender string) error {
	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(config.ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	slog.Info("verifying for", "sender", sender, "retriever", s.ActorUrl())

	actor, err := s.GetActor(sender)
	// actor does not have a pub key -> don't verify
	if actor.PublicKey.PublicKeyPem == "" {
		return nil
	}

	if err != nil {
		slog.Error("unable to retrieve actor for sig verification", "sender", sender)
		return ErrUnableToVerifySignature
	}
	block, _ := pem.Decode([]byte(actor.PublicKey.PublicKeyPem))
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		pubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return err
		}
	}
	slog.Info("retrieved pub key of sender", "actor", actor, "pubKey", pubKey)

	verifier, err := httpsig.NewVerifier(r)
	if err != nil {
		slog.Error("invalid signature", "err", err)
		return ErrUnableToVerifySignature
	}
	if verifier.Verify(pubKey, httpsig.RSA_SHA256) != nil {
		return ErrUnableToVerifySignature
	}
	return nil
}

func (s *ActivityPubService) Accept(act *vocab.Activity) error {
	actor, err := s.GetActor(act.Actor.GetID().String())
	if err != nil {
		return ErrUnableToVerifySignature
	}

	accept := vocab.AcceptNew(vocab.IRI(s.AcccepId()), act)
	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(accept)

	if err != nil {
		slog.Error("marshalling error", "err", err)
		return ErrUnableToVerifySignature
	}

	return s.sendObject(actor, data)
}

func (s *ActivityPubService) AddLike(sender string, liked string, likeId string) error {
	entry, err := s.entryService.FindByUrl(liked)
	if err != nil {
		return err
	}

	actor, err := s.GetActor(sender)
	if err != nil {
		return err
	}

	var like *interactions.Like
	interaction, err := s.interactionRepository.FindById(likeId)
	if err != nil {
		interaction = &interactions.Like{}
	}
	like, ok := interaction.(*interactions.Like)
	if !ok {
		return ErrConflictingId
	}
	existing := like.ID() != ""

	likeMeta := interactions.LikeMetaData{
		SenderUrl:  sender,
		SenderName: actor.Name.String(),
	}
	like.SetID(likeId)
	like.SetMetaData(&likeMeta)
	like.SetEntryID(entry.ID())
	like.SetCreatedAt(time.Now())
	if !existing {
		return s.interactionRepository.Create(like)
	} else {
		return s.interactionRepository.Update(like)
	}
}

func (s *ActivityPubService) RemoveLike(id string) error {
	interaction, err := s.interactionRepository.FindById(id)
	if err != nil {
		return err
	}
	return s.interactionRepository.Delete(interaction)
}

func (s *ActivityPubService) AddRepost(sender string, reposted string, respostId string) error {
	entry, err := s.entryService.FindByUrl(reposted)
	if err != nil {
		return err
	}

	actor, err := s.GetActor(sender)
	if err != nil {
		return err
	}

	var repost *interactions.Repost
	interaction, err := s.interactionRepository.FindById(respostId)
	if err != nil {
		interaction = &interactions.Repost{}
	}
	repost, ok := interaction.(*interactions.Repost)
	if !ok {
		return ErrConflictingId
	}
	existing := repost.ID() != ""

	repostMeta := interactions.RepostMetaData{
		SenderUrl:  sender,
		SenderName: actor.Name.String(),
	}
	repost.SetID(respostId)
	repost.SetMetaData(&repostMeta)
	repost.SetEntryID(entry.ID())
	repost.SetCreatedAt(time.Now())
	if !existing {
		return s.interactionRepository.Create(repost)
	} else {
		return s.interactionRepository.Update(repost)
	}
}

func (s *ActivityPubService) RemoveRepost(id string) error {
	interaction, err := s.interactionRepository.FindById(id)
	if err != nil {
		return err
	}
	return s.interactionRepository.Delete(interaction)
}

func (s *ActivityPubService) AddReply(sender string, replyTo string, replyId string, replyContent string) error {
	entry, err := s.entryService.FindByUrl(replyTo)
	if err != nil {
		return err
	}

	actor, err := s.GetActor(sender)
	if err != nil {
		return err
	}

	var reply *interactions.Reply
	interaction, err := s.interactionRepository.FindById(replyId)
	if err != nil {
		interaction = &interactions.Reply{}
	}
	reply, ok := interaction.(*interactions.Reply)
	if !ok {
		return ErrConflictingId
	}
	existing := reply.ID() != ""

	// clean reply to list of all allowed html tags
	p := bluemonday.NewPolicy()
	// Require URLs to be parseable by net/url.Parse and either:
	//   mailto: http:// or https://
	p.AllowStandardURLs()
	// We only allow <p> and <a href="">
	p.AllowAttrs("href").OnElements("a")
	p.AllowElements("p")
	p.AllowElements("p")
	p.AllowElements("span")
	p.AllowElements("br")
	p.AllowElements("del")
	p.AllowElements("pre")
	p.AllowElements("code")
	p.AllowElements("em")
	p.AllowElements("strong")
	p.AllowElements("b")
	p.AllowElements("i")
	p.AllowElements("u")
	p.AllowElements("ul")
	p.AllowElements("ol")
	p.AllowElements("li")
	p.AllowElements("blockquote")
	cleanReplyContent := p.Sanitize(replyContent)

	replyMeta := interactions.ReplyMetaData{
		SenderUrl:   sender,
		SenderName:  actor.Name.String(),
		OriginalUrl: replyId,
		Content:     template.HTML(cleanReplyContent),
	}
	reply.SetID(replyId)
	reply.SetMetaData(&replyMeta)
	reply.SetEntryID(entry.ID())
	reply.SetCreatedAt(time.Now())
	if !existing {
		return s.interactionRepository.Create(reply)
	} else {
		return s.interactionRepository.Update(reply)
	}
}

func (s *ActivityPubService) RemoveReply(id string) error {
	interaction, err := s.interactionRepository.FindById(id)
	if err != nil {
		return err
	}
	return s.interactionRepository.Delete(interaction)
}

func (s *ActivityPubService) sendObject(to vocab.Actor, data []byte) error {
	siteConfig := model.SiteConfig{}
	apConfig := ActivityPubConfig{}
	s.configRepo.Get(config.ACT_PUB_CONF_NAME, &apConfig)
	s.configRepo.Get(config.SITE_CONFIG, &siteConfig)

	if to.Inbox == nil {
		slog.Error("actor has no inbox", "actor", to)
		return ErrNoActorInbox
	}

	actorUrl, err := url.Parse(to.Inbox.GetID().String())
	if err != nil {
		slog.Error("parse error", "err", err)
		return err
	}

	c := http.Client{}
	req, _ := http.NewRequest("POST", to.Inbox.GetID().String(), bytes.NewReader(data))
	req.Header.Set("Accept", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
    req.Header.Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	req.Header.Set("Host", actorUrl.Host)
	err = s.sign(apConfig.PrivateKey(), s.MainKeyUri(), data, req)
	if err != nil {
		slog.Error("Signing error", "err", err)
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		slog.Error("Sending error", "url", req.URL, "err", err)
		return err
	}
	slog.Info("Request", "host", resp.Request.Header)

	if resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Error sending Note", "method", resp.Request.Method, "url", resp.Request.URL, "status", resp.Status, "body", string(body))
		return err
	}
	body, _ := io.ReadAll(resp.Body)
	slog.Info("Sent Body", "body", string(data))
	slog.Info("Retrieved", "status", resp.Status, "body", string(body))
	return nil
}

/*
 * Notifiers
 */

func (svc *ActivityPubService) NotifyEntryCreated(entry model.Entry) {
	slog.Info("Processing Entry Create for ActivityPub")
	followers, err := svc.AllFollowers()
	if err != nil {
		slog.Error("Cannot retrieve followers")
	}

	object, err := svc.EntryToObject(entry)
	if err != nil {
		slog.Error("Cannot convert object", "err", err)
	}

	create := vocab.CreateNew(object.ID, object)
	create.Actor = object.AttributedTo
	create.To = object.To
	create.CC = object.CC
	create.Published = object.Published
	data, err := ApEncoder.Marshal(create)
	if err != nil {
		slog.Error("marshalling error", "err", err)
	}

	for _, follower := range followers {
		actor, err := svc.GetActor(follower)
		if err != nil {
			slog.Error("Unable to retrieve follower actor", "err", err)
		}
		svc.sendObject(actor, data)
	}
	for _, mentioned := range slices.Concat(create.To, create.CC) {
        href := string(mentioned.GetID())
        if href == vocab.PublicNS.String() || href == svc.FollowersUrl() {
            continue
        }
		actor, err := svc.GetActor(href)
		if err != nil {
			slog.Error("Unable to retrieve mentioned actor", "err", err)
		}
		svc.sendObject(actor, data)
	}
}

func (svc *ActivityPubService) NotifyEntryUpdated(entry model.Entry) {
	slog.Info("Processing Entry Create for ActivityPub")
	followers, err := svc.AllFollowers()
	if err != nil {
		slog.Error("Cannot retrieve followers")
	}

	object, err := svc.EntryToObject(entry)
	if err != nil {
		slog.Error("Cannot convert object", "err", err)
	}

	update := vocab.CreateNew(object.ID, object)
	update.Actor = object.AttributedTo
	update.To = object.To
	update.Published = object.Published
	data, err := ApEncoder.Marshal(update)
	if err != nil {
		slog.Error("marshalling error", "err", err)
	}

	for _, follower := range followers {
		actor, err := svc.GetActor(follower)
		if err != nil {
			slog.Error("Unable to retrieve follower actor", "err", err)
		}
		svc.sendObject(actor, data)
	}

	for _, mentioned := range slices.Concat(update.To, update.CC) {
        href := string(mentioned.GetID())
        if href == vocab.PublicNS.String() || href == svc.FollowersUrl() {
            continue
        }
		actor, err := svc.GetActor(href)
		if err != nil {
			slog.Error("Unable to retrieve mentioned actor", "err", err)
		}
		svc.sendObject(actor, data)
	}

}

func (svc *ActivityPubService) NotifyEntryDeleted(entry model.Entry) {
	obj, err := svc.EntryToObject(entry)
	if err != nil {
		slog.Error("error converting to object", "err", err)
		return
	}

	followers, err := svc.AllFollowers()
	if err != nil {
		slog.Error("Cannot retrieve followers")
	}

	delete := vocab.DeleteNew(obj.ID, obj)
	delete.Actor = obj.AttributedTo
	delete.To = obj.To
	delete.CC = obj.CC
	delete.Published = time.Now()
	data, err := ApEncoder.Marshal(delete)
	if err != nil {
		slog.Error("marshalling error", "err", err)
	}

	for _, follower := range followers {
		actor, err := svc.GetActor(follower)
		if err != nil {
			slog.Error("Unable to retrieve follower actor", "err", err)
		}
		svc.sendObject(actor, data)
	}
	for _, mentioned := range slices.Concat(delete.To, delete.CC) {
        href := string(mentioned.GetID())
        if href == vocab.PublicNS.String() || href == svc.FollowersUrl() {
            continue
        }
		actor, err := svc.GetActor(href)
		if err != nil {
			slog.Error("Unable to retrieve mentioned actor", "err", err)
		}
		svc.sendObject(actor, data)
	}
}

func (svc *ActivityPubService) collectRetrieversOfObject(objId string) ([]string, error) {
    obj, err := svc.GetObject(objId)
    if err != nil {
        return []string{}, err
    }
    
    retrievers := make(map[string]bool)
    
    if obj.AttributedTo != nil && obj.AttributedTo.GetID() != "" {
        retrievers[string(obj.AttributedTo.GetID())] = true
    }
    for _, to := range slices.Concat(obj.To, obj.CC) {
        toId := string(to.GetID())
        if toId != string(vocab.PublicNS) {
            slog.Info("found retriever in object", "obj", obj.ID, "retriever", toId)
            retrievers[toId] = true
        }
    }

    ids := make([]string, 0, len(retrievers))
    for k := range retrievers {
        ids = append(ids, k)
    }

    return ids, err
}

func (svc *ActivityPubService) isActor(id string) bool {
    obj, err := svc.GetObject(id)
    if err != nil {
        slog.Error("could not get object to check if isActor", "object", id, "err", err)
        return false
    }
    isActor := slices.Contains(vocab.ActorTypes, obj.Type)
    slog.Info("isActor", "id", id, "types", vocab.ActorTypes, "objType", obj.Type)
    return isActor
}

func (svc *ActivityPubService) processMentions(obj *vocab.Object, entry model.Entry) error {
    retrievers := make(map[string]bool, 0)
    mentions := make(map[string]bool, 0)

    links, err := ParseLinksFromString(string(entry.Content()))
    slog.Info("Parsed links of entry", "entry", entry.ID(), "num_links", len(links))
    if err != nil {
        slog.Error("Unable to parse links form entry", "entry", entry.ID(), "err", err)
    }
    for _, link := range links {
            slog.Info("Found link in entry", "link", link)
            mentionedObj, err := svc.GetObject(link)
            if err == nil {
                // case Post mentioned
                if mentionedObj.AttributedTo != nil {
                    mentionedActor := mentionedObj.AttributedTo.GetID()
                    slog.Info("Adding mentioned object", "object", mentionedObj.ID, "actor", mentionedActor)
                    retrievers[mentionedActor.GetID().String()] = true
                    mentions[mentionedObj.ID.String()] = true
                    
                    // include all to,cc of mentioned object in cc
                    for _, to := range slices.Concat(mentionedObj.To, mentionedObj.CC) {
                        slog.Info("found retriever in mentioned object", "obj", mentionedObj.ID.String(), "retriever", to.GetID().String())
                        retrievers[to.GetID().String()] = true
                    }

                } else {
                    slog.Info("Adding actor based on link", "actor", mentionedObj)
                    retrievers[mentionedObj.ID.String()] = true
                    mentions[mentionedObj.ID.String()] = true
                }
            } else {
                slog.Info("Unable to get linked object", "err", err)
            }
    }

	if obj.InReplyTo != nil && obj.InReplyTo.GetID() != "" {
        replyRets, err := svc.collectRetrieversOfObject(obj.InReplyTo.GetID().String())
        if err != nil {
            return err
        }
        for _, x := range replyRets {
            retrievers[x] = true
            mentions[x] = true
        }
    }

    for to := range retrievers {
        // remove own followers and public NS
        if to == vocab.PublicNS.String() || to == svc.FollowersUrl() {
            continue
        }
        // only actors should be listed in CC
        if !svc.isActor(to) {
            slog.Info("removing as not actor", "retriever", to)
            continue
        }
        obj.CC = append(obj.CC, vocab.ID(to))
    }
    for to := range mentions {
        if to == vocab.PublicNS.String() || to == svc.FollowersUrl() {
            continue
        }
        mention := vocab.MentionNew(
            vocab.IRI(to),
        )
        mention.Href = vocab.ID(to)
        obj.Tag = append(obj.Tag, mention)
    }

    if len(mentions) == 1 && obj.InReplyTo == nil {
        for k := range mentions{
            slog.Info("entry is only mentioning one object, setting as inReplyTo", "replyTo", k)
            obj.InReplyTo = vocab.IRI(k)
        }
    }
    return nil
}


func (svc *ActivityPubService) EntryToObject(entry model.Entry) (vocab.Object, error) {
	// limit to notes for now

	if activityEntry, ok := entry.(ToActivityPub); ok {
		siteCfg, _ := svc.siteConfigServcie.GetSiteConfig()
		obj := activityEntry.ActivityObject(siteCfg, *svc.binService)
		obj.ID = vocab.ID(entry.FullUrl(siteCfg))
		obj.To = vocab.ItemCollection{
			vocab.PublicNS,
			vocab.IRI(svc.FollowersUrl()),
		}
        obj.CC = vocab.ItemCollection{}
		obj.AttributedTo = vocab.ID(svc.ActorUrl())

		for _, tag := range entry.Tags() {
			hashtag := vocab.LinkNew("", "Hashtag")
			hashtag.Type = "Hashtag"
			hashtag.Href = vocab.IRI(svc.HashtagId(tag))
			hashtag.Name = vocab.NaturalLanguageValues{
				{Value: vocab.Content("#" + tag)},
			}
			obj.Tag = append(obj.Tag, hashtag)
		}

        svc.processMentions(&obj, entry)
		return obj, nil
	}

	slog.Warn("entry type not yet supported for activity pub")
	return vocab.Object{}, ErrEntryTypeNotSupported
}

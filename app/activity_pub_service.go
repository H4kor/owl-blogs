package app

import (
	"bytes"
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
	entrytypes "owl-blogs/entry_types"
	"owl-blogs/interactions"
	"owl-blogs/render"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
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

func (svc *ActivityPubService) HashtagId(hashtag string) string {
	cfg, _ := svc.siteConfigServcie.GetSiteConfig()
	return cfg.FullUrl + "/tags/" + strings.ReplaceAll(hashtag, "#", "")
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

	slog.Info("Signed Request", "req", r.Header)
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
		slog.Error("unable to decode pub key pem", "pubKeyPem", actor.PublicKey.PublicKeyPem)
		return ErrUnableToVerifySignature
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
		interaction = &interactions.Like{}
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
		interaction = &interactions.Repost{}
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
	req.Header.Set("Accept", "application/ld+json")
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

	object, err := svc.entryToObject(entry)
	if err != nil {
		slog.Error("Cannot convert object", "err", err)
	}

	create := vocab.CreateNew(object.ID, object)
	create.Actor = object.AttributedTo
	create.To = object.To
	create.Published = object.Published
	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
		jsonld.Context{
			jsonld.ContextElement{
				Term: "toot",
				IRI:  jsonld.IRI("http://joinmastodon.org/ns#"),
			},
		},
	).Marshal(create)
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
}

func (svc *ActivityPubService) NotifyEntryUpdated(entry model.Entry) {
	slog.Info("Processing Entry Create for ActivityPub")
	followers, err := svc.AllFollowers()
	if err != nil {
		slog.Error("Cannot retrieve followers")
	}

	object, err := svc.entryToObject(entry)
	if err != nil {
		slog.Error("Cannot convert object", "err", err)
	}

	create := vocab.CreateNew(object.ID, object)
	create.Actor = object.AttributedTo
	create.To = object.To
	create.Published = object.Published
	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
		jsonld.Context{
			jsonld.ContextElement{
				Term: "toot",
				IRI:  jsonld.IRI("http://joinmastodon.org/ns#"),
			},
		},
	).Marshal(create)
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
}

func (svc *ActivityPubService) NotifyEntryDeleted(entry model.Entry) {
	obj, err := svc.entryToObject(entry)
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
	delete.Published = time.Now()
	data, err := jsonld.WithContext(
		jsonld.IRI(vocab.ActivityBaseURI),
	).Marshal(delete)
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

}

func (svc *ActivityPubService) entryToObject(entry model.Entry) (vocab.Object, error) {
	// limit to notes for now

	if noteEntry, ok := entry.(*entrytypes.Note); ok {
		return svc.noteToObject(noteEntry), nil
	}
	if imageEntry, ok := entry.(*entrytypes.Image); ok {
		return svc.imageToObject(imageEntry), nil
	}
	if articleEntry, ok := entry.(*entrytypes.Article); ok {
		return svc.articleToObject(articleEntry), nil
	}
	if recipeEntry, ok := entry.(*entrytypes.Recipe); ok {
		return svc.recipeToObject(recipeEntry), nil
	}

	slog.Warn("entry type not yet supported for activity pub")
	return vocab.Object{}, ErrEntryTypeNotSupported
}

func (svc *ActivityPubService) noteToObject(noteEntry *entrytypes.Note) vocab.Object {

	siteCfg, _ := svc.siteConfigServcie.GetSiteConfig()
	content := noteEntry.Content()
	r := regexp.MustCompile("#[a-z0-9_]+")
	matches := r.FindAllString(string(content), -1)
	tags := vocab.ItemCollection{}
	for _, hashtag := range matches {
		tags.Append(vocab.Object{
			ID:   vocab.ID(svc.HashtagId(hashtag)),
			Name: vocab.NaturalLanguageValues{{Value: vocab.Content(hashtag)}},
		})
	}

	note := vocab.Note{
		ID:   vocab.ID(noteEntry.FullUrl(siteCfg)),
		Type: "Note",
		To: vocab.ItemCollection{
			vocab.PublicNS,
			vocab.IRI(svc.FollowersUrl()),
		},
		Published:    *noteEntry.PublishedAt(),
		AttributedTo: vocab.ID(svc.ActorUrl()),
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(content)},
		},
		Tag: tags,
	}
	return note

}

func (svc *ActivityPubService) imageToObject(imageEntry *entrytypes.Image) vocab.Object {
	siteCfg, _ := svc.siteConfigServcie.GetSiteConfig()
	content := imageEntry.Content()

	imgPath := imageEntry.ImageUrl()
	fullImageUrl, _ := url.JoinPath(siteCfg.FullUrl, imgPath)
	binaryFile, err := svc.binService.FindById(imageEntry.MetaData().(*entrytypes.ImageMetaData).ImageId)
	if err != nil {
		slog.Error("cannot get image file")
	}

	attachments := vocab.ItemCollection{}
	attachments = append(attachments, vocab.Document{
		Type:      vocab.DocumentType,
		MediaType: vocab.MimeType(binaryFile.Mime()),
		URL:       vocab.ID(fullImageUrl),
		Name: vocab.NaturalLanguageValues{
			{Value: vocab.Content(content)},
		},
	})

	image := vocab.Image{
		ID:   vocab.ID(imageEntry.FullUrl(siteCfg)),
		Type: "Image",
		To: vocab.ItemCollection{
			vocab.PublicNS,
			vocab.IRI(svc.FollowersUrl()),
		},
		Published:    *imageEntry.PublishedAt(),
		AttributedTo: vocab.ID(svc.ActorUrl()),
		Name: vocab.NaturalLanguageValues{
			{Value: vocab.Content(imageEntry.Title())},
		},
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(imageEntry.Title() + "<br><br>" + string(content))},
		},
		Attachment: attachments,
		// Tag: tags,
	}
	return image

}

func (svc *ActivityPubService) articleToObject(articleEntry *entrytypes.Article) vocab.Object {
	siteCfg, _ := svc.siteConfigServcie.GetSiteConfig()
	content := articleEntry.Content()

	image := vocab.Article{
		ID:   vocab.ID(articleEntry.FullUrl(siteCfg)),
		Type: "Article",
		To: vocab.ItemCollection{
			vocab.PublicNS,
			vocab.IRI(svc.FollowersUrl()),
		},
		Published:    *articleEntry.PublishedAt(),
		AttributedTo: vocab.ID(svc.ActorUrl()),
		Name: vocab.NaturalLanguageValues{
			{Value: vocab.Content(articleEntry.Title())},
		},
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(string(content))},
		},
	}
	return image

}

func (svc *ActivityPubService) recipeToObject(recipeEntry *entrytypes.Recipe) vocab.Object {
	siteCfg, _ := svc.siteConfigServcie.GetSiteConfig()
	content := recipeEntry.Content()

	image := vocab.Article{
		ID:   vocab.ID(recipeEntry.FullUrl(siteCfg)),
		Type: "Article",
		To: vocab.ItemCollection{
			vocab.PublicNS,
			vocab.IRI(svc.FollowersUrl()),
		},
		Published:    *recipeEntry.PublishedAt(),
		AttributedTo: vocab.ID(svc.ActorUrl()),
		Name: vocab.NaturalLanguageValues{
			{Value: vocab.Content(recipeEntry.Title())},
		},
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(string(content))},
		},
	}
	return image

}

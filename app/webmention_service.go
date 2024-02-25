package app

import (
	"fmt"
	"net/url"
	"owl-blogs/app/owlhttp"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/interactions"
	"time"
)

type WebmentionService struct {
	siteConfigService     *SiteConfigService
	InteractionRepository repository.InteractionRepository
	EntryRepository       repository.EntryRepository
	Http                  owlhttp.HttpClient
}

func NewWebmentionService(
	siteConfigService *SiteConfigService,
	interactionRepository repository.InteractionRepository,
	entryRepository repository.EntryRepository,
	http owlhttp.HttpClient,
	bus *EventBus,
) *WebmentionService {
	svc := &WebmentionService{
		siteConfigService:     siteConfigService,
		InteractionRepository: interactionRepository,
		EntryRepository:       entryRepository,
		Http:                  http,
	}
	bus.Subscribe(svc)
	return svc
}

func (s *WebmentionService) GetExistingWebmention(entryId string, source string, target string) (*interactions.Webmention, error) {
	inters, err := s.InteractionRepository.FindAll(entryId)
	if err != nil {
		return nil, err
	}
	for _, interaction := range inters {
		if webm, ok := interaction.(*interactions.Webmention); ok {
			m := webm.MetaData().(*interactions.WebmentionMetaData)
			if m.Source == source && m.Target == target {
				return webm, nil
			}
		}
	}
	return nil, nil
}

func (s *WebmentionService) ProcessWebmention(source string, target string) error {
	resp, err := s.Http.Get(source)
	if err != nil {
		return err
	}

	hEntry, err := ParseHEntry(resp)
	if err != nil {
		return err
	}

	entryId := UrlToEntryId(target)
	_, err = s.EntryRepository.FindById(entryId)
	if err != nil {
		return err
	}

	webmention, err := s.GetExistingWebmention(entryId, source, target)
	if err != nil {
		return err
	}
	if webmention != nil {
		data := interactions.WebmentionMetaData{
			Source: source,
			Target: target,
			Title:  hEntry.Title,
		}
		webmention.SetMetaData(&data)
		webmention.SetEntryID(entryId)
		webmention.SetCreatedAt(time.Now())
		err = s.InteractionRepository.Update(webmention)
		return err
	} else {
		webmention = &interactions.Webmention{}
		data := interactions.WebmentionMetaData{
			Source: source,
			Target: target,
			Title:  hEntry.Title,
		}
		webmention.SetMetaData(&data)
		webmention.SetEntryID(entryId)
		webmention.SetCreatedAt(time.Now())
		err = s.InteractionRepository.Create(webmention)
		return err
	}
}

func (s *WebmentionService) ScanForLinks(entry model.Entry) ([]string, error) {
	content := string(entry.Content())
	return ParseLinksFromString(content)
}

func (s *WebmentionService) FullEntryUrl(entry model.Entry) string {
	siteConfig, _ := s.siteConfigService.GetSiteConfig()

	url, _ := url.JoinPath(
		siteConfig.FullUrl,
		fmt.Sprintf("/posts/%s/", entry.ID()),
	)
	return url
}

func (s *WebmentionService) SendWebmention(entry model.Entry) error {
	links, err := s.ScanForLinks(entry)
	if err != nil {
		return err
	}
	for _, target := range links {
		resp, err := s.Http.Get(target)
		if err != nil {
			continue
		}
		endpoint, err := GetWebmentionEndpoint(resp)
		if err != nil {
			continue
		}
		payload := url.Values{}
		payload.Set("source", s.FullEntryUrl(entry))
		payload.Set("target", target)
		_, err = s.Http.PostForm(endpoint, payload)
		if err != nil {
			continue
		}
		println("Send webmention for target", target)
	}
	return nil
}

func (s *WebmentionService) NotifyEntryCreated(entry model.Entry) {
	s.SendWebmention(entry)
}

func (s *WebmentionService) NotifyEntryUpdated(entry model.Entry) {
	s.SendWebmention(entry)
}

func (s *WebmentionService) NotifyEntryDeleted(entry model.Entry) {
	s.SendWebmention(entry)
}

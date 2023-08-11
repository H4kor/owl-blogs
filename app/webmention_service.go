package app

import (
	"owl-blogs/app/owlhttp"
	"owl-blogs/app/repository"
	"owl-blogs/interactions"
	"time"
)

type WebmentionService struct {
	InteractionRepository repository.InteractionRepository
	EntryRepository       repository.EntryRepository
	Http                  owlhttp.HttpClient
}

func NewWebmentionService(
	interactionRepository repository.InteractionRepository,
	entryRepository repository.EntryRepository,
	http owlhttp.HttpClient,
	bus *EventBus,
) *WebmentionService {
	svc := &WebmentionService{
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

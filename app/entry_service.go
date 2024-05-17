package app

import (
	"errors"
	"fmt"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"regexp"
	"strings"
)

type EntryService struct {
	EntryRepository   repository.EntryRepository
	siteConfigServcie *SiteConfigService
	Bus               *EventBus
}

func NewEntryService(
	entryRepository repository.EntryRepository,
	siteConfigServcie *SiteConfigService,
	bus *EventBus,
) *EntryService {
	return &EntryService{
		EntryRepository:   entryRepository,
		siteConfigServcie: siteConfigServcie,
		Bus:               bus,
	}
}

func (s *EntryService) Create(entry model.Entry) error {
	// try to find a good ID
	m := regexp.MustCompile(`[^a-z0-9-]`)
	prefix := m.ReplaceAllString(strings.ToLower(entry.Title()), "-")
	title := prefix
	counter := 0
	for {
		_, err := s.EntryRepository.FindById(title)
		if err == nil {
			counter += 1
			title = prefix + "-" + fmt.Sprintf("%s-%d", prefix, counter)
		} else {
			break
		}
	}
	entry.SetID(title)

	err := s.EntryRepository.Create(entry)
	if err != nil {
		return err
	}
	s.Bus.NotifyCreated(entry)
	return nil
}

func (s *EntryService) Update(entry model.Entry) error {
	err := s.EntryRepository.Update(entry)
	if err != nil {
		return err
	}
	s.Bus.NotifyUpdated(entry)
	return nil
}

func (s *EntryService) Delete(entry model.Entry) error {
	err := s.EntryRepository.Delete(entry)
	if err != nil {
		return err
	}
	s.Bus.NotifyDeleted(entry)
	return nil
}

func (s *EntryService) FindById(id string) (model.Entry, error) {
	return s.EntryRepository.FindById(id)
}

func (s *EntryService) FindByUrl(url string) (model.Entry, error) {
	cfg, _ := s.siteConfigServcie.GetSiteConfig()
	if !strings.HasPrefix(url, cfg.FullUrl) {
		return nil, errors.New("url does not belong to blog")
	}
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	parts := strings.Split(url, "/")
	id := parts[len(parts)-1]
	return s.FindById(id)
}

func (s *EntryService) filterEntries(entries []model.Entry, published bool, drafts bool) []model.Entry {
	filteredEntries := make([]model.Entry, 0)
	for _, entry := range entries {
		if published && entry.PublishedAt() != nil && !entry.PublishedAt().IsZero() {
			filteredEntries = append(filteredEntries, entry)
		}
		if drafts && (entry.PublishedAt() == nil || entry.PublishedAt().IsZero()) {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	return filteredEntries
}

func (s *EntryService) FindAllByType(types *[]string, published bool, drafts bool) ([]model.Entry, error) {
	entries, err := s.EntryRepository.FindAll(types)
	return s.filterEntries(entries, published, drafts), err

}

func (s *EntryService) FindAll() ([]model.Entry, error) {
	entries, err := s.EntryRepository.FindAll(nil)
	return s.filterEntries(entries, true, true), err
}

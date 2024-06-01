package app

import (
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
	// only notify if the publishing date is set
	// otherwise this is a draft.
	// listeners might publish the entry to other services/platforms
	// this should only happen for publshed content
	if entry.PublishedAt() != nil && !entry.PublishedAt().IsZero() {
		s.Bus.NotifyCreated(entry)
	}
	return nil
}

func (s *EntryService) Update(entry model.Entry) error {
	err := s.EntryRepository.Update(entry)
	if err != nil {
		return err
	}
	// only notify if the publishing date is set
	// otherwise this is a draft.
	// listeners might publish the entry to other services/platforms
	// this should only happen for publshed content
	if entry.PublishedAt() != nil && !entry.PublishedAt().IsZero() {
		s.Bus.NotifyUpdated(entry)
	}
	return nil
}

func (s *EntryService) Delete(entry model.Entry) error {
	err := s.EntryRepository.Delete(entry)
	if err != nil {
		return err
	}
	// deletes should always be notfied
	// a published entry might be converted to a draft before deletion
	// omitting the deletion in this case would prevent deletion on other platforms
	s.Bus.NotifyDeleted(entry)
	return nil
}

func (s *EntryService) FindById(id string) (model.Entry, error) {
	entry, err := s.EntryRepository.FindById(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ErrEntryNotFound
		}
	}
	return entry, nil
}

func (s *EntryService) FindByUrl(url string) (model.Entry, error) {
	cfg, _ := s.siteConfigServcie.GetSiteConfig()
	if !strings.HasPrefix(url, cfg.FullUrl) {
		return nil, ErrEntryNotFound
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

func (s *EntryService) FindAllByTag(tag string, published bool, drafts bool) ([]model.Entry, error) {
	entries, err := s.EntryRepository.FindAllByTag(tag)
	return s.filterEntries(entries, published, drafts), err
}

func (s *EntryService) FindAll() ([]model.Entry, error) {
	entries, err := s.EntryRepository.FindAll(nil)
	return s.filterEntries(entries, true, true), err
}

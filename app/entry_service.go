package app

import (
	"fmt"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"regexp"
	"strings"
)

type EntryService struct {
	EntryRepository repository.EntryRepository
}

func NewEntryService(entryRepository repository.EntryRepository) *EntryService {
	return &EntryService{EntryRepository: entryRepository}
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

	return s.EntryRepository.Create(entry)
}

func (s *EntryService) Update(entry model.Entry) error {
	return s.EntryRepository.Update(entry)
}

func (s *EntryService) Delete(entry model.Entry) error {
	return s.EntryRepository.Delete(entry)
}

func (s *EntryService) FindById(id string) (model.Entry, error) {
	return s.EntryRepository.FindById(id)
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

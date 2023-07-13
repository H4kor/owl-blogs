package app

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
)

type EntryService struct {
	EntryRepository repository.EntryRepository
}

func NewEntryService(entryRepository repository.EntryRepository) *EntryService {
	return &EntryService{EntryRepository: entryRepository}
}

func (s *EntryService) Create(entry model.Entry) error {
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

func (s *EntryService) FindAllByType(types *[]string) ([]model.Entry, error) {
	return s.EntryRepository.FindAll(types)
}

func (s *EntryService) FindAll() ([]model.Entry, error) {
	entries, err := s.EntryRepository.FindAll(nil)
	if err != nil {
		return nil, err
	}
	// filter unpublished entries
	publishedEntries := make([]model.Entry, 0)
	for _, entry := range entries {
		if entry.PublishedAt() != nil && !entry.PublishedAt().IsZero() {
			publishedEntries = append(publishedEntries, entry)
		}
	}
	return publishedEntries, nil
}

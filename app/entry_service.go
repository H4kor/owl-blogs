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
	return s.EntryRepository.FindAll(nil)
}

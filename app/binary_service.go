package app

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
)

type BinaryService struct {
	repo repository.BinaryRepository
}

func NewBinaryFileService(repo repository.BinaryRepository) *BinaryService {
	return &BinaryService{repo: repo}
}

func (s *BinaryService) Create(name string, file []byte) (*model.BinaryFile, error) {
	return s.repo.Create(name, file, nil)
}

func (s *BinaryService) CreateEntryFile(name string, file []byte, entry model.Entry) (*model.BinaryFile, error) {
	return s.repo.Create(name, file, entry)
}

func (s *BinaryService) FindById(id string) (*model.BinaryFile, error) {
	return s.repo.FindById(id)
}

// ListIds list all ids of binary files
// if filter is not empty, the list will be filter to all ids which include the filter filter substring
// ids and filters are compared in lower case
func (s *BinaryService) ListIds(filter string) ([]string, error) {
	return s.repo.ListIds(filter)
}

func (s *BinaryService) Delete(binary *model.BinaryFile) error {
	return s.repo.Delete(binary)
}

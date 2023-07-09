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

func (s *BinaryService) FindById(id string) (*model.BinaryFile, error) {
	return s.repo.FindById(id)
}

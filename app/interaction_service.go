package app

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
)

type InteractionService struct {
	repo repository.InteractionRepository
}

func (s *InteractionService) ListInteractions() ([]model.Interaction, error) {
	return s.repo.ListAllInteractions()
}

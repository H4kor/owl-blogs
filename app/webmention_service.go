package app

import "owl-blogs/app/repository"

type WebmentionService struct {
	InteractionRepository repository.InteractionRepository
	EntryRepository       repository.EntryRepository
}

func NewWebmentionService(
	interactionRepository repository.InteractionRepository,
	entryRepository repository.EntryRepository,
) *WebmentionService {
	return &WebmentionService{
		InteractionRepository: interactionRepository,
		EntryRepository:       entryRepository,
	}
}

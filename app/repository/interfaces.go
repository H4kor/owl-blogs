package repository

import (
	"owl-blogs/domain/model"
)

type EntryRepository interface {
	Create(entry model.Entry) error
	Update(entry model.Entry) error
	Delete(entry model.Entry) error
	FindById(id string) (model.Entry, error)
	FindAll(types *[]string) ([]model.Entry, error)
}

type BinaryRepository interface {
	// Create creates a new binary file
	// The name is the original file name, and is not unique
	// BinaryFile.Id is a unique identifier
	Create(name string, data []byte, entry model.Entry) (*model.BinaryFile, error)
	FindById(id string) (*model.BinaryFile, error)
	FindByNameForEntry(name string, entry model.Entry) (*model.BinaryFile, error)
	ListIds() ([]string, error)
	Delete(binary *model.BinaryFile) error
}

type AuthorRepository interface {
	// Create creates a new author
	// It returns an error if the name is already taken
	Create(name string, passwordHash string) (*model.Author, error)

	Update(author *model.Author) error
	// FindByName finds an author by name
	// It returns an error if the author is not found
	FindByName(name string) (*model.Author, error)
}

type ConfigRepository interface {
	Get(name string, config interface{}) error
	Update(name string, siteConfig interface{}) error
}

type InteractionRepository interface {
	Create(interaction model.Interaction) error
	Update(interaction model.Interaction) error
	Delete(interaction model.Interaction) error
	FindById(id string) (model.Interaction, error)
	FindAll(entryId string) ([]model.Interaction, error)
}

package repository

import "owl-blogs/domain/model"

type EntryRepository interface {
	RegisterEntryType(entry model.Entry) error
	Create(entry model.Entry) error
	Update(entry model.Entry) error
	Delete(entry model.Entry) error
	FindById(id string) (model.Entry, error)
	FindAll(types *[]string) ([]model.Entry, error)
}

package repository

import (
	"owl-blogs/domain/model"
	"time"
)

type EntryRepository interface {
	Create(entry model.Entry, publishedAt *time.Time, metaData model.EntryMetaData) error
	Update(entry model.Entry) error
	Delete(entry model.Entry) error
	FindById(id string) (model.Entry, error)
	FindAll(types *[]string) ([]model.Entry, error)
}

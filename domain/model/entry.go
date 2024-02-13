package model

import (
	"time"
)

type EntryContent string

type Entry interface {
	ID() string
	Content() EntryContent
	PublishedAt() *time.Time
	AuthorId() string
	MetaData() EntryMetaData

	// Optional: can return empty string
	Title() string

	SetID(id string)
	SetPublishedAt(publishedAt *time.Time)
	SetMetaData(metaData EntryMetaData)
	SetAuthorId(authorId string)
}

type EntryMetaData interface {
	Form(binSvc BinaryStorageInterface) string
	ParseFormData(data HttpFormData, binSvc BinaryStorageInterface) (EntryMetaData, error)
}

type EntryBase struct {
	id          string
	publishedAt *time.Time
	authorId    string
}

func (e *EntryBase) ID() string {
	return e.id
}

func (e *EntryBase) PublishedAt() *time.Time {
	return e.publishedAt
}

func (e *EntryBase) SetID(id string) {
	e.id = id
}

func (e *EntryBase) SetPublishedAt(publishedAt *time.Time) {
	e.publishedAt = publishedAt
}

func (e *EntryBase) AuthorId() string {
	return e.authorId
}

func (e *EntryBase) SetAuthorId(authorId string) {
	e.authorId = authorId
}

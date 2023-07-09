package model

import "time"

type EntryContent string

type Entry interface {
	ID() string
	Content() EntryContent
	PublishedAt() *time.Time
	MetaData() interface{}

	// Optional: can return empty string
	Title() string

	SetID(id string)
	SetPublishedAt(publishedAt *time.Time)
	SetMetaData(metaData interface{})
}

type EntryMetaData interface {
}

type EntryBase struct {
	id          string
	publishedAt *time.Time
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

package model

import "time"

type EntryContent string

type Entry interface {
	ID() string
	Content() EntryContent
	PublishedAt() *time.Time
	MetaData() interface{}
	Create(id string, content string, publishedAt *time.Time, metaData EntryMetaData) error
}

type EntryMetaData interface {
}

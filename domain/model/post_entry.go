package model

import "time"

type PostEntry struct {
	id          string
	publishedAt *time.Time
	meta        PostEntryMetaData
}

type PostEntryMetaData struct {
	Title   string `owl:"inputType=text"`
	Content string `owl:"inputType=text widget=textarea"`
}

func (e *PostEntry) ID() string {
	return e.id
}

func (e *PostEntry) Content() EntryContent {
	return EntryContent(e.meta.Content)
}

func (e *PostEntry) PublishedAt() *time.Time {
	return e.publishedAt
}

func (e *PostEntry) MetaData() interface{} {
	return &e.meta
}

func (e *PostEntry) Create(id string, publishedAt *time.Time, metaData EntryMetaData) error {
	e.id = id
	e.publishedAt = publishedAt
	e.meta = *metaData.(*PostEntryMetaData)
	return nil
}

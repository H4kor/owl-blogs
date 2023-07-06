package test

import (
	"owl-blogs/domain/model"
	"time"
)

type MockEntryMetaData struct {
	Str    string
	Number int
	Date   time.Time
}

type MockEntry struct {
	id          string
	publishedAt *time.Time
	metaData    *MockEntryMetaData
}

func (e *MockEntry) ID() string {
	return e.id
}

func (e *MockEntry) Content() model.EntryContent {
	return model.EntryContent(e.metaData.Str)
}

func (e *MockEntry) PublishedAt() *time.Time {
	return e.publishedAt
}

func (e *MockEntry) MetaData() interface{} {
	return e.metaData
}

func (e *MockEntry) Create(id string, publishedAt *time.Time, metaData model.EntryMetaData) error {
	e.id = id
	e.publishedAt = publishedAt
	e.metaData = metaData.(*MockEntryMetaData)
	return nil
}

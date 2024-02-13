package test

import (
	"owl-blogs/domain/model"
	"owl-blogs/web/forms"
	"time"
)

type MockEntryMetaData struct {
	forms.DefaultForm

	Str    string
	Number int
	Date   time.Time
	Title  string
}

type MockEntry struct {
	model.EntryBase
	metaData *MockEntryMetaData
}

func (e *MockEntry) Content() model.EntryContent {
	return model.EntryContent(e.metaData.Str)
}

func (e *MockEntry) MetaData() model.EntryMetaData {
	return e.metaData
}

func (e *MockEntry) SetMetaData(metaData model.EntryMetaData) {
	e.metaData = metaData.(*MockEntryMetaData)
}

func (e *MockEntry) Title() string {
	return e.metaData.Title
}

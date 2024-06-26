package test

import (
	"html/template"
	"owl-blogs/domain/model"
	"time"
)

type MockEntryMetaData struct {
	Str    string
	Number int
	Date   time.Time
	Title  string
}

// Form implements model.EntryMetaData.
func (*MockEntryMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	panic("unimplemented")
}

// ParseFormData implements model.EntryMetaData.
func (*MockEntryMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	panic("unimplemented")
}

type MockEntry struct {
	model.EntryBase
	metaData *MockEntryMetaData
}

func (e *MockEntry) Content() template.HTML {
	return template.HTML(e.metaData.Str)
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

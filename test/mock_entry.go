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
	model.EntryBase
	metaData *MockEntryMetaData
}

func (e *MockEntry) Content() model.EntryContent {
	return model.EntryContent(e.metaData.Str)
}

func (e *MockEntry) MetaData() interface{} {
	return e.metaData
}

func (e *MockEntry) SetMetaData(metaData interface{}) {
	e.metaData = metaData.(*MockEntryMetaData)
}

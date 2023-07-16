package entrytypes

import "owl-blogs/domain/model"

type Note struct {
	model.EntryBase
	meta NoteMetaData
}

type NoteMetaData struct {
	Content string `owl:"inputType=text widget=textarea"`
}

func (e *Note) Title() string {
	return ""
}

func (e *Note) Content() model.EntryContent {
	return model.EntryContent(e.meta.Content)
}

func (e *Note) MetaData() interface{} {
	return &e.meta
}

func (e *Note) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*NoteMetaData)
}

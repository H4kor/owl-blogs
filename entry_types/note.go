package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

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
	str, err := render.RenderTemplateToString("entry/Note", e)
	if err != nil {
		fmt.Println(err)
	}
	return model.EntryContent(str)
}

func (e *Note) MetaData() interface{} {
	return &e.meta
}

func (e *Note) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*NoteMetaData)
}

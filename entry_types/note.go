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

// Form implements model.EntryMetaData.
func (meta *NoteMetaData) Form(binSvc model.BinaryStorageInterface) string {
	f, _ := render.RenderTemplateToString("forms/Note", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (*NoteMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) (model.EntryMetaData, error) {
	return &NoteMetaData{
		Content: data.FormValue("content"),
	}, nil
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

func (e *Note) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Note) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*NoteMetaData)
}

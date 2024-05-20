package entrytypes

import (
	"fmt"
	"html/template"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Note struct {
	model.EntryBase
	meta NoteMetaData
}

type NoteMetaData struct {
	Content string
}

// Form implements model.EntryMetaData.
func (meta *NoteMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	f, _ := render.RenderTemplateToString("forms/Note", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (meta *NoteMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	meta.Content = data.FormValue("content")
	return nil
}

func (e *Note) Title() string {
	return ""
}

func (e *Note) Content() template.HTML {
	str, err := render.RenderTemplateToString("entry/Note", e)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (e *Note) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Note) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*NoteMetaData)
}

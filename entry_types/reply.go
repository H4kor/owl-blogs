package entrytypes

import (
	"fmt"
	"html/template"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Reply struct {
	model.EntryBase
	meta ReplyMetaData
}

type ReplyMetaData struct {
	Title   string
	Url     string
	Content string
}

// Form implements model.EntryMetaData.
func (meta *ReplyMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	f, _ := render.RenderTemplateToString("forms/Reply", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (meta *ReplyMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	meta.Title = data.FormValue("title")
	meta.Url = data.FormValue("url")
	meta.Content = data.FormValue("content")
	return nil
}

func (e *Reply) Title() string {
	return "Re: " + e.meta.Title
}

func (e *Reply) Content() template.HTML {
	str, err := render.RenderTemplateToString("entry/Reply", e)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (e *Reply) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Reply) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*ReplyMetaData)
}

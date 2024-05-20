package entrytypes

import (
	"fmt"
	"html/template"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Page struct {
	model.EntryBase
	meta PageMetaData
}

type PageMetaData struct {
	Title   string
	Content string
}

// Form implements model.EntryMetaData.
func (meta *PageMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	f, _ := render.RenderTemplateToString("forms/Page", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (meta *PageMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	meta.Title = data.FormValue("title")
	meta.Content = data.FormValue("content")
	return nil
}

func (e *Page) Title() string {
	return e.meta.Title
}

func (e *Page) Content() template.HTML {
	str, err := render.RenderTemplateToString("entry/Page", e)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (e *Page) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Page) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*PageMetaData)
}

package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Page struct {
	model.EntryBase
	meta PageMetaData
}

type PageMetaData struct {
	Title   string `owl:"inputType=text"`
	Content string `owl:"inputType=text widget=textarea"`
}

// Form implements model.EntryMetaData.
func (meta *PageMetaData) Form(binSvc model.BinaryStorageInterface) string {
	f, _ := render.RenderTemplateToString("forms/Page", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (*PageMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) (model.EntryMetaData, error) {
	return &PageMetaData{
		Title:   data.FormValue("title"),
		Content: data.FormValue("content"),
	}, nil
}

func (e *Page) Title() string {
	return e.meta.Title
}

func (e *Page) Content() model.EntryContent {
	str, err := render.RenderTemplateToString("entry/Page", e)
	if err != nil {
		fmt.Println(err)
	}
	return model.EntryContent(str)
}

func (e *Page) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Page) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*PageMetaData)
}

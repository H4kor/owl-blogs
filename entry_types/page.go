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

func (e *Page) MetaData() interface{} {
	return &e.meta
}

func (e *Page) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*PageMetaData)
}

package entrytypes

import (
	"fmt"
	"html/template"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Bookmark struct {
	model.EntryBase
	meta BookmarkMetaData
}

type BookmarkMetaData struct {
	Title   string
	Url     string
	Content string
}

// Form implements model.EntryMetaData.
func (meta *BookmarkMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	f, _ := render.RenderTemplateToString("forms/Bookmark", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (meta *BookmarkMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	meta.Title = data.FormValue("title")
	meta.Url = data.FormValue("url")
	meta.Content = data.FormValue("content")
	return nil
}

func (e *Bookmark) Title() string {
	return e.meta.Title
}

func (e *Bookmark) Content() template.HTML {
	str, err := render.RenderTemplateToString("entry/Bookmark", e)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (e *Bookmark) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Bookmark) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*BookmarkMetaData)
}

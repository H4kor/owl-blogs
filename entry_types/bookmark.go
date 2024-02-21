package entrytypes

import (
	"fmt"
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
func (meta *BookmarkMetaData) Form(binSvc model.BinaryStorageInterface) string {
	f, _ := render.RenderTemplateToString("forms/Bookmark", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (*BookmarkMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) (model.EntryMetaData, error) {
	return &BookmarkMetaData{
		Title:   data.FormValue("title"),
		Url:     data.FormValue("url"),
		Content: data.FormValue("content"),
	}, nil
}

func (e *Bookmark) Title() string {
	return e.meta.Title
}

func (e *Bookmark) Content() model.EntryContent {
	str, err := render.RenderTemplateToString("entry/Bookmark", e)
	if err != nil {
		fmt.Println(err)
	}
	return model.EntryContent(str)
}

func (e *Bookmark) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Bookmark) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*BookmarkMetaData)
}

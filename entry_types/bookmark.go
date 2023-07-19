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
	Title   string `owl:"inputType=text"`
	Url     string `owl:"inputType=text"`
	Content string `owl:"inputType=text widget=textarea"`
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

func (e *Bookmark) MetaData() interface{} {
	return &e.meta
}

func (e *Bookmark) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*BookmarkMetaData)
}

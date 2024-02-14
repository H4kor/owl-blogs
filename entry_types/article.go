package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Article struct {
	model.EntryBase
	meta ArticleMetaData
}

type ArticleMetaData struct {
	Title   string `owl:"inputType=text"`
	Content string `owl:"inputType=text widget=textarea"`
}

// Form implements model.EntryMetaData.
func (meta *ArticleMetaData) Form(binSvc model.BinaryStorageInterface) string {
	f, _ := render.RenderTemplateToString("forms/Article", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (*ArticleMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) (model.EntryMetaData, error) {
	return &ArticleMetaData{
		Title:   data.FormValue("title"),
		Content: data.FormValue("content"),
	}, nil
}

func (e *Article) Title() string {
	return e.meta.Title
}

func (e *Article) Content() model.EntryContent {
	str, err := render.RenderTemplateToString("entry/Article", e)
	if err != nil {
		fmt.Println(err)
	}
	return model.EntryContent(str)
}

func (e *Article) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Article) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*ArticleMetaData)
}

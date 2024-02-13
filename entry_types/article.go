package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"owl-blogs/web/forms"
)

type Article struct {
	model.EntryBase
	meta ArticleMetaData
}

type ArticleMetaData struct {
	forms.DefaultForm
	Title   string `owl:"inputType=text"`
	Content string `owl:"inputType=text widget=textarea"`
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

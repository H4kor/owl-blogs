package entrytypes

import (
	"fmt"
	"html/template"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"

	vocab "github.com/go-ap/activitypub"
)

type Article struct {
	model.EntryBase
	meta ArticleMetaData
}

type ArticleMetaData struct {
	Title   string
	Content string
}

// Form implements model.EntryMetaData.
func (meta *ArticleMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	f, _ := render.RenderTemplateToString("forms/Article", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (meta *ArticleMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	meta.Title = data.FormValue("title")
	meta.Content = data.FormValue("content")
	return nil
}

func (e *Article) Title() string {
	return e.meta.Title
}

func (e *Article) Content() template.HTML {
	str, err := render.RenderTemplateToString("entry/Article", e)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (e *Article) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Article) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*ArticleMetaData)
}

func (e *Article) ActivityObject(siteCfg model.SiteConfig, binSvc app.BinaryService) vocab.Object {
	content := e.Content()

	obj := vocab.Article{
		Type:      "Article",
		Published: *e.PublishedAt(),
		Name: vocab.NaturalLanguageValues{
			{Value: vocab.Content(e.Title())},
		},
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(string(content))},
		},
	}
	return obj

}

func (e *Article) Tags() []string {
	return extractTags(e.meta.Content)
}

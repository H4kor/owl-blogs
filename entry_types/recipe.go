package entrytypes

import (
	"fmt"
	"html/template"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strings"

	vocab "github.com/go-ap/activitypub"
)

type Recipe struct {
	model.EntryBase
	meta RecipeMetaData
}

type RecipeMetaData struct {
	Title       string
	Yield       string
	Duration    string
	Ingredients []string
	Content     string
}

// Form implements model.EntryMetaData.
func (meta *RecipeMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	f, _ := render.RenderTemplateToString("forms/Recipe", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (meta *RecipeMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	ings := strings.Split(data.FormValue("ingredients"), "\n")
	clean := make([]string, 0)
	for _, ing := range ings {
		if strings.TrimSpace(ing) != "" {
			clean = append(clean, strings.TrimSpace(ing))
		}
	}
	meta.Title = data.FormValue("title")
	meta.Yield = data.FormValue("yield")
	meta.Duration = data.FormValue("duration")
	meta.Ingredients = clean
	meta.Content = data.FormValue("content")
	return nil
}

func (e *Recipe) Title() string {
	return e.meta.Title
}

func (e *Recipe) Content() template.HTML {
	str, err := render.RenderTemplateToString("entry/Recipe", e)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (e *Recipe) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Recipe) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*RecipeMetaData)
}

func (e *Recipe) ActivityObject(siteCfg model.SiteConfig, binSvc app.BinaryService) vocab.Object {
	content := e.Content()

	image := vocab.Article{
		Type:      "Article",
		Published: *e.PublishedAt(),
		Name: vocab.NaturalLanguageValues{
			{Value: vocab.Content(e.Title())},
		},
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(string(content))},
		},
	}
	return image

}

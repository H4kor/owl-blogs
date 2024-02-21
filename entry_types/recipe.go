package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strings"
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
func (meta *RecipeMetaData) Form(binSvc model.BinaryStorageInterface) string {
	f, _ := render.RenderTemplateToString("forms/Recipe", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (*RecipeMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) (model.EntryMetaData, error) {
	ings := strings.Split(data.FormValue("ingredients"), "\n")
	clean := make([]string, 0)
	for _, ing := range ings {
		if strings.TrimSpace(ing) != "" {
			clean = append(clean, strings.TrimSpace(ing))
		}
	}
	return &RecipeMetaData{
		Title:       data.FormValue("title"),
		Yield:       data.FormValue("yield"),
		Duration:    data.FormValue("duration"),
		Ingredients: clean,
		Content:     data.FormValue("content"),
	}, nil
}

func (e *Recipe) Title() string {
	return e.meta.Title
}

func (e *Recipe) Content() model.EntryContent {
	str, err := render.RenderTemplateToString("entry/Recipe", e)
	if err != nil {
		fmt.Println(err)
	}
	return model.EntryContent(str)
}

func (e *Recipe) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Recipe) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*RecipeMetaData)
}

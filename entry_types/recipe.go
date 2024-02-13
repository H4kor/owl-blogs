package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"owl-blogs/web/forms"
)

type Recipe struct {
	model.EntryBase
	meta RecipeMetaData
}

type RecipeMetaData struct {
	forms.DefaultForm

	Title       string   `owl:"inputType=text"`
	Yield       string   `owl:"inputType=text"`
	Duration    string   `owl:"inputType=text"`
	Ingredients []string `owl:"inputType=text widget=textlist"`
	Content     string   `owl:"inputType=text widget=textarea"`
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

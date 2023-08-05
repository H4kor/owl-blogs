package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Recipe struct {
	model.EntryBase
	meta RecipeMetaData
}

type RecipeMetaData struct {
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

func (e *Recipe) MetaData() interface{} {
	return &e.meta
}

func (e *Recipe) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*RecipeMetaData)
}

package model

import (
	"fmt"
	"owl-blogs/render"
)

type Recipe struct {
	EntryBase
	meta RecipeMetaData
}

type RecipeMetaData struct {
	Title       string   `owl:"inputType=text"`
	Yield       string   `owl:"inputType=text"`
	Duration    string   `owl:"inputType=text"`
	Ingredients []string `owl:"inputType=text widget=textarea"`
	Content     string   `owl:"inputType=text widget=textarea"`
}

func (e *Recipe) Title() string {
	return e.meta.Title
}

func (e *Recipe) Content() EntryContent {
	str, err := render.RenderTemplateToString("entry/Recipe", e)
	if err != nil {
		fmt.Println(err)
	}
	return EntryContent(str)
}

func (e *Recipe) MetaData() interface{} {
	return &e.meta
}

func (e *Recipe) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*RecipeMetaData)
}

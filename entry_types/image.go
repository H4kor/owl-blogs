package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"owl-blogs/web/forms"
)

type Image struct {
	model.EntryBase
	meta ImageMetaData
}

type ImageMetaData struct {
	forms.DefaultForm

	ImageId string `owl:"inputType=file"`
	Title   string `owl:"inputType=text"`
	Content string `owl:"inputType=text widget=textarea"`
}

func (e *Image) Title() string {
	return e.meta.Title
}

func (e *Image) Content() model.EntryContent {
	str, err := render.RenderTemplateToString("entry/Image", e)
	if err != nil {
		fmt.Println(err)
	}
	return model.EntryContent(str)
}

func (e *Image) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Image) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*ImageMetaData)
}

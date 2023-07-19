package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Image struct {
	model.EntryBase
	meta ImageMetaData
}

type ImageMetaData struct {
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

func (e *Image) MetaData() interface{} {
	return &e.meta
}

func (e *Image) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*ImageMetaData)
}

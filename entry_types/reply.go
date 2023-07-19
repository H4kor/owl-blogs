package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Reply struct {
	model.EntryBase
	meta ReplyMetaData
}

type ReplyMetaData struct {
	Title   string `owl:"inputType=text"`
	Url     string `owl:"inputType=text"`
	Content string `owl:"inputType=text widget=textarea"`
}

func (e *Reply) Title() string {
	return "Re: " + e.meta.Title
}

func (e *Reply) Content() model.EntryContent {
	str, err := render.RenderTemplateToString("entry/Reply", e)
	if err != nil {
		fmt.Println(err)
	}
	return model.EntryContent(str)
}

func (e *Reply) MetaData() interface{} {
	return &e.meta
}

func (e *Reply) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*ReplyMetaData)
}

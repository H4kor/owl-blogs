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

// Form implements model.EntryMetaData.
func (meta *ReplyMetaData) Form(binSvc model.BinaryStorageInterface) string {
	f, _ := render.RenderTemplateToString("forms/Reply", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (*ReplyMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) (model.EntryMetaData, error) {
	return &ReplyMetaData{
		Title:   data.FormValue("title"),
		Url:     data.FormValue("url"),
		Content: data.FormValue("content"),
	}, nil
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

func (e *Reply) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Reply) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*ReplyMetaData)
}

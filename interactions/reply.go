package interactions

import (
	"fmt"
	"html/template"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Reply struct {
	model.InteractionBase
	meta ReplyMetaData
}

type ReplyMetaData struct {
	SenderUrl   string
	SenderName  string
	OriginalUrl string
	Content     template.HTML
}

func (i *Reply) Content() template.HTML {
	str, err := render.RenderTemplateToString("interaction/Reply", i)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (i *Reply) MetaData() interface{} {
	return &i.meta
}

func (i *Reply) SetMetaData(metaData interface{}) {
	i.meta = *metaData.(*ReplyMetaData)
}

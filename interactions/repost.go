package interactions

import (
	"fmt"
	"html/template"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Repost struct {
	model.InteractionBase
	meta RepostMetaData
}

type RepostMetaData struct {
	SenderUrl  string
	SenderName string
}

func (i *Repost) Content() template.HTML {
	str, err := render.RenderTemplateToString("interaction/Repost", i)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (i *Repost) MetaData() interface{} {
	return &i.meta
}

func (i *Repost) SetMetaData(metaData interface{}) {
	i.meta = *metaData.(*RepostMetaData)
}

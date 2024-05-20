package interactions

import (
	"fmt"
	"html/template"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Like struct {
	model.InteractionBase
	meta LikeMetaData
}

type LikeMetaData struct {
	SenderUrl  string
	SenderName string
}

func (i *Like) Content() template.HTML {
	str, err := render.RenderTemplateToString("interaction/Like", i)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (i *Like) MetaData() interface{} {
	return &i.meta
}

func (i *Like) SetMetaData(metaData interface{}) {
	i.meta = *metaData.(*LikeMetaData)
}

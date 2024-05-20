package interactions

import (
	"fmt"
	"html/template"
	"owl-blogs/domain/model"
	"owl-blogs/render"
)

type Webmention struct {
	model.InteractionBase
	meta WebmentionMetaData
}

type WebmentionMetaData struct {
	Source string
	Target string
	Title  string
}

func (i *Webmention) Content() template.HTML {
	str, err := render.RenderTemplateToString("interaction/Webmention", i)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (i *Webmention) MetaData() interface{} {
	return &i.meta
}

func (i *Webmention) SetMetaData(metaData interface{}) {
	i.meta = *metaData.(*WebmentionMetaData)
}

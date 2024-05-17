package interactions

import (
	"fmt"
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

func (i *Like) Content() model.InteractionContent {
	str, err := render.RenderTemplateToString("interaction/Like", i)
	if err != nil {
		fmt.Println(err)
	}
	return model.InteractionContent(str)
}

func (i *Like) MetaData() interface{} {
	return &i.meta
}

func (i *Like) SetMetaData(metaData interface{}) {
	i.meta = *metaData.(*LikeMetaData)
}

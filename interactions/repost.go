package interactions

import (
	"fmt"
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

func (i *Repost) Content() model.InteractionContent {
	str, err := render.RenderTemplateToString("interaction/Repost", i)
	if err != nil {
		fmt.Println(err)
	}
	return model.InteractionContent(str)
}

func (i *Repost) MetaData() interface{} {
	return &i.meta
}

func (i *Repost) SetMetaData(metaData interface{}) {
	i.meta = *metaData.(*RepostMetaData)
}

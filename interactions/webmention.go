package interactions

import "owl-blogs/domain/model"

type Webmention struct {
	model.InteractionBase
	meta WebmentionInteractionMetaData
}

type WebmentionInteractionMetaData struct {
	Source string
	Target string
	Title  string
}

func (i *Webmention) Content() model.InteractionContent {
	return model.InteractionContent(i.meta.Source)
}

func (i *Webmention) MetaData() interface{} {
	return &i.meta
}

func (i *Webmention) SetMetaData(metaData interface{}) {
	i.meta = *metaData.(*WebmentionInteractionMetaData)
}

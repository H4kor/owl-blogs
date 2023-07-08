package model

type Image struct {
	EntryBase
	meta ImageMetaData
}

type ImageMetaData struct {
	ImageId string `owl:"inputType=file"`
	Content string `owl:"inputType=text widget=textarea"`
}

func (e *Image) Content() EntryContent {
	return EntryContent(e.meta.Content)
}

func (e *Image) MetaData() interface{} {
	return &e.meta
}

func (e *Image) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*ImageMetaData)
}

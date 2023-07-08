package model

type ImageEntry struct {
	EntryBase
	meta ImageEntryMetaData
}

type ImageEntryMetaData struct {
	ImageId string `owl:"inputType=file"`
	Content string `owl:"inputType=text widget=textarea"`
}

func (e *ImageEntry) Content() EntryContent {
	return EntryContent(e.meta.Content)
}

func (e *ImageEntry) MetaData() interface{} {
	return &e.meta
}

func (e *ImageEntry) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*ImageEntryMetaData)
}

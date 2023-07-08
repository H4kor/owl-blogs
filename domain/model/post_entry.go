package model

type PostEntry struct {
	EntryBase
	meta PostEntryMetaData
}

type PostEntryMetaData struct {
	Title   string `owl:"inputType=text"`
	Content string `owl:"inputType=text widget=textarea"`
}

func (e *PostEntry) Content() EntryContent {
	return EntryContent(e.meta.Content)
}

func (e *PostEntry) MetaData() interface{} {
	return &e.meta
}

func (e *PostEntry) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*PostEntryMetaData)
}

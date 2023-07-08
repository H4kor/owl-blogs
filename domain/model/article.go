package model

type Article struct {
	EntryBase
	meta ArticleMetaData
}

type ArticleMetaData struct {
	Title   string `owl:"inputType=text"`
	Content string `owl:"inputType=text widget=textarea"`
}

func (e *Article) Content() EntryContent {
	return EntryContent(e.meta.Content)
}

func (e *Article) MetaData() interface{} {
	return &e.meta
}

func (e *Article) SetMetaData(metaData interface{}) {
	e.meta = *metaData.(*ArticleMetaData)
}

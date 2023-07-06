package model

import "time"

type ImageEntry struct {
	id          string
	publishedAt *time.Time
	meta        ImageEntryMetaData
}

type ImageEntryMetaData struct {
	ImagePath string       `owl:"inputType=file"`
	Content   EntryContent `owl:"inputType=text widget=textarea"`
}

func (e *ImageEntry) ID() string {
	return e.id
}

func (e *ImageEntry) Content() EntryContent {
	return e.meta.Content
}

func (e *ImageEntry) PublishedAt() *time.Time {
	return e.publishedAt
}

func (e *ImageEntry) MetaData() interface{} {
	return &ImageEntryMetaData{
		ImagePath: e.meta.ImagePath,
		Content:   e.meta.Content,
	}
}

func (e *ImageEntry) Create(id string, publishedAt *time.Time, metaData EntryMetaData) error {
	e.id = id
	e.publishedAt = publishedAt
	e.meta = *metaData.(*ImageEntryMetaData)
	return nil
}

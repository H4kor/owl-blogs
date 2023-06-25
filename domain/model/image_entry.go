package model

import "time"

type ImageEntry struct {
	id          string
	content     EntryContent
	publishedAt *time.Time
	ImagePath   string
}

type ImageEntryMetaData struct {
	ImagePath string `owl:"type=upload"`
}

func (e *ImageEntry) ID() string {
	return e.id
}

func (e *ImageEntry) Content() EntryContent {
	return e.content
}

func (e *ImageEntry) PublishedAt() *time.Time {
	return e.publishedAt
}

func (e *ImageEntry) MetaData() interface{} {
	return &ImageEntryMetaData{
		ImagePath: e.ImagePath,
	}
}

func (e *ImageEntry) Create(id string, content string, publishedAt *time.Time, metaData EntryMetaData) error {
	e.id = id
	e.content = EntryContent(content)
	e.publishedAt = publishedAt
	e.ImagePath = metaData.(*ImageEntryMetaData).ImagePath
	return nil
}

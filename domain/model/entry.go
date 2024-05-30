package model

import (
	"html/template"
	"net/url"
	"time"
)

type Entry interface {
	ID() string
	Content() template.HTML
	PublishedAt() *time.Time
	AuthorId() string
	MetaData() EntryMetaData

	// Optional: can return empty string
	Title() string
	ImageUrl() string
	Tags() []string

	SetID(id string)
	SetPublishedAt(publishedAt *time.Time)
	SetMetaData(metaData EntryMetaData)
	SetAuthorId(authorId string)

	FullUrl(cfg SiteConfig) string
}

type EntryMetaData interface {
	Formable
}

type EntryBase struct {
	id          string
	publishedAt *time.Time
	authorId    string
}

func (e *EntryBase) ID() string {
	return e.id
}

func (e *EntryBase) PublishedAt() *time.Time {
	return e.publishedAt
}

func (e *EntryBase) ImageUrl() string {
	return ""
}

func (e *EntryBase) Tags() []string {
	return []string{}
}

func (e *EntryBase) SetID(id string) {
	e.id = id
}

func (e *EntryBase) SetPublishedAt(publishedAt *time.Time) {
	e.publishedAt = publishedAt
}

func (e *EntryBase) AuthorId() string {
	return e.authorId
}

func (e *EntryBase) SetAuthorId(authorId string) {
	e.authorId = authorId
}

func (e *EntryBase) FullUrl(cfg SiteConfig) string {
	u, _ := url.JoinPath(cfg.FullUrl, "/posts/", e.ID(), "/")
	return u
}

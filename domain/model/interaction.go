package model

import (
	"html/template"
	"time"
)

// Interaction is a generic interface for all interactions with entries
// These interactions can be:
// - Webmention, Pingback, Trackback
// - Likes, Comments on third party sites
// - Comments on the site itself
type Interaction interface {
	ID() string
	EntryID() string
	Content() template.HTML
	CreatedAt() time.Time
	MetaData() interface{}

	SetID(id string)
	SetEntryID(entryID string)
	SetCreatedAt(createdAt time.Time)
	SetMetaData(metaData interface{})
}

type InteractionBase struct {
	id        string
	entryID   string
	createdAt time.Time
}

func (i *InteractionBase) ID() string {
	return i.id
}

func (i *InteractionBase) EntryID() string {
	return i.entryID
}

func (i *InteractionBase) CreatedAt() time.Time {
	return i.createdAt
}

func (i *InteractionBase) SetID(id string) {
	i.id = id
}

func (i *InteractionBase) SetEntryID(entryID string) {
	i.entryID = entryID
}

func (i *InteractionBase) SetCreatedAt(createdAt time.Time) {
	i.createdAt = createdAt
}

package entrytypes

import (
	"fmt"
	"html/template"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"regexp"
	"strings"

	vocab "github.com/go-ap/activitypub"
)

type Note struct {
	model.EntryBase
	meta NoteMetaData
}

type NoteMetaData struct {
	Content string
}

// Form implements model.EntryMetaData.
func (meta *NoteMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	f, _ := render.RenderTemplateToString("forms/Note", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (meta *NoteMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	meta.Content = data.FormValue("content")
	return nil
}

func (e *Note) Title() string {
	return ""
}

func (e *Note) Content() template.HTML {
	str, err := render.RenderTemplateToString("entry/Note", e)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (e *Note) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Note) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*NoteMetaData)
}

func (e *Note) ActivityObject(siteCfg model.SiteConfig, binSvc app.BinaryService) vocab.Object {

	content := e.Content()
	tags := vocab.ItemCollection{}

	note := vocab.Note{
		ID:        vocab.ID(e.FullUrl(siteCfg)),
		Type:      "Note",
		Published: *e.PublishedAt(),
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(content)},
		},
		Tag: tags,
	}
	return note

}

func (e *Note) Tags() []string {
	// TODO: move into more generic structure
	// should also be usable elsewhere
	r := regexp.MustCompile("#[a-zA-Z0-9_]+")
	content := e.meta.Content
	matches := r.FindAllString(string(content), -1)
	tags := make([]string, 0)
	for _, hashtag := range matches {
		tag, _ := strings.CutPrefix(hashtag, "#")
		tags = append(tags, tag)
	}
	return tags
}

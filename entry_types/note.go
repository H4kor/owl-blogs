package entrytypes

import (
	"fmt"
	"html/template"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"

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
	// TODO: move into more generic structure
	// r := regexp.MustCompile("#[a-zA-Z0-9_]+")
	// matches := r.FindAllString(string(content), -1)
	// should also be usable elsewhere
	// for _, hashtag := range matches {
	// 	tags.Append(vocab.Object{
	// 		ID:   vocab.ID(svc.HashtagId(hashtag)),
	// 		Name: vocab.NaturalLanguageValues{{Value: vocab.Content(hashtag)}},
	// 	})
	// }

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

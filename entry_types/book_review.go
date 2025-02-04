package entrytypes

import (
	"errors"
	"fmt"
	"html/template"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"strconv"

	vocab "github.com/go-ap/activitypub"
)


type BookReview struct {
    model.EntryBase
    meta BookReviewMetaData
}

type BookReviewMetaData struct {
    Book string
    BookUrl string
    Author string
    AuthorUrl string
    // 1 to 5
    Rating int
    Content string
}


func (meta *BookReviewMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	f, _ := render.RenderTemplateToString("forms/BookReview", meta)
	return f
}

func (meta *BookReviewMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	meta.Book = data.FormValue("book")
	meta.BookUrl = data.FormValue("book_url")
	meta.Author = data.FormValue("author")
	meta.AuthorUrl = data.FormValue("author_url")
	meta.Content = data.FormValue("content")
    var err error
    meta.Rating, err = strconv.Atoi(data.FormValue("rating"))
    if meta.Rating < 1 || meta.Rating > 5 {
        err = errors.New("rating must be between 1 and 5")
    }
	return err
}


func (e *BookReview) Title() string {
    return "Review: " + e.meta.Book + " by " + e.meta.Author
}

func (e *BookReview) Content() template.HTML {
	str, err := render.RenderTemplateToString("entry/BookReview", e)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (e *BookReview) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *BookReview) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*BookReviewMetaData)
}

func (e *BookReview) ActivityObject(siteCfg model.SiteConfig, binSvc app.BinaryService) vocab.Object {
	content := e.Content()

	obj := vocab.Article{
		Type:      "BookReview",
		Published: *e.PublishedAt(),
		Name: vocab.NaturalLanguageValues{
			{Value: vocab.Content(e.Title())},
		},
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(string(content))},
		},
	}
	return obj

}

func (e *BookReview) Tags() []string {
	return extractTags(e.meta.Content)
}

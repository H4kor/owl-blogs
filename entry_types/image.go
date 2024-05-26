package entrytypes

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/url"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"

	vocab "github.com/go-ap/activitypub"
)

type Image struct {
	model.EntryBase
	meta ImageMetaData
}

type ImageMetaData struct {
	ImageId string
	Title   string
	Content string
}

// Form implements model.EntryMetaData.
func (meta *ImageMetaData) Form(binSvc model.BinaryStorageInterface) template.HTML {
	f, _ := render.RenderTemplateToString("forms/Image", meta)
	return f
}

// ParseFormData implements model.EntryMetaData.
func (meta *ImageMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) error {
	file, err := data.FormFile("image")
	var imgId = meta.ImageId
	if err != nil && imgId == "" {
		return err
	} else if err == nil {
		fileData, err := file.Open()
		if err != nil {
			return err
		}
		defer fileData.Close()

		fileBytes := make([]byte, file.Size)
		_, err = fileData.Read(fileBytes)
		if err != nil {
			return err
		}
		bin, err := binSvc.Create(file.Filename, fileBytes)
		if err != nil {
			return err
		}
		imgId = bin.Id
	}

	meta.ImageId = imgId
	meta.Title = data.FormValue("title")
	meta.Content = data.FormValue("content")
	return nil
}

func (e *Image) Title() string {
	return e.meta.Title
}

func (e *Image) ImageUrl() string {
	return "/media/" + e.meta.ImageId
}

func (e *Image) Content() template.HTML {
	str, err := render.RenderTemplateToString("entry/Image", e)
	if err != nil {
		fmt.Println(err)
	}
	return template.HTML(str)
}

func (e *Image) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Image) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*ImageMetaData)
}

func (e *Image) ActivityObject(siteCfg model.SiteConfig, binSvc app.BinaryService) vocab.Object {
	content := e.Content()

	imgPath := e.ImageUrl()
	fullImageUrl, _ := url.JoinPath(siteCfg.FullUrl, imgPath)
	binaryFile, err := binSvc.FindById(e.MetaData().(*ImageMetaData).ImageId)
	if err != nil {
		slog.Error("cannot get image file")
	}

	attachments := vocab.ItemCollection{}
	attachments = append(attachments, vocab.Document{
		Type:      vocab.DocumentType,
		MediaType: vocab.MimeType(binaryFile.Mime()),
		URL:       vocab.ID(fullImageUrl),
		Name: vocab.NaturalLanguageValues{
			{Value: vocab.Content(content)},
		},
	})

	image := vocab.Image{
		Type:      "Image",
		Published: *e.PublishedAt(),
		Name: vocab.NaturalLanguageValues{
			{Value: vocab.Content(e.Title())},
		},
		Content: vocab.NaturalLanguageValues{
			{Value: vocab.Content(e.Title() + "<br><br>" + string(content))},
		},
		Attachment: attachments,
		// Tag: tags,
	}
	return image

}

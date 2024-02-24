package entrytypes

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/render"
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
func (meta *ImageMetaData) Form(binSvc model.BinaryStorageInterface) string {
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

func (e *Image) Content() model.EntryContent {
	str, err := render.RenderTemplateToString("entry/Image", e)
	if err != nil {
		fmt.Println(err)
	}
	return model.EntryContent(str)
}

func (e *Image) MetaData() model.EntryMetaData {
	return &e.meta
}

func (e *Image) SetMetaData(metaData model.EntryMetaData) {
	e.meta = *metaData.(*ImageMetaData)
}

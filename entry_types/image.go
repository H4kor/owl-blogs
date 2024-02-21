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
func (meta *ImageMetaData) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) (model.EntryMetaData, error) {
	file, err := data.FormFile("image")
	var imgId = meta.ImageId
	if err != nil && imgId == "" {
		return nil, err
	} else if err == nil {
		fileData, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer fileData.Close()

		fileBytes := make([]byte, file.Size)
		_, err = fileData.Read(fileBytes)
		if err != nil {
			return nil, err
		}
		bin, err := binSvc.Create(file.Filename, fileBytes)
		if err != nil {
			return nil, err
		}
		imgId = bin.Id
	}

	return &ImageMetaData{
		ImageId: imgId,
		Title:   data.FormValue("title"),
		Content: data.FormValue("content"),
	}, nil
}

func (e *Image) Title() string {
	return e.meta.Title
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

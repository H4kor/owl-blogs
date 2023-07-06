package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/web/editor"
	"time"

	"github.com/gofiber/fiber/v2"
)

type EditorHandler struct {
	entrySvc *app.EntryService
}

func NewEditorHandler(entryService *app.EntryService) *EditorHandler {
	return &EditorHandler{entrySvc: entryService}
}

func (h *EditorHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	form := editor.NewEntryForm(&model.ImageEntry{})
	htmlForm, err := form.HtmlForm()
	if err != nil {
		return err
	}
	return c.SendString(htmlForm)
}

func (h *EditorHandler) HandlePost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	form := editor.NewEntryForm(&model.ImageEntry{})
	// get form data
	metaData, err := form.Parse(c)
	if err != nil {
		return err
	}

	// create entry
	now := time.Now()
	entry := &model.ImageEntry{}
	err = h.entrySvc.Create(entry, &now, metaData.MetaData())
	if err != nil {
		return err
	}

	return c.SendString("Hello, Editor!")
}

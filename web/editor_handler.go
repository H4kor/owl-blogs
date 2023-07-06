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
	registry *app.EntryTypeRegistry
}

func NewEditorHandler(entryService *app.EntryService, registry *app.EntryTypeRegistry) *EditorHandler {
	return &EditorHandler{entrySvc: entryService, registry: registry}
}

func (h *EditorHandler) paramToEntry(c *fiber.Ctx) (model.Entry, error) {
	typeName := c.Params("editor")
	entryType, err := h.registry.Type(typeName)
	if err != nil {
		return nil, err
	}
	return entryType, nil
}

func (h *EditorHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entryType, err := h.paramToEntry(c)
	if err != nil {
		return err
	}

	form := editor.NewEntryForm(entryType)
	htmlForm, err := form.HtmlForm()
	if err != nil {
		return err
	}
	return c.SendString(htmlForm)
}

func (h *EditorHandler) HandlePost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entryType, err := h.paramToEntry(c)
	if err != nil {
		return err
	}

	form := editor.NewEntryForm(entryType)
	// get form data
	metaData, err := form.Parse(c)
	if err != nil {
		return err
	}

	// create entry
	now := time.Now()
	entry := entryType
	err = h.entrySvc.Create(entry, &now, metaData.MetaData())
	if err != nil {
		return err
	}

	return c.SendString("Hello, Editor!")
}

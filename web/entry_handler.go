package web

import (
	"owl-blogs/app"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type EntryHandler struct {
	entrySvc *app.EntryService
	registry *app.EntryTypeRegistry
}

func NewEntryHandler(entryService *app.EntryService, registry *app.EntryTypeRegistry) *EntryHandler {
	return &EntryHandler{entrySvc: entryService, registry: registry}
}

func (h *EntryHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entryId := c.Params("post")
	entry, err := h.entrySvc.FindById(entryId)
	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(c, "views/entry", entry)
}

package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type EntryHandler struct {
	entrySvc  *app.EntryService
	authorSvc *app.AuthorService
	registry  *app.EntryTypeRegistry
}

type entryData struct {
	Entry  model.Entry
	Author *model.Author
}

func NewEntryHandler(entryService *app.EntryService, registry *app.EntryTypeRegistry, authorService *app.AuthorService) *EntryHandler {
	return &EntryHandler{entrySvc: entryService, authorSvc: authorService, registry: registry}
}

func (h *EntryHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entryId := c.Params("post")
	entry, err := h.entrySvc.FindById(entryId)
	if err != nil {
		return err
	}

	author, err := h.authorSvc.FindByName("h4kor")
	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(c, "views/entry", entryData{Entry: entry, Author: author})
}

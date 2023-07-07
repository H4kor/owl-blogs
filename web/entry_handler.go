package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"text/template"

	"github.com/gofiber/fiber/v2"
)

type EntryHandler struct {
	entrySvc *app.EntryService
	registry *app.EntryTypeRegistry
}

func NewEntryHandler(entryService *app.EntryService, registry *app.EntryTypeRegistry) *EntryHandler {
	return &EntryHandler{entrySvc: entryService, registry: registry}
}

func (h *EntryHandler) getTemplate(entry model.Entry) (*template.Template, error) {
	name, err := h.registry.TypeName(entry)
	if err != nil {
		return nil, err
	}
	return template.ParseFS(
		templates,
		"templates/base.tmpl",
		"templates/views/entry/"+name+".tmpl",
	)
}

func (h *EntryHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entryId := c.Params("post")
	entry, err := h.entrySvc.FindById(entryId)
	if err != nil {
		return err
	}

	template, err := h.getTemplate(entry)
	if err != nil {
		return err
	}

	return template.ExecuteTemplate(c, "base", entry)
}

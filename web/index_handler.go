package web

import (
	"owl-blogs/app"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type IndexHandler struct {
	entrySvc *app.EntryService
}

func NewIndexHandler(entryService *app.EntryService) *IndexHandler {
	return &IndexHandler{entrySvc: entryService}
}

func (h *IndexHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	entries, err := h.entrySvc.FindAll()

	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(c, "views/index", entries)

}

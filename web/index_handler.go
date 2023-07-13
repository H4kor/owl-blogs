package web

import (
	"owl-blogs/app"
	"owl-blogs/render"
	"sort"

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

	// sort entries by date descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PublishedAt().After(*entries[j].PublishedAt())
	})

	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(c, "views/index", entries)

}

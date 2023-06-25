package web

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type IndexHandler struct {
	entrySvc *app.EntryService
}

func NewIndexHandler(entryService *app.EntryService) *IndexHandler {
	return &IndexHandler{entrySvc: entryService}
}

func (h *IndexHandler) Handle(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

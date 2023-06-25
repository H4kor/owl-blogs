package web

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type MediaHandler struct {
	entrySvc *app.EntryService
}

func NewMediaHandler(entryService *app.EntryService) *MediaHandler {
	return &MediaHandler{entrySvc: entryService}
}

func (h *MediaHandler) Handle(c *fiber.Ctx) error {
	return c.SendString("Hello, Media!")
}

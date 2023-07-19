package web

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type RSSHandler struct {
	entrySvc *app.EntryService
}

func NewRSSHandler(entryService *app.EntryService) *RSSHandler {
	return &RSSHandler{entrySvc: entryService}
}

func (h *RSSHandler) Handle(c *fiber.Ctx) error {
	return c.SendString("Hello, RSS!")
}

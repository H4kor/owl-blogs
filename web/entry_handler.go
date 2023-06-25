package web

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type EntryHandler struct {
	entrySvc *app.EntryService
}

func NewEntryHandler(entryService *app.EntryService) *EntryHandler {
	return &EntryHandler{entrySvc: entryService}
}

func (h *EntryHandler) Handle(c *fiber.Ctx) error {
	return c.SendString("Hello, RSS!")
}

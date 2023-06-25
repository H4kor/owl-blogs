package web

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type ListHandler struct {
	entrySvc *app.EntryService
}

func NewListHandler(entryService *app.EntryService) *ListHandler {
	return &ListHandler{entrySvc: entryService}
}

func (h *ListHandler) Handle(c *fiber.Ctx) error {
	return c.SendString("Hello, List!")
}

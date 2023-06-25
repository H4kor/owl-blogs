package web

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type EditorHandler struct {
	entrySvc *app.EntryService
}

func NewEditorHandler(entryService *app.EntryService) *EditorHandler {
	return &EditorHandler{entrySvc: entryService}
}

func (h *EditorHandler) HandleGet(c *fiber.Ctx) error {
	return c.SendString("Hello, Editor!")
}

func (h *EditorHandler) HandlePost(c *fiber.Ctx) error {
	return c.SendString("Hello, Editor!")
}

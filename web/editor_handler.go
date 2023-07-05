package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/web/editor"

	"github.com/gofiber/fiber/v2"
)

type EditorHandler struct {
	entrySvc *app.EntryService
}

func NewEditorHandler(entryService *app.EntryService) *EditorHandler {
	return &EditorHandler{entrySvc: entryService}
}

func (h *EditorHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	formService := editor.NewEditorFormService(&model.ImageEntry{})
	form, _ := formService.HtmlForm()
	return c.SendString(form)
}

func (h *EditorHandler) HandlePost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.SendString("Hello, Editor!")
}

package web

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type LoginHandler struct {
	entrySvc *app.EntryService
}

func NewLoginHandler(entryService *app.EntryService) *LoginHandler {
	return &LoginHandler{entrySvc: entryService}
}

func (h *LoginHandler) HandleGet(c *fiber.Ctx) error {
	return c.SendString("Hello, Login!")
}

func (h *LoginHandler) HandlePost(c *fiber.Ctx) error {
	return c.SendString("Hello, Login!")
}

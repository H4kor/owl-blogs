package web

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type MediaHandler struct {
	binaryService *app.BinaryService
}

func NewMediaHandler(binaryService *app.BinaryService) *MediaHandler {
	return &MediaHandler{binaryService: binaryService}
}

func (h *MediaHandler) Handle(c *fiber.Ctx) error {
	id := c.Params("+")
	binary, err := h.binaryService.FindById(id)
	if err != nil {
		return err
	}
	return c.Send(binary.Data)
}

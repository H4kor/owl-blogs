package web

import (
	"net/url"
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
	// urldecode
	id, err := url.PathUnescape(id)
	if err != nil {
		return err
	}
	binary, err := h.binaryService.FindById(id)
	c.Set("Content-Type", binary.Mime())
	if err != nil {
		return err
	}
	return c.Send(binary.Data)
}

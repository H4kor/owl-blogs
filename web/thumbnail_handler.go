package web

import (
	"net/url"
	"owl-blogs/app"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ThumbnailHandler struct {
	binaryService    *app.BinaryService
	thumbnailService *app.ThumbnailService
}

func NewThumbnailHandler(binaryService *app.BinaryService, thumbnailService *app.ThumbnailService) *ThumbnailHandler {
	return &ThumbnailHandler{binaryService: binaryService, thumbnailService: thumbnailService}
}

func (h *ThumbnailHandler) Handle(c *fiber.Ctx) error {
	fileName := c.Params("+")

	dotIdx := strings.LastIndex(fileName, ".")
	if dotIdx == -1 {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	binaryFileId := fileName[:dotIdx]

	// urldecode
	binaryFileId, err := url.PathUnescape(binaryFileId)
	if err != nil {
		return err
	}

	thumbnail, err := h.thumbnailService.GetThumbnailForBinaryFileId(binaryFileId)

	if err != nil {
		binary, err := h.binaryService.FindById(binaryFileId)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		thumbnail, err = h.thumbnailService.CreateThumbnailForBinary(binary)
		if err != nil {
			// cannot create thumbnail for binary. Deliver binary instead
			c.Set("Content-Type", binary.Mime())
			return c.Send(binary.Data)
		}
	}

	c.Set("Content-Type", thumbnail.MimeType)
	return c.Send(thumbnail.Data)
}

package web

import (
	"log/slog"
	"net/url"
	"owl-blogs/app"

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
	binaryFileId := c.Params("+")
	// urldecode
	binaryFileId, err := url.PathUnescape(binaryFileId)
	if err != nil {
		return err
	}

	thumbnail, err := h.thumbnailService.GetThumbnailForBinaryFileId(binaryFileId)

	if err != nil {
		slog.Info("Could not get thumbnail", "error", err)
		binary, err := h.binaryService.FindById(binaryFileId)
		if err != nil {
			slog.Info("Could not find binary file", "error", err)
			return c.SendStatus(fiber.StatusNotFound)
		}
		slog.Info("Generating new thumbnail")
		thumbnail, err = h.thumbnailService.CreateThumbnailForBinary(binary)
		if err != nil {
			slog.Warn("Could not create thumbnail", "error", err)
			// cannot create thumbnail for binary. Deliver binary instead
			c.Set("Content-Type", binary.Mime())
			return c.Send(binary.Data)
		}
	}

	c.Set("Content-Type", thumbnail.MimeType)
	return c.Send(thumbnail.Data)
}

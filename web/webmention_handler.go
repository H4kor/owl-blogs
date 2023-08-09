package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"

	"github.com/gofiber/fiber/v2"
)

type WebmentionHandler struct {
	configRepo        repository.ConfigRepository
	webmentionService *app.WebmentionService
}

func NewWebmentionHandler(
	webmentionService *app.WebmentionService,
	configRepo repository.ConfigRepository,
) *WebmentionHandler {
	return &WebmentionHandler{
		webmentionService: webmentionService,
		configRepo:        configRepo,
	}
}

func (h *WebmentionHandler) Handle(c *fiber.Ctx) error {
	target := c.FormValue("target")
	source := c.FormValue("source")

	println("target", target)
	println("source", source)

	if target == "" {
		return c.Status(400).SendString("target is required")
	}
	if source == "" {
		return c.Status(400).SendString("source is required")
	}

	if len(target) < 7 || (target[:7] != "http://" && target[:8] != "https://") {
		return c.Status(400).SendString("target must be a valid URL")
	}

	if len(source) < 7 || (source[:7] != "http://" && source[:8] != "https://") {
		return c.Status(400).SendString("source must be a valid URL")
	}

	if source == target {
		return c.Status(400).SendString("source and target must be different")
	}

	err := h.webmentionService.ProcessWebmention(source, target)
	if err != nil {
		return err
	}

	return c.SendString("ok")

}

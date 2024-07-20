package web

import (
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type NotFoundPageData struct {
	Msg string
}

func Render404PageWithMessage(data NotFoundPageData, c *fiber.Ctx) error {

	return render.RenderTemplateWithBase(
		c.Status(fiber.StatusNotFound),
		"views/404",
		data,
	)
}

package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/render"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoginHandler struct {
	configRepo    repository.ConfigRepository
	authorService *app.AuthorService
}

func NewLoginHandler(
	authorService *app.AuthorService,
	configRepo repository.ConfigRepository,
) *LoginHandler {
	return &LoginHandler{
		authorService: authorService,
		configRepo:    configRepo,
	}
}

func (h *LoginHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return render.RenderTemplateWithBase(c, "views/login", nil)
}

func (h *LoginHandler) HandlePost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	name := c.FormValue("name")
	password := c.FormValue("password")

	valid := h.authorService.Authenticate(name, password)
	if !valid {
		return c.Redirect("/auth/login")
	}

	token, err := h.authorService.CreateToken(name)
	if err != nil {
		return err
	}

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Redirect("/admin/")

}

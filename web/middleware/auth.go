package middleware

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	authorService *app.AuthorService
}

func NewAuthMiddleware(authorService *app.AuthorService) *AuthMiddleware {
	return &AuthMiddleware{authorService: authorService}
}

func (m *AuthMiddleware) Handle(c *fiber.Ctx) error {
	// get token from cookie
	token := c.Cookies("token")
	if token == "" {
		return c.Redirect("/auth/login")
	}

	// check token
	valid, name := m.authorService.ValidateToken(token)
	if !valid {
		return c.Redirect("/auth/login")
	}

	// set author name to context
	c.Locals("author", name)

	return c.Next()
}

package middleware

import (
	"owl-blogs/app"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	authorService *app.AuthorService
}

type UserMiddleware struct {
	authorService *app.AuthorService
}

func NewAuthMiddleware(authorService *app.AuthorService) *AuthMiddleware {
	return &AuthMiddleware{authorService: authorService}
}

func NewUserMiddleware(authorService *app.AuthorService) *UserMiddleware {
	return &UserMiddleware{authorService: authorService}
}

func (m *AuthMiddleware) Handle(c *fiber.Ctx) error {
	if c.Locals("author") == nil {
		return c.Redirect("/auth/login")
	}

	return c.Next()
}

func (m *UserMiddleware) Handle(c *fiber.Ctx) error {
	// get token from cookie
	token := c.Cookies("token")
	if token == "" {
		return c.Next()
	}

	// check token
	valid, name := m.authorService.ValidateToken(token)
	if !valid {
		return c.Next()
	}

	// set author name to context
	c.Locals("author", name)

	return c.Next()
}

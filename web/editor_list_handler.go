package web

import (
	"embed"
	"owl-blogs/app"
	"text/template"

	"github.com/gofiber/fiber/v2"
)

//go:embed templates
var templates embed.FS

type EditorListHandler struct {
	registry *app.EntryTypeRegistry
	ts       *template.Template
}

type EditorListContext struct {
	Types []string
}

func NewEditorListHandler(registry *app.EntryTypeRegistry) *EditorListHandler {
	ts, err := template.ParseFS(
		templates,
		"templates/base.tmpl",
		"templates/views/editor_list.tmpl",
	)

	if err != nil {
		panic(err)
	}

	return &EditorListHandler{
		registry: registry,
		ts:       ts,
	}

}

func (h *EditorListHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	types := h.registry.Types()
	typeNames := []string{}

	for _, t := range types {
		name, _ := h.registry.TypeName(t)
		typeNames = append(typeNames, name)
	}

	return h.ts.ExecuteTemplate(c, "base", &EditorListContext{Types: typeNames})
}

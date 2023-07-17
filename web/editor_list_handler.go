package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type EditorListHandler struct {
	configRepo repository.SiteConfigRepository
	registry   *app.EntryTypeRegistry
}

type EditorListContext struct {
	Types []string
}

func NewEditorListHandler(registry *app.EntryTypeRegistry,
	configRepo repository.SiteConfigRepository) *EditorListHandler {
	return &EditorListHandler{
		registry:   registry,
		configRepo: configRepo,
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

	return render.RenderTemplateWithBase(c, getConfig(h.configRepo), "views/editor_list", &EditorListContext{Types: typeNames})
}

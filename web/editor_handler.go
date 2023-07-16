package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"owl-blogs/web/editor"
	"time"

	"github.com/gofiber/fiber/v2"
)

type EditorHandler struct {
	entrySvc *app.EntryService
	binSvc   *app.BinaryService
	registry *app.EntryTypeRegistry
}

func NewEditorHandler(
	entryService *app.EntryService,
	registry *app.EntryTypeRegistry,
	binService *app.BinaryService,
) *EditorHandler {
	return &EditorHandler{
		entrySvc: entryService,
		registry: registry,
		binSvc:   binService,
	}
}

func (h *EditorHandler) paramToEntry(c *fiber.Ctx) (model.Entry, error) {
	typeName := c.Params("editor")
	entryType, err := h.registry.Type(typeName)
	if err != nil {
		return nil, err
	}
	return entryType, nil
}

func (h *EditorHandler) HandleGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entryType, err := h.paramToEntry(c)
	if err != nil {
		return err
	}

	form := editor.NewEntryForm(entryType, h.binSvc)
	htmlForm, err := form.HtmlForm()
	if err != nil {
		return err
	}
	return render.RenderTemplateWithBase(c, "views/editor", htmlForm)
}

func (h *EditorHandler) HandlePost(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entryType, err := h.paramToEntry(c)
	if err != nil {
		return err
	}

	form := editor.NewEntryForm(entryType, h.binSvc)
	// get form data
	entry, err := form.Parse(c)
	if err != nil {
		return err
	}

	// create entry
	now := time.Now()
	entry.SetPublishedAt(&now)
	entry.SetAuthorId(c.Locals("author").(string))

	err = h.entrySvc.Create(entry)
	if err != nil {
		return err
	}
	return c.Redirect("/posts/" + entry.ID() + "/")

}

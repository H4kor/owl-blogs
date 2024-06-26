package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"time"

	"github.com/gofiber/fiber/v2"
)

type EditorHandler struct {
	configRepo repository.ConfigRepository
	entrySvc   *app.EntryService
	binSvc     *app.BinaryService
	registry   *app.EntryTypeRegistry
}

func NewEditorHandler(
	entryService *app.EntryService,
	registry *app.EntryTypeRegistry,
	binService *app.BinaryService,
	configRepo repository.ConfigRepository,
) *EditorHandler {
	return &EditorHandler{
		entrySvc:   entryService,
		registry:   registry,
		binSvc:     binService,
		configRepo: configRepo,
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

func (h *EditorHandler) HandleGetNew(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entryType, err := h.paramToEntry(c)
	if err != nil {
		return err
	}
	htmlForm := entryType.MetaData().Form(h.binSvc)
	return render.RenderTemplateWithBase(c, "views/editor", htmlForm)
}

func (h *EditorHandler) HandlePostNew(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entry, err := h.paramToEntry(c)
	if err != nil {
		return err
	}

	entryMeta := entry.MetaData()
	err = entryMeta.ParseFormData(c, h.binSvc)
	if err != nil {
		return err
	}

	// create entry
	now := time.Now()
	entry.SetMetaData(entryMeta)
	published := c.FormValue("action") == "Publish"
	if published {
		entry.SetPublishedAt(&now)
	} else {
		entry.SetPublishedAt(nil)
	}
	entry.SetAuthorId(c.Locals("author").(string))

	err = h.entrySvc.Create(entry)
	if err != nil {
		return err
	}
	return c.Redirect("/posts/" + entry.ID() + "/")

}

func (h *EditorHandler) HandleGetEdit(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	id := c.Params("id")
	entry, err := h.entrySvc.FindById(id)
	if err != nil {
		return err
	}

	htmlForm := entry.MetaData().Form(h.binSvc)
	return render.RenderTemplateWithBase(c, "views/editor", htmlForm)
}

func (h *EditorHandler) HandlePostEdit(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	id := c.Params("id")
	entry, err := h.entrySvc.FindById(id)
	if err != nil {
		return err
	}

	// get form data
	meta := entry.MetaData()
	err = meta.ParseFormData(c, h.binSvc)
	if err != nil {
		return err
	}

	published := c.FormValue("action") == "Publish"
	if !published {
		entry.SetPublishedAt(nil)
	} else if entry.PublishedAt() == nil || entry.PublishedAt().IsZero() {
		now := time.Now()
		entry.SetPublishedAt(&now)
	}
	// update entry
	entry.SetMetaData(meta)
	err = h.entrySvc.Update(entry)
	if err != nil {
		return err
	}
	return c.Redirect("/posts/" + entry.ID() + "/")
}

func (h *EditorHandler) HandlePostDelete(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	id := c.Params("id")
	entry, err := h.entrySvc.FindById(id)
	if err != nil {
		return err
	}

	confirm := c.FormValue("confirm")
	if confirm != "on" {
		return c.Redirect("/posts/" + entry.ID() + "/")
	}

	err = h.entrySvc.Delete(entry)
	if err != nil {
		return err
	}
	return c.Redirect("/")
}

func (h *EditorHandler) HandlePostUnpublish(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	id := c.Params("id")
	entry, err := h.entrySvc.FindById(id)
	if err != nil {
		return err
	}
	entry.SetPublishedAt(nil)
	err = h.entrySvc.Update(entry)
	if err != nil {
		return err
	}
	return c.Redirect("/posts/" + entry.ID() + "/")
}

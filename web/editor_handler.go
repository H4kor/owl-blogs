package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"owl-blogs/web/forms"
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

	form := forms.NewForm(entryType.MetaData(), h.binSvc)
	htmlForm, err := form.HtmlForm()
	if err != nil {
		return err
	}
	return render.RenderTemplateWithBase(c, getSiteConfig(h.configRepo), "views/editor", htmlForm)
}

func (h *EditorHandler) HandlePostNew(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	entry, err := h.paramToEntry(c)
	if err != nil {
		return err
	}

	form := forms.NewForm(entry.MetaData(), h.binSvc)
	// get form data
	entryMeta, err := form.Parse(c)
	if err != nil {
		return err
	}

	// create entry
	now := time.Now()
	entry.SetMetaData(entryMeta)
	entry.SetPublishedAt(&now)
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

	form := forms.NewForm(entry.MetaData(), h.binSvc)
	htmlForm, err := form.HtmlForm()
	if err != nil {
		return err
	}
	return render.RenderTemplateWithBase(c, getSiteConfig(h.configRepo), "views/editor", htmlForm)
}

func (h *EditorHandler) HandlePostEdit(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	id := c.Params("id")
	entry, err := h.entrySvc.FindById(id)
	if err != nil {
		return err
	}

	form := forms.NewForm(entry.MetaData(), h.binSvc)
	// get form data
	meta, err := form.Parse(c)
	if err != nil {
		return err
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

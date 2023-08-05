package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type EntryHandler struct {
	configRepo repository.ConfigRepository
	entrySvc   *app.EntryService
	authorSvc  *app.AuthorService
	registry   *app.EntryTypeRegistry
}

type entryData struct {
	Entry    model.Entry
	Author   *model.Author
	LoggedIn bool
}

func NewEntryHandler(
	entryService *app.EntryService,
	registry *app.EntryTypeRegistry,
	authorService *app.AuthorService,
	configRepo repository.ConfigRepository,
) *EntryHandler {
	return &EntryHandler{
		entrySvc:   entryService,
		authorSvc:  authorService,
		registry:   registry,
		configRepo: configRepo,
	}
}

func (h *EntryHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	loggedIn := c.Locals("author") != nil

	entryId := c.Params("post")
	entry, err := h.entrySvc.FindById(entryId)
	if err != nil {
		return err
	}

	if !loggedIn {
		if entry.PublishedAt() == nil || entry.PublishedAt().IsZero() {
			return fiber.NewError(fiber.StatusNotFound, "Entry not found")
		}
	}

	author, err := h.authorSvc.FindByName(entry.AuthorId())
	if err != nil {
		author = &model.Author{}
	}

	return render.RenderTemplateWithBase(
		c,
		getSiteConfig(h.configRepo),
		"views/entry",
		entryData{
			Entry:    entry,
			Author:   author,
			LoggedIn: loggedIn,
		},
	)
}

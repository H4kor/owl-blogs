package web

import (
	"errors"
	"net/http"
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type EntryHandler struct {
	configRepo      repository.ConfigRepository
	entrySvc        *app.EntryService
	authorSvc       *app.AuthorService
	registry        *app.EntryTypeRegistry
	interactionRepo repository.InteractionRepository
}

type entryData struct {
	Entry        model.Entry
	Author       *model.Author
	LoggedIn     bool
	Interactions []model.Interaction
}

func NewEntryHandler(
	entryService *app.EntryService,
	registry *app.EntryTypeRegistry,
	authorService *app.AuthorService,
	configRepo repository.ConfigRepository,
	interactionRepo repository.InteractionRepository,
) *EntryHandler {
	return &EntryHandler{
		entrySvc:        entryService,
		authorSvc:       authorService,
		registry:        registry,
		configRepo:      configRepo,
		interactionRepo: interactionRepo,
	}
}

func (h *EntryHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	url := c.OriginalURL()
	if len(url) == 0 || url[len(url)-1] != '/' {
		return c.Redirect(url+"/", http.StatusMovedPermanently)
	}

	loggedIn := c.Locals("author") != nil

	entryId := c.Params("post")
	entry, err := h.entrySvc.FindById(entryId)
	if err != nil {
		if errors.Is(err, app.ErrEntryNotFound) {
			return Render404PageWithMessage(NotFoundPageData{
				Msg: "Entry not found",
			}, c)
		}
		return err
	}

	if !loggedIn {
		if entry.PublishedAt() == nil || entry.PublishedAt().IsZero() {
			return Render404PageWithMessage(NotFoundPageData{
				Msg: "Entry not found",
			}, c)
		}
	}

	author, err := h.authorSvc.FindByName(entry.AuthorId())
	if err != nil {
		author = &model.Author{}
	}

	inters, _ := h.interactionRepo.FindAll(entry.ID())

	return render.RenderTemplateWithBase(
		c,
		"views/entry",
		entryData{
			Entry:        entry,
			Author:       author,
			LoggedIn:     loggedIn,
			Interactions: inters,
		},
	)
}

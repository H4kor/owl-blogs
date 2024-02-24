package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type AdminInteractionHandler struct {
	interactionRepo repository.InteractionRepository
	configRepo      repository.ConfigRepository
}

func NewAdminInteractionHandler(configRepo repository.ConfigRepository, interactionRepo repository.InteractionRepository) *AdminInteractionHandler {
	return &AdminInteractionHandler{
		interactionRepo: interactionRepo,
		configRepo:      configRepo,
	}
}

func (h *AdminInteractionHandler) HandleGet(c *fiber.Ctx) error {
	siteConfig := getSiteConfig(h.configRepo)

	filter := c.Query("filter", "")

	interactions, err := h.interactionRepo.ListAllInteractions()
	if err != nil {
		return err
	}
	pageData := paginate(c, interactions, 50)

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return render.RenderTemplateWithBase(c, siteConfig, "views/interaction_manager", fiber.Map{
		"Interactions": pageData.items,
		"Page":         pageData.page,
		"NextPage":     pageData.page + 1,
		"PrevPage":     pageData.page - 1,
		"FirstPage":    pageData.page == 1,
		"LastPage":     pageData.lastPage,
		"Filter":       filter,
	})

}

func (h *AdminInteractionHandler) HandleDelete(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	id := c.Params("id")
	inter, err := h.interactionRepo.FindById(id)
	entryId := inter.EntryID()
	if err != nil {
		return err
	}

	confirm := c.FormValue("confirm")
	if confirm != "on" {
		return c.Redirect("/posts/" + inter.ID() + "/")
	}

	err = h.interactionRepo.Delete(inter)
	if err != nil {
		return err
	}
	return c.Redirect("/posts/" + entryId + "/")
}

package web

import (
	"owl-blogs/app/repository"

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

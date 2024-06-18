package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/render"

	"github.com/gofiber/fiber/v2"
)

type AdminFollowerHandler struct {
	interactionRepo repository.FollowerRepository
	configRepo      repository.ConfigRepository
}

func NewAdminFollowerHandler(configRepo repository.ConfigRepository, interactionRepo repository.FollowerRepository) *AdminFollowerHandler {
	return &AdminFollowerHandler{
		interactionRepo: interactionRepo,
		configRepo:      configRepo,
	}
}

func (h *AdminFollowerHandler) HandleGet(c *fiber.Ctx) error {
	filter := c.Query("filter", "")

	Followers, err := h.interactionRepo.All()
	if err != nil {
		return err
	}
	pageData := paginate(c, Followers, 50)

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return render.RenderTemplateWithBase(c, "views/follower_manager", fiber.Map{
		"Followers": pageData.items,
		"Page":      pageData.page,
		"NextPage":  pageData.page + 1,
		"PrevPage":  pageData.page - 1,
		"FirstPage": pageData.page == 1,
		"LastPage":  pageData.lastPage,
		"Filter":    filter,
	})

}

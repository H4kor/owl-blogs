package web

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type BinaryManageHandler struct {
	configRepo repository.ConfigRepository
	service    *app.BinaryService
}

func NewBinaryManageHandler(configRepo repository.ConfigRepository, service *app.BinaryService) *BinaryManageHandler {
	return &BinaryManageHandler{
		configRepo: configRepo,
		service:    service,
	}
}

func (h *BinaryManageHandler) Handle(c *fiber.Ctx) error {
	siteConfig := getSiteConfig(h.configRepo)

	filter := c.Query("filter", "")

	allIds, err := h.service.ListIds(filter)
	sort.Slice(allIds, func(i, j int) bool {
		return strings.ToLower(allIds[i]) < strings.ToLower(allIds[j])
	})
	if err != nil {
		return err
	}
	pageData := paginate(c, allIds, 50)

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return render.RenderTemplateWithBase(c, siteConfig, "views/binary_manager", fiber.Map{
		"Binaries":  pageData.items,
		"Page":      pageData.page,
		"NextPage":  pageData.page + 1,
		"PrevPage":  pageData.page - 1,
		"FirstPage": pageData.page == 1,
		"LastPage":  pageData.lastPage,
		"Filter":    filter,
	})

}

func (h *BinaryManageHandler) saveFileUpload(c *fiber.Ctx) (*model.BinaryFile, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	content := make([]byte, file.Size)
	_, err = reader.Read(content)
	if err != nil {
		return nil, err
	}

	binary, err := h.service.Create(file.Filename, content)
	if err != nil {
		return nil, err
	}
	return binary, nil
}

func (h *BinaryManageHandler) HandleUpload(c *fiber.Ctx) error {
	binary, err := h.saveFileUpload(c)
	if err != nil {
		return err
	}
	return c.Redirect("/media/" + binary.Id)
}

func (h *BinaryManageHandler) HandleUploadApi(c *fiber.Ctx) error {
	binary, err := h.saveFileUpload(c)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{
		"location": "/media/" + binary.Id,
	})
}

func (h *BinaryManageHandler) HandleDelete(c *fiber.Ctx) error {
	id := c.FormValue("file")
	binary, err := h.service.FindById(id)
	if err != nil {
		return err
	}

	confirm := c.FormValue("confirm")
	if confirm != "on" {
		return c.Redirect("/admin/binaries/")
	}

	err = h.service.Delete(binary)
	if err != nil {
		return err
	}

	return c.Redirect("/admin/binaries/")
}

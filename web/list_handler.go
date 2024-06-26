package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ListHandler struct {
	siteConfigService *app.SiteConfigService
	entrySvc          *app.EntryService
}

func NewListHandler(
	entryService *app.EntryService,
	siteConfigService *app.SiteConfigService,
) *ListHandler {
	return &ListHandler{
		entrySvc:          entryService,
		siteConfigService: siteConfigService,
	}
}

type listRenderData struct {
	List      model.EntryList
	Entries   []model.Entry
	Page      int
	NextPage  int
	PrevPage  int
	FirstPage bool
	LastPage  bool
}

func (h *ListHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()
	if err != nil {
		return err
	}
	listId := c.Params("list")
	list := model.EntryList{}
	for _, l := range siteConfig.Lists {
		if l.Id == listId {
			list = l
		}
	}
	if list.Id == "" {
		return c.SendStatus(404)
	}

	entries, err := h.entrySvc.FindAllByType(&list.Include, true, false)
	if err != nil {
		return err
	}

	// sort entries by date descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PublishedAt().After(*entries[j].PublishedAt())
	})

	// pagination
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 1
	}
	limit := 10
	offset := (pageNum - 1) * limit
	lastPage := false
	if offset > len(entries) {
		offset = len(entries)
		lastPage = true
	}
	if offset+limit > len(entries) {
		limit = len(entries) - offset
		lastPage = true
	}
	entries = entries[offset : offset+limit]

	if err != nil {
		return err
	}

	return render.RenderTemplateWithBase(c, "views/list", listRenderData{
		List:      list,
		Entries:   entries,
		Page:      pageNum,
		NextPage:  pageNum + 1,
		PrevPage:  pageNum - 1,
		FirstPage: pageNum == 1,
		LastPage:  lastPage,
	})

}

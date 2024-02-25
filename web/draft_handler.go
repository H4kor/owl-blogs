package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type DraftHandler struct {
	siteConfigService *app.SiteConfigService
	entrySvc          *app.EntryService
}

func NewDraftHandler(
	entryService *app.EntryService,
	siteConfigService *app.SiteConfigService,
) *DraftHandler {
	return &DraftHandler{
		entrySvc:          entryService,
		siteConfigService: siteConfigService,
	}
}

type DraftRenderData struct {
	Entries   []model.Entry
	Page      int
	NextPage  int
	PrevPage  int
	FirstPage bool
	LastPage  bool
}

func (h *DraftHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	siteConfig, err := h.siteConfigService.GetSiteConfig()
	if err != nil {
		return err
	}

	entries, err := h.entrySvc.FindAllByType(&siteConfig.PrimaryListInclude, false, true)
	if err != nil {
		return err
	}

	// sort entries by date descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Title() < entries[j].Title()
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

	return render.RenderTemplateWithBase(c, "views/draft_list", DraftRenderData{
		Entries:   entries,
		Page:      pageNum,
		NextPage:  pageNum + 1,
		PrevPage:  pageNum - 1,
		FirstPage: pageNum == 1,
		LastPage:  lastPage,
	})

}

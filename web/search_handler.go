package web

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/render"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	siteConfigService *app.SiteConfigService
	entrySvc          *app.EntryService
}

type searchRenderData struct {
	SearchTerm string
	Entries    []model.Entry
	Page       int
	NextPage   int
	PrevPage   int
	FirstPage  bool
	LastPage   bool
}

func NewSearchHandler(
	entryService *app.EntryService,
	siteConfigService *app.SiteConfigService,
) *SearchHandler {
	return &SearchHandler{
		entrySvc:          entryService,
		siteConfigService: siteConfigService,
	}
}

func (h *SearchHandler) Handle(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	var entries []model.Entry
	var err error

	// search all entries
	searchTerm := c.Query("query")
	if searchTerm != "" {
		entries, err = h.entrySvc.SearchEntries(searchTerm)
		if err != nil {
			return err
		}
	} else {
		entries = make([]model.Entry, 0)
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
		// if the page number is not an integer -> remove query param by redirect
		return c.Redirect(c.Path(), 301)
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

	return render.RenderTemplateWithBase(c, "views/search", searchRenderData{
		SearchTerm: searchTerm,
		Entries:    entries,
		Page:       pageNum,
		NextPage:   pageNum + 1,
		PrevPage:   pageNum - 1,
		FirstPage:  pageNum == 1,
		LastPage:   lastPage,
	})

}

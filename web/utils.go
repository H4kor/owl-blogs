package web

import (
	"owl-blogs/app/repository"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getSiteConfig(repo repository.ConfigRepository) model.SiteConfig {
	siteConfig := model.SiteConfig{}
	err := repo.Get(config.SITE_CONFIG, &siteConfig)
	if err != nil {
		panic(err)
	}
	return siteConfig
}

type paginationData[T any] struct {
	items    []T
	page     uint
	lastPage bool
}

func paginate[T any](c *fiber.Ctx, items []T, limit int) paginationData[T] {
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 1
	}
	offset := (pageNum - 1) * limit
	lastPage := false
	if offset > len(items) {
		offset = len(items)
		lastPage = true
	}
	if offset+limit > len(items) {
		limit = len(items) - offset
		lastPage = true
	}
	items = items[offset : offset+limit]

	return paginationData[T]{
		items:    items,
		page:     uint(pageNum),
		lastPage: lastPage,
	}
}

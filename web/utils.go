package web

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

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

func isActivityPub(ctx *fiber.Ctx) bool {
	slog.Info(
		"AP",
		"accept", ctx.Request().Header.Peek("Accept"),
		"ct", ctx.Request().Header.Peek("Content-Type"))
	accepts := (strings.Contains(string(ctx.Request().Header.Peek("Accept")), "application/activity+json") ||
		strings.Contains(string(ctx.Request().Header.Peek("Accept")), "application/ld+json"))
	req_content := (strings.Contains(string(ctx.Request().Header.Peek("Content-Type")), "application/activity+json") ||
		strings.Contains(string(ctx.Request().Header.Peek("Content-Type")), "application/ld+json"))
	return accepts || req_content
}

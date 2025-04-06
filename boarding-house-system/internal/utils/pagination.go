package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type Pagination struct {
	Page     int
	PageSize int
}

func GetPagination(ctx fiber.Ctx) Pagination {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

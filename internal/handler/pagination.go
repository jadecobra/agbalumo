package handler

import (
	"strconv"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

// Pagination holds pagination parameters.
type Pagination struct {
	Page        int
	Limit       int
	Offset      int
	HasNextPage bool
}

// GetPagination extracts pagination parameters from query string.
func GetPagination(c echo.Context, defaultLimit int) Pagination {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit := defaultLimit
	offset := (page - 1) * limit
	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// ConvertCounts converts category counts map to string keys and calculates total.
func ConvertCounts(counts map[domain.Category]int) (map[string]int, int) {
	strCounts := make(map[string]int)
	totalCount := 0
	for cat, count := range counts {
		strCounts[string(cat)] = count
		totalCount += count
	}
	return strCounts, totalCount
}

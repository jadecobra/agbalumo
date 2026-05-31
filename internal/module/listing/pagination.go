package listing

import (
	"strconv"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

// GetPagination extracts pagination parameters from query string.
func GetPagination(c echo.Context, defaultLimit int) domain.Pagination {
	page, _ := strconv.Atoi(c.QueryParam(domain.ParamPage))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = defaultLimit
	}
	offset := (page - 1) * limit
	return domain.Pagination{
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

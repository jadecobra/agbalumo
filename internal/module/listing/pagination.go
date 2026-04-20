package listing

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
	TotalCount  int
	TotalPages  int
	HasNextPage bool
}

// GetPageRange returns a slice of page numbers to display in the UI.
func (p Pagination) GetPageRange() []int {
	if p.TotalPages <= 1 {
		return nil
	}

	// Show up to 5 surrounding pages
	start := p.Page - 2
	if start < 1 {
		start = 1
	}
	end := start + 4
	if end > p.TotalPages {
		end = p.TotalPages
		start = end - 4
		if start < 1 {
			start = 1
		}
	}

	var pages []int
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	return pages
}

// GetPagination extracts pagination parameters from query string.
func GetPagination(c echo.Context, defaultLimit int) Pagination {
	page, _ := strconv.Atoi(c.QueryParam(domain.ParamPage))
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

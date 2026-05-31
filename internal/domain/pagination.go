package domain

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

package listing

import (
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestPagination_GetPageRange(t *testing.T) {
	tests := []struct {
		name       string
		expected   []int
		pagination Pagination
	}{
		{
			name: "SinglePage",
			pagination: Pagination{
				Page:       1,
				TotalPages: 1,
			},
			expected: nil,
		},
		{
			name: "TwoPages",
			pagination: Pagination{
				Page:       1,
				TotalPages: 2,
			},
			expected: []int{1, 2},
		},
		{
			name: "FirstPage_ManyPages",
			pagination: Pagination{
				Page:       1,
				TotalPages: 10,
			},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name: "MiddlePage_ManyPages",
			pagination: Pagination{
				Page:       5,
				TotalPages: 10,
			},
			expected: []int{3, 4, 5, 6, 7},
		},
		{
			name: "LastPage_ManyPages",
			pagination: Pagination{
				Page:       10,
				TotalPages: 10,
			},
			expected: []int{6, 7, 8, 9, 10},
		},
		{
			name: "NearEnd_ManyPages",
			pagination: Pagination{
				Page:       9,
				TotalPages: 10,
			},
			expected: []int{6, 7, 8, 9, 10},
		},
		{
			name: "SecondPage_SmallCount",
			pagination: Pagination{
				Page:       2,
				TotalPages: 3,
			},
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pagination.GetPageRange())
		})
	}
}

func TestConvertCounts(t *testing.T) {
	counts := map[domain.Category]int{
		domain.Business: 10,
		domain.Job:      5,
	}

	strCounts, total := ConvertCounts(counts)

	assert.Equal(t, 15, total)
	assert.Equal(t, 10, strCounts[string(domain.Business)])
	assert.Equal(t, 5, strCounts[string(domain.Job)])
}

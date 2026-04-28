package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildListingWhere_LocationFallback(t *testing.T) {
	tests := []struct {
		name          string
		expectedWhere string
		filters       ListingFilters
	}{
		{
			name: "fallback to city string match when radius search has zero coordinates",
			filters: ListingFilters{
				City:        "Dallas",
				Radius:      10,
				IncludedLat: 0,
				IncludedLng: 0,
			},
			expectedWhere: "(city = ? OR address LIKE ?)",
		},
	}

	repo := &SQLiteRepository{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			where, _ := repo.buildListingWhere(tt.filters)
			assert.Contains(t, where, tt.expectedWhere)
		})
	}
}

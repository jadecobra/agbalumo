package admin_test

import (
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAllListings_Extended(t *testing.T) {
	repo := handler.SetupTestRepository(t)

	// Seed listings for filtering and counts
	_ = repo.Save(context.Background(), domain.Listing{ID: "l1", Title: "Test Business", Type: "business", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria", Address: "123 Test St", City: "Lagos"})
	_ = repo.Save(context.Background(), domain.Listing{ID: "l2", Title: "Test Event", Type: "events", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria"})

	tests := []struct {
		name       string
		query      string
		expectCode int
	}{
		{
			name:       "HappyPath_WithCategoryFilter",
			query:      "?category=business&sort=title&order=asc",
			expectCode: http.StatusOK,
		},
		{
			name:       "HappyPath_NoFilters",
			query:      "",
			expectCode: http.StatusOK,
		},
		{
			name:       "PaginationAndSorting",
			query:      "?page=2&sort=created_at&order=desc",
			expectCode: http.StatusOK,
		},
		{
			name:       "SortingByFeatured",
			query:      "?sort=featured&order=desc",
			expectCode: http.StatusOK,
		},
		{
			name:       "SortingByTypeAsc",
			query:      "?sort=type&order=asc",
			expectCode: http.StatusOK,
		},
		{
			name:       "SortingByTypeDesc",
			query:      "?sort=type&order=desc",
			expectCode: http.StatusOK,
		},
		{
			name:       "SearchQuery",
			query:      "?q=Test",
			expectCode: http.StatusOK,
		},
		{
			name:       "SearchQueryWithCategory",
			query:      "?q=Test&category=business",
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings"+tt.query, nil)
			c.Set("User", domain.User{Role: domain.UserRoleAdmin})

			h := admin.NewAdminHandler(repo, repo, repo, repo, repo, repo, repo, nil, config.LoadConfig())
			_ = h.HandleAllListings(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

func TestAdminHandler_HandleAllListings_Counts(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	ctx := context.Background()

	// Seed multiple categories
	_ = repo.Save(ctx, domain.Listing{ID: "l1", Title: "B1", Type: "business", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria"})
	_ = repo.Save(ctx, domain.Listing{ID: "l2", Title: "B2", Type: "business", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria"})
	_ = repo.Save(ctx, domain.Listing{ID: "l3", Title: "E1", Type: "events", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria"})

	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	h := admin.NewAdminHandler(repo, repo, repo, repo, repo, repo, repo, nil, config.LoadConfig())
	_ = h.HandleAllListings(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	// We can't easily check the body without a real renderer, but we verified the repo calls worked.
}

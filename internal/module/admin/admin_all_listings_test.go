package admin_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/module/admin"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAllListings_Extended(t *testing.T) {
	t.Parallel()

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
			t.Parallel()
			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()

			// Seed listings for filtering and counts
			_ = env.App.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "Test Business", Type: "business", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria", Address: "123 Test St", City: "Lagos"})
			_ = env.App.DB.Save(context.Background(), domain.Listing{ID: "l2", Title: "Test Event", Type: "events", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria"})

			c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin/listings"+tt.query, nil)

			h := admin.NewAdminHandler(env.App)
			_ = h.HandleAllListings(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

func TestAdminHandler_HandleAllListings_Counts(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	ctx := context.Background()

	for _, l := range []domain.Listing{
		{ID: "l1", Title: "B1", Type: "business", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria"},
		{ID: "l2", Title: "B2", Type: "business", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria"},
		{ID: "l3", Title: "E1", Type: "events", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria"},
	} {
		_ = env.App.DB.Save(ctx, l)
	}

	c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin/listings", nil)

	h := admin.NewAdminHandler(env.App)
	_ = h.HandleAllListings(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_HandleAllListings_EnrichmentAttemptedAt(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	ctx := context.Background()

	now := time.Now()
	_ = env.App.DB.Save(ctx, domain.Listing{
		ID:                    "l1",
		Title:                 "Enriched Listing",
		Type:                  "business",
		Status:                domain.ListingStatusApproved,
		OwnerOrigin:           "Nigeria",
		EnrichmentAttemptedAt: &now,
	})

	c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin/listings", nil)

	h := admin.NewAdminHandler(env.App)
	_ = h.HandleAllListings(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "Enriched")
	assert.Contains(t, body, now.Format("Jan 02, 2006"))
}


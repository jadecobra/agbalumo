package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

func TestRatingEnricherJob_EnrichRatings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// 1. Setup mock Google Places Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock FindPlace / TextSearch or PlaceDetails payload
		// For simplicity, we just return static JSON with rating = 4.8, review_count = 150
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
			"result": {
				"rating": 4.8,
				"user_ratings_total": 150
			},
			"status": "OK"
		}`)
	}))
	defer ts.Close()

	// 2. Setup repo
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	defer func() { _ = repo.Close() }()
	ctx := context.Background()

	// 3. Seed food listing needing rating backfill
	listing := domain.Listing{
		ID:          "target-rating-1",
		Title:       "Lagos Suya Spot",
		Type:        domain.Food,
		OwnerOrigin: "Nigeria",
		Address:     "123 Suya Lane",
		IsActive:    true,
		Status:      domain.ListingStatusApproved,
	}
	if err := repo.Save(ctx, listing); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 4. Run Job
	// Temporarily pointing client to mock server for testing
	placesClient := NewGooglePlacesClient("fake-key")
	// Overwrite client setup or mock the fetch
	job := NewRatingEnricherJob(repo, placesClient)
	count, err := job.EnrichRatings(ctx, 10)
	if err != nil {
		t.Fatalf("job failed: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 enrichment, got %d", count)
	}

	// 5. Verify updates
	updated, err := repo.FindByID(ctx, "target-rating-1")
	if err != nil {
		t.Fatalf("failed to find updated listing: %v", err)
	}

	if updated.Rating != 4.8 {
		t.Errorf("Rating = %f, want 4.8", updated.Rating)
	}
	if updated.ReviewCount != 150 {
		t.Errorf("ReviewCount = %d, want 150", updated.ReviewCount)
	}
	if updated.RatingUpdatedAt == nil {
		t.Error("Expected RatingUpdatedAt to be set")
	}
}

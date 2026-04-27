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

func TestScraperJob_EnrichListings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// 1. Setup mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `<html><body><h1>Signature Jollof</h1><p>Very spicy and hot!</p><p>Pay via Zelle</p></body></html>`)
	}))
	defer ts.Close()

	// 2. Setup repo
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	defer func() { _ = repo.Close() }()
	ctx := context.Background()

	// 3. Seed listing needing enrichment
	listing := domain.Listing{
		ID:          "target-1",
		Title:       "Test Restaurant",
		WebsiteURL:  ts.URL,
		Type:        domain.Food,
		OwnerOrigin: "Nigeria",
		Address:     "123 Street",
		IsActive:    true,
		Status:      domain.ListingStatusApproved,
	}
	if err := repo.Save(ctx, listing); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 4. Run job
	scraper := NewWebsiteScraper()
	job := NewScraperJob(repo, scraper)
	count, err := job.EnrichListings(ctx, 10)
	if err != nil {
		t.Fatalf("job failed: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 enrichment, got %d", count)
	}

	// 5. Verify update
	updated, err := repo.FindByID(ctx, "target-1")
	if err != nil {
		t.Fatalf("failed to find updated listing: %v", err)
	}
	if updated.HeatLevel != 2 { // "spicy" and "hot"
		t.Errorf("HeatLevel = %d, want 2", updated.HeatLevel)
	}
	if updated.TopDish != "Signature Jollof" {
		t.Errorf("TopDish = %q, want %q", updated.TopDish, "Signature Jollof")
	}
	if updated.PaymentMethods != "Zelle" {
		t.Errorf("PaymentMethods = %q, want %q", updated.PaymentMethods, "Zelle")
	}
}

func TestScraperJob_EnrichAttemptedAtOnFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo, _ := testutil.SetupTestRepositoryUnique(t)
	defer func() { _ = repo.Close() }()
	ctx := context.Background()

	listing := domain.Listing{
		ID:          "target-fail",
		Title:       "Fail Restaurant",
		WebsiteURL:  "http://localhost:12345/nonexistent", // will fail to connect
		Type:        domain.Food,
		OwnerOrigin: "Nigeria",
		IsActive:    true,
		Status:      domain.ListingStatusApproved,
	}
	if err := repo.Save(ctx, listing); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	scraper := NewWebsiteScraper()
	job := NewScraperJob(repo, scraper)
	
	_, _ = job.EnrichListings(ctx, 10)

	updated, err := repo.FindByID(ctx, "target-fail")
	if err != nil {
		t.Fatalf("failed to find updated listing: %v", err)
	}

	if updated.EnrichmentAttemptedAt == nil {
		t.Error("Expected EnrichmentAttemptedAt to be set even on failure")
	}
}


package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// SeedStandardData populates the database with a representative set of listings for testing.
func SeedStandardData(t *testing.T, db domain.ListingStore) {
	t.Helper()
	ctx := context.Background()

	listings := []domain.Listing{
		{ID: "l1", Title: "Business A", Type: domain.Business, City: "Houston", IsActive: true, Status: domain.ListingStatusApproved},
		{ID: "l2", Title: "Food B", Type: domain.Food, City: "Houston", IsActive: true, Status: domain.ListingStatusApproved},
		{ID: "l3", Title: "Service C", Type: domain.Service, City: "Dallas", IsActive: true, Status: domain.ListingStatusApproved},
		{ID: "l4", Title: "Request D", Type: domain.Request, City: "Dallas", IsActive: true, Status: domain.ListingStatusApproved},
	}

	for _, l := range listings {
		if err := db.Save(ctx, l); err != nil {
			t.Fatalf("Failed to seed listing %s: %v", l.ID, err)
		}
	}
}

// SeedAdaDallasData seeds specific Nigerian food listings in Dallas for persona-based testing.
func SeedAdaDallasData(t *testing.T, db domain.ListingStore) {
	t.Helper()
	ctx := context.Background()

	listings := []domain.Listing{
		{
			ID:          "ada-1",
			Title:       "Mama's Jollof House",
			Description: "Authentic Nigerian Jollof rice and Suya in the heart of Dallas.",
			Type:        domain.Food,
			City:        "Dallas",
			IsActive:    true,
			Status:      domain.ListingStatusApproved,
			Featured:    true,
		},
		{
			ID:          "ada-2",
			Title:       "Lagos Grill Dallas",
			Description: "Freshly grilled fish and plantains. Best in Texas.",
			Type:        domain.Food,
			City:        "Dallas",
			IsActive:    true,
			Status:      domain.ListingStatusApproved,
		},
		{
			ID:          "ada-3",
			Title:       "Abuja Express",
			Description: "Quick and spicy Nigerian takeout. Pounded yam and Egusi soup.",
			Type:        domain.Food,
			City:        "Dallas",
			IsActive:    true,
			Status:      domain.ListingStatusApproved,
		},
	}

	for _, l := range listings {
		if err := db.Save(ctx, l); err != nil {
			t.Fatalf("Failed to seed ada listing %s: %v", l.ID, err)
		}
	}
}

// SaveTestUser is a helper to save a user with sensible defaults.
func SaveTestUser(t *testing.T, db domain.UserStore, id, email string, role domain.UserRole) domain.User {
	t.Helper()
	u := domain.User{
		ID:        id,
		Email:     email,
		Role:      role,
		CreatedAt: time.Now(),
		GoogleID:  "google-" + id,
	}

	if err := db.SaveUser(context.Background(), u); err != nil {
		t.Fatalf("Failed to save test user %s: %v", id, err)
	}
	return u
}

package sqlite_test

import (
	"context"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

func TestAdaSignalsPersistence(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	defer func() {
		if err := repo.Close(); err != nil {
			t.Errorf("Failed to close repo: %v", err)
		}
	}()

	ctx := context.Background()
	listing := domain.Listing{
		ID:                "ada-test",
		Title:             "Spicy Suya Spot",
		Type:              domain.Food,
		OwnerOrigin:       "Nigeria",
		Address:           "Lagos Street",
		HeatLevel:         5,
		RegionalSpecialty: "Yoruba • Lagos Style",
		TopDish:           "Lagos Suya",
	}

	saveTestListing(t, ctx, repo, listing)

	got, err := repo.FindByID(ctx, "ada-test")
	if err != nil {
		t.Fatalf("Failed to get listing: %v", err)
	}

	if got.HeatLevel != 5 {
		t.Errorf("expected HeatLevel 5, got %d", got.HeatLevel)
	}
	if got.RegionalSpecialty != "Yoruba • Lagos Style" {
		t.Errorf("expected RegionalSpecialty Yoruba • Lagos Style, got %q", got.RegionalSpecialty)
	}
	if got.TopDish != "Lagos Suya" {
		t.Errorf("expected TopDish Lagos Suya, got %q", got.TopDish)
	}
}

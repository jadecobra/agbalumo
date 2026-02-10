package seeder_test

import (
	"context"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/jadecobra/agbalumo/internal/seeder"
)

func TestSeedAll(t *testing.T) {
	saveCount := 0
	mockRepo := &mock.MockListingRepository{
		SaveFn: func(ctx context.Context, l domain.Listing) error {
			saveCount++
			return nil
		},
	}

	seeder.SeedAll(context.Background(), mockRepo)

	// We expect a significant number of listings to be seeded.
	// As of writing, there are roughly 10 per category * 7 categories = ~70
	if saveCount < 50 {
		t.Errorf("Expected at least 50 listings seeded, got %d", saveCount)
	}
}

func TestEnsureSeeded_Empty(t *testing.T) {
	saveCalled := false
	mockRepo := &mock.MockListingRepository{
		FindAllFn: func(ctx context.Context, filterType, queryText string, includeInactive bool) ([]domain.Listing, error) {
			return []domain.Listing{}, nil // Empty
		},
		SaveFn: func(ctx context.Context, l domain.Listing) error {
			saveCalled = true
			return nil
		},
	}

	seeder.EnsureSeeded(context.Background(), mockRepo)

	if !saveCalled {
		t.Error("Expected Save to be called when database is empty")
	}
}

func TestEnsureSeeded_NotEmpty(t *testing.T) {
	saveCalled := false
	mockRepo := &mock.MockListingRepository{
		FindAllFn: func(ctx context.Context, filterType, queryText string, includeInactive bool) ([]domain.Listing, error) {
			return []domain.Listing{{Title: "Existing"}}, nil
		},
		SaveFn: func(ctx context.Context, l domain.Listing) error {
			saveCalled = true
			return nil
		},
	}

	seeder.EnsureSeeded(context.Background(), mockRepo)

	if saveCalled {
		t.Error("Expected Save NOT to be called when database is not empty")
	}
}

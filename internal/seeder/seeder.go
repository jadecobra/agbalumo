package seeder

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
)

//go:embed listings.json
var seedData embed.FS

type seedSource struct {
	Group    string           `json:"group"`
	Listings []domain.Listing `json:"listings"`
}

// SeedAll inserts all predefined data into the repository.
func SeedAll(ctx context.Context, repo domain.ListingStore) {
	sources, err := getListingData()
	if err != nil {
		slog.Error("Failed to load seed data", "error", err)
		return
	}
	for _, source := range sources {
		seedGroup(ctx, repo, source.Group, source.Listings)
	}
}

// EnsureSeeded checks if the database is empty, and if so, seeds it.
func EnsureSeeded(ctx context.Context, repo domain.ListingStore) {
	listings, _, err := repo.FindAll(ctx, "", "", "", "", "", true, 1, 0)
	if err != nil {
		slog.Error("Failed to check existing listings", "error", err)
		return
	}

	if len(listings) == 0 {
		slog.Info("Database empty. Seeding data...")
		SeedAll(ctx, repo)
	}
}

func seedGroup(ctx context.Context, repo domain.ListingStore, name string, listings []domain.Listing) {
	slog.Info("Seeding", "group", name)
	for _, l := range listings {
		l.ID = uuid.New().String()
		l.CreatedAt = time.Now()
		l.IsActive = true
		if l.Type == domain.Request || l.Type == domain.Event {
			l.Deadline = time.Now().Add(30 * 24 * time.Hour)
		}

		if err := repo.Save(ctx, l); err != nil {
			slog.Error("Error saving listing", "title", l.Title, "error", err)
		} else {
			fmt.Printf("Saved: %s\n", l.Title)
		}
	}
}

func getListingData() ([]seedSource, error) {
	data, err := seedData.ReadFile("listings.json")
	if err != nil {
		return nil, err
	}

	var sources []seedSource
	if err := json.Unmarshal(data, &sources); err != nil {
		return nil, err
	}

	return sources, nil
}

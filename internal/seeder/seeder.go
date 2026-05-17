package seeder

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
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
	listings, _, err := repo.FindAll(ctx, "", "", "", 0.0, 0.0, 0.0, "", "", true, 1, 0)
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
		if err := processSeedListing(ctx, repo, l); err != nil {
			slog.Error("Error saving listing", "title", l.Title, "error", err)
		} else {
			fmt.Printf("Saved: %s\n", l.Title)
		}
	}
}

func processSeedListing(ctx context.Context, repo domain.ListingStore, l domain.Listing) error {
	// Use deterministic UUID based on title to keep E2E snapshots stable
	l.ID = uuid.NewMD5(uuid.NameSpaceURL, []byte(l.Title)).String()
	l.CreatedAt = time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	l.IsActive = true
	l.Status = domain.ListingStatusApproved
	if l.Type == domain.Request || l.Type == domain.Event {
		l.Deadline = l.CreatedAt.Add(30 * 24 * time.Hour)
	}

	// Map DFW cities to realistic coordinates
	if l.Latitude == 0 && l.Longitude == 0 {
		if lat, lng := getDFWCoords(l.City); lat != 0 {
			l.Latitude = lat
			l.Longitude = lng
		}
	}

	return repo.Save(ctx, l)
}

func getDFWCoords(city string) (float64, float64) {
	switch strings.ToLower(city) {
	case "dallas":
		return 32.7767, -96.7970
	case "plano":
		return 33.0198, -96.6989
	case "fort worth":
		return 32.7555, -97.3308
	case "arlington":
		return 32.7357, -97.1081
	case "irving":
		return 32.8140, -96.9489
	case "garland":
		return 32.9126, -96.6389
	case "grand prairie":
		return 32.7460, -96.9978
	case "richardson":
		return 32.9482, -96.7297
	case "mckinney":
		return 33.1972, -96.6398
	case "mesquite":
		return 32.7668, -96.5992
	case "frisco":
		return 33.1507, -96.8236
	}
	return 0, 0
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

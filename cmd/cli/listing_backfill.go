package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/spf13/cobra"
)

var listingBackfillCitiesCmd = &cobra.Command{
	Use:   "backfill-cities",
	Short: "Backfill missing city data for listings using geocoding",
	Long: `Iterates through all listings that have an empty city field but have an address,
and uses the Google Geocoding API to attempt to populate the city.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		if cfg.GoogleMapsAPIKey == "" {
			slog.Error("GOOGLE_MAPS_API_KEY is not set. Cannot perform geocoding.")
			os.Exit(1)
		}

		repo := InitRepo()
		geocodingSvc := service.NewGoogleGeocodingService(cfg.GoogleMapsAPIKey)
		ctx := context.Background()

		// Get all listings
		listings, _, err := repo.FindAll(ctx, "", "", "", 0.0, 0.0, 0.0, "", "", false, 0, 0)
		if err != nil {
			slog.Error("Failed to fetch listings", "error", err)
			os.Exit(1)
		}

		updatedCount := 0
		errorCount := 0

		for _, l := range listings {
			changed := false
			if l.City == "" && l.Address != "" {
				slog.Info("Backfilling city for listing", "id", l.ID, "address", l.Address)
				city, err := geocodingSvc.GetCity(ctx, l.Address)
				if err != nil {
					slog.Error("Failed to geocode address city", "id", l.ID, "address", l.Address, "error", err)
					errorCount++
				} else if city != "" {
					l.City = city
					changed = true
					slog.Info("Successfully backfilled city", "id", l.ID, "city", city)
				}
			}
			if (l.Latitude == 0.0 || l.Longitude == 0.0) && l.Address != "" {
				slog.Info("Backfilling coordinates for listing", "id", l.ID, "address", l.Address)
				lat, lng, err := geocodingSvc.Geocode(ctx, l.Address)
				if err != nil {
					slog.Error("Failed to geocode address coordinates", "id", l.ID, "address", l.Address, "error", err)
					errorCount++
				} else if lat != 0.0 && lng != 0.0 {
					l.Latitude = lat
					l.Longitude = lng
					changed = true
					slog.Info("Successfully backfilled coordinates", "id", l.ID, "lat", lat, "lng", lng)
				} else if l.City != "" {
					lat, lng, err = geocodingSvc.Geocode(ctx, l.City)
					if err == nil && lat != 0.0 && lng != 0.0 {
						l.Latitude = lat
						l.Longitude = lng
						changed = true
						slog.Info("Successfully backfilled city-fallback coordinates", "id", l.ID, "lat", lat, "lng", lng)
					}
				}
			}
			if changed {
				if err := repo.Save(ctx, l); err != nil {
					slog.Error("Failed to save backfilled listing", "id", l.ID, "error", err)
					errorCount++
				} else {
					updatedCount++
				}
			}
		}

		fmt.Printf("Backfill complete. Updated: %d, Errors: %d\n", updatedCount, errorCount)
	},
}

func init() {
}

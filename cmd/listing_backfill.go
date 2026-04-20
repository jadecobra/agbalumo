package cmd

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

		repo := initRepo()
		geocodingSvc := service.NewGoogleGeocodingService(cfg.GoogleMapsAPIKey)
		ctx := context.Background()

		// Get all listings
		listings, _, err := repo.FindAll(ctx, "", "", "", "", "", false, 0, 0)
		if err != nil {
			slog.Error("Failed to fetch listings", "error", err)
			os.Exit(1)
		}

		updatedCount := 0
		errorCount := 0

		for _, l := range listings {
			if l.City == "" && l.Address != "" {
				slog.Info("Backfilling city for listing", "id", l.ID, "address", l.Address)
				city, err := geocodingSvc.GetCity(ctx, l.Address)
				if err != nil {
					slog.Error("Failed to geocode address", "id", l.ID, "address", l.Address, "error", err)
					errorCount++
					continue
				}

				if city != "" {
					l.City = city
					if err := repo.Save(ctx, l); err != nil {
						slog.Error("Failed to save listing with updated city", "id", l.ID, "error", err)
						errorCount++
						continue
					}
					updatedCount++
					slog.Info("Successfully backfilled city", "id", l.ID, "city", city)
				} else {
					slog.Info("No city found for listing address", "id", l.ID, "address", l.Address)
				}
			}
		}

		fmt.Printf("Backfill complete. Updated: %d, Errors: %d\n", updatedCount, errorCount)
	},
}
func init() {
	rootCmd.AddCommand(listingBackfillCitiesCmd)
}

package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var locationBackfillCmd = &cobra.Command{
	Use:   "location-backfill",
	Short: "Backfill missing city, state, and country from address using domain heuristics",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = godotenv.Load(".env")
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = ".tester/data/agbalumo.db"
		}
		repo, err := sqlite.NewSQLiteRepository(dbURL)
		if err != nil {
			return err
		}
		defer func() { _ = repo.Close() }()

		ctx := cmd.Context()
		listings, _, err := repo.FindAll(ctx, "", "", "", 0.0, 0.0, 0.0, "", "", true, 10000, 0)
		if err != nil {
			return err
		}

		count := 0
		for _, l := range listings {
			originalCity := l.City
			originalState := l.State
			originalCountry := l.Country

			if l.Address == "" {
				continue
			}

			if l.City == "" {
				l.City = domain.ExtractCityFromAddress(l.Address)
			}
			if l.State == "" {
				l.State = domain.ExtractStateFromAddress(l.Address)
			}
			if l.Country == "" {
				l.Country = domain.ExtractCountryFromAddress(l.Address)
			}

			if l.City != originalCity || l.State != originalState || l.Country != originalCountry {
				err := repo.Save(ctx, l)
				if err != nil {
					fmt.Printf("Error saving listing %s: %v\n", l.ID, err)
					continue
				}
				count++
			}
		}
		fmt.Printf("✅ Success! Backfilled location data for %d listings.\n", count)
		return nil
	},
}

var enrichCmd = &cobra.Command{
	Use:   "enrich",
	Short: "Run the scraper job manually to enrich listings",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load .env if present
		_ = godotenv.Load(".env")
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = ".tester/data/agbalumo.db"
		}
		repo, err := sqlite.NewSQLiteRepository(dbURL)
		if err != nil {
			return err
		}
		defer func() { _ = repo.Close() }()

		scraper := service.NewWebsiteScraper()
		job := service.NewScraperJob(repo, scraper, service.NewGeminiHoursExtractor(os.Getenv("GEMINI_API_KEY"), nil))

		fmt.Println("🚀 Starting Manual Enrichment Job...")
		count, err := job.EnrichListings(cmd.Context(), 50)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Success! Enriched %d listings with sensory signals.\n", count)
		return nil
	},
}

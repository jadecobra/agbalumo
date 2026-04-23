package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	_ "modernc.org/sqlite"
)

const countryUSA = "USA"

func main() {
	fmt.Printf("Backfilling Locations for Listings in: %s\n", domain.DefaultDatabaseURL)

	repo, err := sqlite.NewSQLiteRepository(domain.DefaultDatabaseURL)
	if err != nil {
		log.Fatalf("failed to create repo: %v", err)
	}
	defer func() {
		if closeErr := repo.Close(); closeErr != nil {
			log.Printf("error closing repo: %v", closeErr)
		}
	}()

	ctx := context.Background()

	listings, _, err := repo.FindAll(ctx, "", "", "", 0.0, 0.0, 0.0, "", "", false, 2000, 0)
	if err != nil {
		log.Fatalf("failed to get listings: %v", err)
	}

	fmt.Printf("Processing %d listings...\n", len(listings))
	updated := 0
	for _, l := range listings {
		if l.City == "" {
			continue
		}

		if processListing(ctx, repo, l) {
			updated++
		}
	}

	fmt.Printf("Backfill complete. Records updated: %d\n", updated)
}

func processListing(ctx context.Context, repo domain.ListingRepository, l domain.Listing) bool {
	changed := false

	// Ensure State/Country is present
	if l.Country == "" {
		l.Country = countryUSA
		changed = true
	}

	// Try to extract State from address if empty
	if l.State == "" && l.Address != "" {
		state := extractState(l.Address)
		if state != "" {
			l.State = state
			changed = true
		}
	}

	if changed {
		if err := repo.Save(ctx, l); err != nil {
			log.Printf("Failed to update listing %s: %v", l.ID, err)
			return false
		}
		return true
	}
	return false
}

func extractState(address string) string {
	// Simple comma-based or space-based extraction for state abbreviations
	// TX, TX 75201, Texas (abbreviate to TX if known) etc.
	re := regexp.MustCompile(`(?i)\b(TX|NY|CA|FL|GA|IL|WA|NJ|MD|MA)\b`)
	match := re.FindString(address)
	if match != "" {
		return strings.ToUpper(match)
	}
	return ""
}

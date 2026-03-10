package cmd

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/spf13/cobra"
)

var listingCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new listing",
	Long: `Create a new listing in the agbalumo directory. Mandatory fields include 
the title. Other fields like type, description, and contact information 
can be specified via flags.`,
	Example: `  # Create a basic business listing
  agbalumo listing create --title "Lagos Chop House" --type Service

  # Create a job listing with a deadline
  agbalumo listing create --title "Backend Developer" --type Job --deadline 2026-12-31`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listing := domain.Listing{
			ID:              generateID(),
			OwnerID:         flagOwnerID,
			OwnerOrigin:     flagOrigin,
			Type:            domain.Category(flagType),
			Title:           flagTitle,
			Description:     flagDescription,
			City:            flagCity,
			Address:         flagAddress,
			ContactEmail:    flagEmail,
			ContactPhone:    flagPhone,
			ContactWhatsApp: flagWhatsApp,
			WebsiteURL:      flagWebsite,
			ImageURL:        flagImageURL,
			CreatedAt:       time.Now(),
			IsActive:        true,
			Status:          domain.ListingStatusApproved,
			Skills:          flagSkills,
			JobApplyURL:     flagApplyURL,
			Company:         flagCompany,
			PayRange:        flagPayRange,
		}

		if flagDeadline != "" {
			if t, err := time.Parse("2006-01-02", flagDeadline); err == nil {
				listing.Deadline = t
			} else {
				slog.Warn("Invalid deadline format, expected YYYY-MM-DD", "error", err)
			}
		}
		if flagEventStart != "" {
			if t, err := time.Parse("2006-01-02T15:04", flagEventStart); err == nil {
				listing.EventStart = t
			}
		}
		if flagEventEnd != "" {
			if t, err := time.Parse("2006-01-02T15:04", flagEventEnd); err == nil {
				listing.EventEnd = t
			}
		}
		if flagJobStart != "" {
			if t, err := time.Parse("2006-01-02T15:04", flagJobStart); err == nil {
				listing.JobStartDate = t
			}
		}

		if err := repo.Save(context.Background(), listing); err != nil {
			slog.Error("Failed to create listing", "error", err)
			os.Exit(1)
		}

		if flagJSON {
			data, _ := json.MarshalIndent(listing, "", "  ")
			cmd.Println(string(data))
			return
		}

		cmd.Printf("Listing created successfully: %s\n", listing.ID)
		printListing(cmd, listing)
	},
}

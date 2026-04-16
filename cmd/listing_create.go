package cmd

import (
	"context"
	"encoding/json"
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
			ID:          generateID(),
			OwnerID:     flagOwnerID,
			OwnerOrigin: flagOrigin,
			Type:        domain.Category(flagType),
			CreatedAt:   time.Now(),
			IsActive:    true,
			Status:      domain.ListingStatusApproved,
		}

		applyListingUpdates(&listing)

		exitOnErr(repo.Save(context.Background(), listing), domain.MsgFailedToCreateListing)

		if !flagText {
			data, _ := json.MarshalIndent(listing, "", "  ")
			cmd.Println(string(data))
			return
		}

		cmd.Printf("Listing created successfully: %s\n", listing.ID)
		printListing(cmd, listing)
	},
}

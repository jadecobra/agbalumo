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

var listingUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a listing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		if err != nil {
			slog.Error("Listing not found", "error", err)
			os.Exit(1)
		}

		if flagTitle != "" {
			listing.Title = flagTitle
		}
		if flagDescription != "" {
			listing.Description = flagDescription
		}
		if flagCity != "" {
			listing.City = flagCity
		}
		if flagAddress != "" {
			listing.Address = flagAddress
		}
		if flagEmail != "" {
			listing.ContactEmail = flagEmail
		}
		if flagPhone != "" {
			listing.ContactPhone = flagPhone
		}
		if flagWhatsApp != "" {
			listing.ContactWhatsApp = flagWhatsApp
		}
		if flagWebsite != "" {
			listing.WebsiteURL = flagWebsite
		}
		if flagImageURL != "" {
			listing.ImageURL = flagImageURL
		}
		if flagRemoveImage {
			listing.ImageURL = ""
		}
		if flagDeadline != "" {
			if t, err := time.Parse("2006-01-02", flagDeadline); err == nil {
				listing.Deadline = t
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
		if flagSkills != "" {
			listing.Skills = flagSkills
		}
		if flagJobStart != "" {
			if t, err := time.Parse("2006-01-02T15:04", flagJobStart); err == nil {
				listing.JobStartDate = t
			}
		}
		if flagApplyURL != "" {
			listing.JobApplyURL = flagApplyURL
		}
		if flagCompany != "" {
			listing.Company = flagCompany
		}
		if flagPayRange != "" {
			listing.PayRange = flagPayRange
		}

		if err := repo.Save(context.Background(), listing); err != nil {
			slog.Error(domain.MsgFailedToUpdateListing, "error", err)
			os.Exit(1)
		}

		if !flagText {
			data, _ := json.MarshalIndent(listing, "", "  ")
			cmd.Println(string(data))
			return
		}

		cmd.Printf("Listing updated successfully: %s\n", listing.ID)
		printListing(cmd, listing)
	},
}

package cmd

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/spf13/cobra"
)

// applyDate sets *dst to the parsed date only if the flag is non-empty and parseable.
func applyDate(flag, name string, dst *time.Time) {
	if flag != "" {
		if t := parseDate(flag, name); !t.IsZero() {
			*dst = t
		}
	}
}

// applyDateTime sets *dst to the parsed datetime only if the flag is non-empty and parseable.
func applyDateTime(flag, name string, dst *time.Time) {
	if flag != "" {
		if t := parseDateTime(flag, name); !t.IsZero() {
			*dst = t
		}
	}
}

var listingUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a listing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		exitOnErr(err, "Listing not found")

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
		applyDate(flagDeadline, "deadline", &listing.Deadline)
		applyDateTime(flagEventStart, "event-start", &listing.EventStart)
		applyDateTime(flagEventEnd, "event-end", &listing.EventEnd)
		applyDateTime(flagJobStart, "job-start", &listing.JobStartDate)
		if flagSkills != "" {
			listing.Skills = flagSkills
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

		exitOnErr(repo.Save(context.Background(), listing), domain.MsgFailedToUpdateListing)

		if !flagText {
			data, _ := json.MarshalIndent(listing, "", "  ")
			cmd.Println(string(data))
			return
		}

		cmd.Printf("Listing updated successfully: %s\n", listing.ID)
		printListing(cmd, listing)
	},
}

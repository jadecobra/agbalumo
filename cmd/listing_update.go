package cmd

import (
	"context"
	"encoding/json"

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
		if flagDeadline != "" {
			if t := parseDate(flagDeadline, "deadline"); !t.IsZero() {
				listing.Deadline = t
			}
		}
		if flagEventStart != "" {
			if t := parseDateTime(flagEventStart, "event-start"); !t.IsZero() {
				listing.EventStart = t
			}
		}
		if flagEventEnd != "" {
			if t := parseDateTime(flagEventEnd, "event-end"); !t.IsZero() {
				listing.EventEnd = t
			}
		}
		if flagSkills != "" {
			listing.Skills = flagSkills
		}
		if flagJobStart != "" {
			if t := parseDateTime(flagJobStart, "job-start"); !t.IsZero() {
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

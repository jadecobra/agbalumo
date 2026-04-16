package cmd

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/spf13/cobra"
)

// Listing flags are now defined in shared.go

var listingCmd = &cobra.Command{
	Use:   "listing",
	Short: "Manage listings",
	Long: `The listing command provides subcommands to create, list, retrieve, update, and delete 
listings in the agbalumo directory. Listings represent businesses, services, 
products, jobs, events, or requests within the community.`,
}

func init() {
	listingCmd.AddCommand(listingCreateCmd)
	listingCmd.AddCommand(listingListCmd)
	listingCmd.AddCommand(listingGetCmd)
	listingCmd.AddCommand(listingUpdateCmd)
	listingCmd.AddCommand(listingDeleteCmd)
	listingCmd.AddCommand(listingBackfillCitiesCmd)

	rootCmd.AddCommand(listingCmd)

	bindListingFlags(listingCreateCmd, false)
	bindListingFlags(listingUpdateCmd, true)

	_ = listingCreateCmd.MarkFlagRequired(domain.FieldTitle)
}

func generateID() string {
	return fmt.Sprintf("cli-%s", uuid.New().String()[:8])
}

func printListing(cmd *cobra.Command, l domain.Listing) {
	cmd.Println("==================================")
	cmd.Printf("ID:              %s\n", l.ID)
	cmd.Printf("Title:           %s\n", l.Title)
	cmd.Printf("Type:            %s\n", l.Type)
	cmd.Printf("Origin:          %s\n", l.OwnerOrigin)
	cmd.Printf("Status:          %s\n", l.Status)
	cmd.Printf("Featured:        %v\n", l.Featured)
	cmd.Printf("Description:     %s\n", l.Description)
	cmd.Printf("City:            %s\n", l.City)
	cmd.Printf("Address:         %s\n", l.Address)
	cmd.Printf("Hours:           %s\n", l.HoursOfOperation)
	cmd.Printf("Email:           %s\n", l.ContactEmail)
	cmd.Printf("Phone:           %s\n", l.ContactPhone)
	cmd.Printf("WhatsApp:        %s\n", l.ContactWhatsApp)
	cmd.Printf("Website:         %s\n", l.WebsiteURL)
	cmd.Printf("Image URL:       %s\n", l.ImageURL)
	cmd.Printf("Created:         %s\n", l.CreatedAt.Format(time.RFC3339))
	if !l.Deadline.IsZero() {
		cmd.Printf("Deadline:        %s\n", l.Deadline.Format(domain.DateFormat))
	}
	if !l.EventStart.IsZero() {
		cmd.Printf("Event Start:     %s\n", l.EventStart.Format(time.RFC3339))
	}
	if !l.EventEnd.IsZero() {
		cmd.Printf("Event End:       %s\n", l.EventEnd.Format(time.RFC3339))
	}
	if l.Type == domain.Job {
		cmd.Printf("Company:         %s\n", l.Company)
		cmd.Printf("Skills:          %s\n", l.Skills)
		cmd.Printf("Job Start:       %s\n", l.JobStartDate.Format(time.RFC3339))
		cmd.Printf("Apply URL:       %s\n", l.JobApplyURL)
		cmd.Printf("Pay Range:       %s\n", l.PayRange)
	}
}

func printListingSummary(cmd *cobra.Command, l domain.Listing) {
	cmd.Printf("[%s] %s - %s (%s) [%s]\n", l.ID, l.Title, l.Type, l.City, l.Status)
}

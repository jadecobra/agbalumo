package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/spf13/cobra"
)

var (
	flagTitle       string
	flagType        string
	flagOrigin      string
	flagDescription string
	flagCity        string
	flagAddress     string
	flagEmail       string
	flagPhone       string
	flagWhatsApp    string
	flagWebsite     string
	flagOwnerID     string
	flagImageURL    string
	flagRemoveImage bool
	flagDeadline    string
	flagEventStart  string
	flagEventEnd    string
	flagSkills      string
	flagJobStart    string
	flagApplyURL    string
	flagCompany     string
	flagPayRange    string
)

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

	listingCreateCmd.Flags().StringVarP(&flagTitle, "title", "t", "", "Listing title (required)")
	listingCreateCmd.Flags().StringVarP(&flagType, "type", "y", "Business", "Listing type (Business, Service, Product, Food, Event, Job, Request)")
	listingCreateCmd.Flags().StringVarP(&flagOrigin, "origin", "o", "Nigeria", "Owner origin/country")
	listingCreateCmd.Flags().StringVarP(&flagDescription, "description", "d", "", "Listing description")
	listingCreateCmd.Flags().StringVarP(&flagCity, "city", "c", "", "City")
	listingCreateCmd.Flags().StringVarP(&flagAddress, "address", "a", "", "Address")
	listingCreateCmd.Flags().StringVarP(&flagEmail, "email", "e", "", "Contact email")
	listingCreateCmd.Flags().StringVarP(&flagPhone, "phone", "p", "", "Contact phone")
	listingCreateCmd.Flags().StringVarP(&flagWhatsApp, "whatsapp", "w", "", "WhatsApp number")
	listingCreateCmd.Flags().StringVarP(&flagWebsite, "website", "s", "", "Website URL")
	listingCreateCmd.Flags().StringVarP(&flagImageURL, "image-url", "i", "", "Image URL")
	listingCreateCmd.Flags().StringVarP(&flagOwnerID, "owner-id", "", "", "Owner user ID")
	listingCreateCmd.Flags().StringVar(&flagDeadline, "deadline", "", "Deadline (YYYY-MM-DD)")
	listingCreateCmd.Flags().StringVar(&flagEventStart, "event-start", "", "Event start (YYYY-MM-DDTHH:MM)")
	listingCreateCmd.Flags().StringVar(&flagEventEnd, "event-end", "", "Event end (YYYY-MM-DDTHH:MM)")
	listingCreateCmd.Flags().StringVar(&flagSkills, "skills", "", "Required skills")
	listingCreateCmd.Flags().StringVar(&flagJobStart, "job-start", "", "Job start date (YYYY-MM-DDTHH:MM)")
	listingCreateCmd.Flags().StringVar(&flagApplyURL, "apply-url", "", "Job application URL")
	listingCreateCmd.Flags().StringVar(&flagCompany, "company", "", "Company name")
	listingCreateCmd.Flags().StringVar(&flagPayRange, "pay-range", "", "Pay range")

	listingUpdateCmd.Flags().StringVarP(&flagTitle, "title", "t", "", "New title")
	listingUpdateCmd.Flags().StringVarP(&flagDescription, "description", "d", "", "New description")
	listingUpdateCmd.Flags().StringVarP(&flagCity, "city", "c", "", "New city")
	listingUpdateCmd.Flags().StringVarP(&flagAddress, "address", "a", "", "New address")
	listingUpdateCmd.Flags().StringVarP(&flagEmail, "email", "e", "", "New email")
	listingUpdateCmd.Flags().StringVarP(&flagPhone, "phone", "p", "", "New phone")
	listingUpdateCmd.Flags().StringVarP(&flagWhatsApp, "whatsapp", "w", "", "New WhatsApp")
	listingUpdateCmd.Flags().StringVarP(&flagWebsite, "website", "s", "", "New website")
	listingUpdateCmd.Flags().StringVarP(&flagImageURL, "image-url", "i", "", "New image URL")
	listingUpdateCmd.Flags().BoolVar(&flagRemoveImage, "remove-image", false, "Remove listing image")
	listingUpdateCmd.Flags().StringVar(&flagDeadline, "deadline", "", "New deadline (YYYY-MM-DD)")
	listingUpdateCmd.Flags().StringVar(&flagEventStart, "event-start", "", "New event start")
	listingUpdateCmd.Flags().StringVar(&flagEventEnd, "event-end", "", "New event end")
	listingUpdateCmd.Flags().StringVar(&flagSkills, "skills", "", "New skills")
	listingUpdateCmd.Flags().StringVar(&flagJobStart, "job-start", "", "New job start")
	listingUpdateCmd.Flags().StringVar(&flagApplyURL, "apply-url", "", "New apply URL")
	listingUpdateCmd.Flags().StringVar(&flagCompany, "company", "", "New company")
	listingUpdateCmd.Flags().StringVar(&flagPayRange, "pay-range", "", "New pay range")

	_ = listingCreateCmd.MarkFlagRequired("title")
}

func initRepo() *sqlite.SQLiteRepository {
	dbPath := getDatabaseURL()
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	return repo
}

func getDatabaseURL() string {
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}
	return ".tester/data/agbalumo.db"
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
		cmd.Printf("Deadline:        %s\n", l.Deadline.Format("2006-01-02"))
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

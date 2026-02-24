package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

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
)

var listingCmd = &cobra.Command{
	Use:   "listing",
	Short: "Manage listings",
}

var listingCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new listing",
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
			CreatedAt:       time.Now(),
			IsActive:        true,
			Status:          domain.ListingStatusApproved,
		}

		if err := repo.Save(context.Background(), listing); err != nil {
			slog.Error("Failed to create listing", "error", err)
			os.Exit(1)
		}

		fmt.Printf("Listing created successfully: %s\n", listing.ID)
		printListing(listing)
	},
}

var listingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all listings",
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listings, err := repo.FindAll(context.Background(), "", "", "", "", false, 100, 0)
		if err != nil {
			slog.Error("Failed to list listings", "error", err)
			os.Exit(1)
		}

		if len(listings) == 0 {
			fmt.Println("No listings found")
			return
		}

		fmt.Printf("Found %d listings:\n\n", len(listings))
		for _, l := range listings {
			printListingSummary(l)
		}
	},
}

var listingGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a listing by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		listing, err := repo.FindByID(context.Background(), args[0])
		if err != nil {
			slog.Error("Failed to get listing", "error", err)
			os.Exit(1)
		}

		printListing(listing)
	},
}

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

		if err := repo.Save(context.Background(), listing); err != nil {
			slog.Error("Failed to update listing", "error", err)
			os.Exit(1)
		}

		fmt.Printf("Listing updated successfully: %s\n", listing.ID)
		printListing(listing)
	},
}

var listingDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a listing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := initRepo()

		if err := repo.Delete(context.Background(), args[0]); err != nil {
			slog.Error("Failed to delete listing", "error", err)
			os.Exit(1)
		}

		fmt.Printf("Listing deleted successfully: %s\n", args[0])
	},
}

func init() {
	listingCmd.AddCommand(listingCreateCmd)
	listingCmd.AddCommand(listingListCmd)
	listingCmd.AddCommand(listingGetCmd)
	listingCmd.AddCommand(listingUpdateCmd)
	listingCmd.AddCommand(listingDeleteCmd)

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
	listingCreateCmd.Flags().StringVarP(&flagOwnerID, "owner-id", "", "", "Owner user ID")

	listingUpdateCmd.Flags().StringVarP(&flagTitle, "title", "t", "", "New title")
	listingUpdateCmd.Flags().StringVarP(&flagDescription, "description", "d", "", "New description")
	listingUpdateCmd.Flags().StringVarP(&flagCity, "city", "c", "", "New city")
	listingUpdateCmd.Flags().StringVarP(&flagAddress, "address", "a", "", "New address")
	listingUpdateCmd.Flags().StringVarP(&flagEmail, "email", "e", "", "New email")
	listingUpdateCmd.Flags().StringVarP(&flagPhone, "phone", "p", "", "New phone")
	listingUpdateCmd.Flags().StringVarP(&flagWhatsApp, "whatsapp", "w", "", "New WhatsApp")
	listingUpdateCmd.Flags().StringVarP(&flagWebsite, "website", "s", "", "New website")

	listingCreateCmd.MarkFlagRequired("title")
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
	return "agbalumo.db"
}

func generateID() string {
	return fmt.Sprintf("cli-%d", time.Now().UnixNano())
}

func printListing(l domain.Listing) {
	fmt.Println("==================================")
	fmt.Printf("ID:          %s\n", l.ID)
	fmt.Printf("Title:       %s\n", l.Title)
	fmt.Printf("Type:        %s\n", l.Type)
	fmt.Printf("Origin:      %s\n", l.OwnerOrigin)
	fmt.Printf("Status:      %s\n", l.Status)
	fmt.Printf("Featured:    %v\n", l.Featured)
	fmt.Printf("Description: %s\n", l.Description)
	fmt.Printf("City:        %s\n", l.City)
	fmt.Printf("Address:     %s\n", l.Address)
	fmt.Printf("Email:       %s\n", l.ContactEmail)
	fmt.Printf("Phone:       %s\n", l.ContactPhone)
	fmt.Printf("WhatsApp:    %s\n", l.ContactWhatsApp)
	fmt.Printf("Website:     %s\n", l.WebsiteURL)
	fmt.Printf("Created:     %s\n", l.CreatedAt.Format(time.RFC3339))
	if !l.Deadline.IsZero() {
		fmt.Printf("Deadline:    %s\n", l.Deadline.Format(time.RFC3339))
	}
}

func printListingSummary(l domain.Listing) {
	fmt.Printf("[%s] %s - %s (%s) [%s]\n", l.ID, l.Title, l.Type, l.City, l.Status)
}

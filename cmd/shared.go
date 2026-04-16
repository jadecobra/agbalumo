package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

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

const (
	layoutDate     = "2006-01-02"
	layoutDateTime = "2006-01-02T15:04"

	defaultOrigin = "Nigeria"
	defaultType   = "Business"
)

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

func exitOnErr(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}

// bindListingFlags adds all common listing flags to the given command.
func bindListingFlags(cmd *cobra.Command, isUpdate bool) {
	f := cmd.Flags()
	if isUpdate {
		f.StringVarP(&flagTitle, "title", "t", "", "New title")
		f.StringVarP(&flagDescription, "description", "d", "", "New description")
		f.StringVarP(&flagCity, "city", "c", "", "New city")
		f.StringVarP(&flagAddress, "address", "a", "", "New address")
		f.StringVarP(&flagEmail, "email", "e", "", "New email")
		f.StringVarP(&flagPhone, "phone", "p", "", "New phone")
		f.StringVarP(&flagWhatsApp, "whatsapp", "w", "", "New WhatsApp")
		f.StringVarP(&flagWebsite, "website", "s", "", "New website")
		f.StringVarP(&flagImageURL, "image-url", "i", "", "New image URL")
		f.BoolVar(&flagRemoveImage, "remove-image", false, "Remove listing image")
		f.StringVar(&flagDeadline, "deadline", "", "New deadline (YYYY-MM-DD)")
		f.StringVar(&flagEventStart, "event-start", "", "New event start")
		f.StringVar(&flagEventEnd, "event-end", "", "New event end")
		f.StringVar(&flagSkills, "skills", "", "New skills")
		f.StringVar(&flagJobStart, "job-start", "", "New job start")
		f.StringVar(&flagApplyURL, "apply-url", "", "New apply URL")
		f.StringVar(&flagCompany, "company", "", "New company")
		f.StringVar(&flagPayRange, "pay-range", "", "New pay range")
	} else {
		f.StringVarP(&flagTitle, "title", "t", "", "Listing title (required)")
		f.StringVarP(&flagType, "type", "y", defaultType, "Listing type (Business, Service, Product, Food, Event, Job, Request)")
		f.StringVarP(&flagOrigin, "origin", "o", defaultOrigin, "Owner origin/country")
		f.StringVarP(&flagDescription, "description", "d", "", "Listing description")
		f.StringVarP(&flagCity, "city", "c", "", "City")
		f.StringVarP(&flagAddress, "address", "a", "", "Address")
		f.StringVarP(&flagEmail, "email", "e", "", "Contact email")
		f.StringVarP(&flagPhone, "phone", "p", "", "Contact phone")
		f.StringVarP(&flagWhatsApp, "whatsapp", "w", "", "WhatsApp number")
		f.StringVarP(&flagWebsite, "website", "s", "", "Website URL")
		f.StringVarP(&flagImageURL, "image-url", "i", "", "Image URL")
		f.StringVarP(&flagOwnerID, "owner-id", "", "", "Owner user ID")
		f.StringVar(&flagDeadline, "deadline", "", "Deadline (YYYY-MM-DD)")
		f.StringVar(&flagEventStart, "event-start", "", "Event start (YYYY-MM-DDTHH:MM)")
		f.StringVar(&flagEventEnd, "event-end", "", "Event end (YYYY-MM-DDTHH:MM)")
		f.StringVar(&flagSkills, "skills", "", "Required skills")
		f.StringVar(&flagJobStart, "job-start", "", "Job start date (YYYY-MM-DDTHH:MM)")
		f.StringVar(&flagApplyURL, "apply-url", "", "Job application URL")
		f.StringVar(&flagCompany, "company", "", "Company name")
		f.StringVar(&flagPayRange, "pay-range", "", "Pay range")
	}
}

func parseDate(val string, label string) time.Time {
	if val == "" {
		return time.Time{}
	}
	t, err := time.Parse(layoutDate, val)
	if err != nil {
		slog.Warn(fmt.Sprintf("Invalid %s format, expected YYYY-MM-DD", label), "error", err)
		return time.Time{}
	}
	return t
}

func parseDateTime(val string, label string) time.Time {
	if val == "" {
		return time.Time{}
	}
	t, err := time.Parse(layoutDateTime, val)
	if err != nil {
		slog.Warn(fmt.Sprintf("Invalid %s format, expected YYYY-MM-DDTHH:MM", label), "error", err)
		return time.Time{}
	}
	return t
}

func printListResponse(cmd *cobra.Command, items any, count int, emptyMsg string) bool {
	if count == 0 {
		if !flagText {
			cmd.Println("[]")
		} else {
			cmd.Println(emptyMsg)
		}
		return true
	}

	if !flagText {
		data, _ := json.MarshalIndent(items, "", "  ")
		cmd.Println(string(data))
		return true
	}

	return false
}

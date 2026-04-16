package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/spf13/cobra"

	"github.com/jadecobra/agbalumo/internal/domain"
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
	layoutDate     = domain.DateFormat
	layoutDateTime = domain.DateTimeFormat

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
	if dbURL := os.Getenv(domain.EnvKeyDatabaseURL); dbURL != "" {
		return dbURL
	}
	return domain.DefaultDatabaseURL
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
		f.StringVarP(&flagTitle, domain.FieldTitle, "t", "", "New title")
		f.StringVarP(&flagDescription, domain.FieldDescription, "d", "", "New description")
		f.StringVarP(&flagCity, domain.FieldCity, "c", "", "New city")
		f.StringVarP(&flagAddress, domain.FieldAddress, "a", "", "New address")
		f.StringVarP(&flagEmail, domain.FieldEmail, "e", "", "New email")
		f.StringVarP(&flagPhone, domain.FieldPhone, "p", "", "New phone")
		f.StringVarP(&flagWhatsApp, domain.FieldWhatsApp, "w", "", "New WhatsApp")
		f.StringVarP(&flagWebsite, domain.FieldWebsite, "s", "", "New website")
		f.StringVarP(&flagImageURL, domain.FieldImageURL, "i", "", "New image URL")
		f.BoolVar(&flagRemoveImage, "remove-image", false, "Remove listing image")
		f.StringVar(&flagDeadline, domain.FieldDeadline, "", "New deadline (YYYY-MM-DD)")
		f.StringVar(&flagEventStart, domain.FieldEventStart, "", "New event start")
		f.StringVar(&flagEventEnd, domain.FieldEventEnd, "", "New event end")
		f.StringVar(&flagSkills, domain.FieldSkills, "", "New skills")
		f.StringVar(&flagJobStart, domain.FieldJobStart, "", "New job start")
		f.StringVar(&flagApplyURL, domain.FieldApplyURL, "", "New apply URL")
		f.StringVar(&flagCompany, domain.FieldCompany, "", "New company")
		f.StringVar(&flagPayRange, domain.FieldPayRange, "", "New pay range")
	} else {
		f.StringVarP(&flagTitle, domain.FieldTitle, "t", "", "Listing title (required)")
		f.StringVarP(&flagType, domain.FieldType, "y", defaultType, "Listing type (Business, Service, Product, Food, Event, Job, Request)")
		f.StringVarP(&flagOrigin, "origin", "o", defaultOrigin, "Owner origin/country")
		f.StringVarP(&flagDescription, domain.FieldDescription, "d", "", "Listing description")
		f.StringVarP(&flagCity, domain.FieldCity, "c", "", "City")
		f.StringVarP(&flagAddress, domain.FieldAddress, "a", "", "Address")
		f.StringVarP(&flagEmail, domain.FieldEmail, "e", "", "Contact email")
		f.StringVarP(&flagPhone, domain.FieldPhone, "p", "", "Contact phone")
		f.StringVarP(&flagWhatsApp, domain.FieldWhatsApp, "w", "", "WhatsApp number")
		f.StringVarP(&flagWebsite, domain.FieldWebsite, "s", "", "Website URL")
		f.StringVarP(&flagImageURL, domain.FieldImageURL, "i", "", "Image URL")
		f.StringVarP(&flagOwnerID, "owner-id", "", "", "Owner user ID")
		f.StringVar(&flagDeadline, domain.FieldDeadline, "", "Deadline (YYYY-MM-DD)")
		f.StringVar(&flagEventStart, domain.FieldEventStart, "", "Event start (YYYY-MM-DDTHH:MM)")
		f.StringVar(&flagEventEnd, domain.FieldEventEnd, "", "Event end (YYYY-MM-DDTHH:MM)")
		f.StringVar(&flagSkills, domain.FieldSkills, "", "Required skills")
		f.StringVar(&flagJobStart, domain.FieldJobStart, "", "Job start date (YYYY-MM-DDTHH:MM)")
		f.StringVar(&flagApplyURL, domain.FieldApplyURL, "", "Job application URL")
		f.StringVar(&flagCompany, domain.FieldCompany, "", "Company name")
		f.StringVar(&flagPayRange, domain.FieldPayRange, "", "Pay range")
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

package listing

import (
	"net/http"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

const datetimeLocalFormat = "2006-01-02T15:04"

type ListingFormRequest struct {
	Title             string `form:"title"`
	Type              string `form:"type"`
	OwnerOrigin       string `form:"owner_origin"`
	Description       string `form:"description"`
	City              string `form:"city"`
	Address           string `form:"address"`
	HoursOfOperation  string `form:"hours_of_operation"`
	ContactEmail      string `form:"contact_email"`
	ContactPhone      string `form:"contact_phone"`
	ContactWhatsApp   string `form:"contact_whatsapp"`
	WebsiteURL        string `form:"website_url"`
	DeadlineDate      string `form:"deadline_date"`
	EventStart        string `form:"event_start"`
	EventEnd          string `form:"event_end"`
	Skills            string `form:"skills"`
	JobStartDate      string `form:"job_start_date"`
	JobApplyURL       string `form:"job_apply_url"`
	Company           string `form:"company"`
	PayRange          string `form:"pay_range"`
	RemoveImage       bool   `form:"remove_image"`
	HeatLevel         int    `form:"heat_level"`
	RegionalSpecialty string `form:"regional_specialty"`
	TopDish           string `form:"top_dish"`
}

// ToListing maps the DTO fields directly to the domain Listing and parses dates.
func (req *ListingFormRequest) ToListing(l *domain.Listing) error {
	l.Title = req.Title
	l.Type = domain.Category(req.Type)
	l.OwnerOrigin = req.OwnerOrigin
	l.Description = req.Description
	l.City = req.City
	l.Address = req.Address
	l.HoursOfOperation = req.HoursOfOperation
	l.ContactEmail = req.ContactEmail
	l.ContactPhone = req.ContactPhone
	l.ContactWhatsApp = req.ContactWhatsApp
	l.WebsiteURL = domain.NormalizeURL(req.WebsiteURL)
	l.Skills = req.Skills
	l.JobApplyURL = domain.NormalizeURL(req.JobApplyURL)
	l.Company = req.Company
	l.PayRange = req.PayRange
	l.HeatLevel = req.HeatLevel
	l.RegionalSpecialty = req.RegionalSpecialty
	l.TopDish = req.TopDish

	if err := parseDeadline(req, l); err != nil {
		return err
	}
	if err := parseEventDates(req, l); err != nil {
		return err
	}
	if err := parseJobStartDate(req, l); err != nil {
		return err
	}

	return nil
}

func (h *ListingHandler) bindAndMapListing(c echo.Context, l *domain.Listing) error {
	var req ListingFormRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Request")
	}

	if err := req.ToListing(l); err != nil {
		return err
	}

	if err := h.handleImageUpload(c, l); err != nil {
		return err
	}

	return nil
}

func parseDeadline(req *ListingFormRequest, l *domain.Listing) error {
	if l.Type != domain.Request {
		return nil
	}
	return assignFormDate(req.DeadlineDate, "2006-01-02", "Invalid Date Format", &l.Deadline)
}

func parseEventDates(req *ListingFormRequest, l *domain.Listing) error {
	if l.Type != domain.Event {
		return nil
	}

	if err := assignFormDate(req.EventStart, datetimeLocalFormat, "Invalid Start Date Format", &l.EventStart); err != nil {
		return err
	}
	return assignFormDate(req.EventEnd, datetimeLocalFormat, "Invalid End Date Format", &l.EventEnd)
}

func parseJobStartDate(req *ListingFormRequest, l *domain.Listing) error {
	if l.Type != domain.Job {
		return nil
	}
	return assignFormDate(req.JobStartDate, datetimeLocalFormat, "Invalid Job Start Date Format", &l.JobStartDate)
}

func assignFormDate(val, format, errMsg string, target *time.Time) error {
	parsed, err := parseFormDate(val, format, errMsg)
	if err == nil && !parsed.IsZero() {
		*target = parsed
	}
	return err
}

func parseFormDate(val, format, errMsg string) (time.Time, error) {
	if val == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(format, val)
	if err != nil {
		return time.Time{}, echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}
	return parsed, nil
}

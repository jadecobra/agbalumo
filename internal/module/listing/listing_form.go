package listing

import (
	"net/http"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

type ListingFormRequest struct {
	Title            string `form:"title"`
	Type             string `form:"type"`
	OwnerOrigin      string `form:"owner_origin"`
	Description      string `form:"description"`
	City             string `form:"city"`
	Address          string `form:"address"`
	HoursOfOperation string `form:"hours_of_operation"`
	ContactEmail     string `form:"contact_email"`
	ContactPhone     string `form:"contact_phone"`
	ContactWhatsApp  string `form:"contact_whatsapp"`
	WebsiteURL       string `form:"website_url"`
	DeadlineDate     string `form:"deadline_date"`
	EventStart       string `form:"event_start"`
	EventEnd         string `form:"event_end"`
	Skills           string `form:"skills"`
	JobStartDate     string `form:"job_start_date"`
	JobApplyURL      string `form:"job_apply_url"`
	Company          string `form:"company"`
	PayRange         string `form:"pay_range"`
	RemoveImage      bool   `form:"remove_image"`
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
	l.WebsiteURL = normalizeURL(req.WebsiteURL)
	l.Skills = req.Skills
	l.JobApplyURL = normalizeURL(req.JobApplyURL)
	l.Company = req.Company
	l.PayRange = req.PayRange

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
	if l.Type == domain.Request && req.DeadlineDate != "" {
		parsedTime, err := time.Parse("2006-01-02", req.DeadlineDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Date Format")
		}
		l.Deadline = parsedTime
	}
	return nil
}

func parseEventDates(req *ListingFormRequest, l *domain.Listing) error {
	if l.Type == domain.Event {
		if req.EventStart != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04", req.EventStart)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid Start Date Format")
			}
			l.EventStart = parsedTime
		}
		if req.EventEnd != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04", req.EventEnd)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid End Date Format")
			}
			l.EventEnd = parsedTime
		}
	}
	return nil
}

func parseJobStartDate(req *ListingFormRequest, l *domain.Listing) error {
	if l.Type == domain.Job && req.JobStartDate != "" {
		parsedTime, err := time.Parse("2006-01-02T15:04", req.JobStartDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Job Start Date Format")
		}
		l.JobStartDate = parsedTime
	}
	return nil
}

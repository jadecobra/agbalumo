package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidDeadline = errors.New("request deadline cannot exceed 90 days")
	ErrMissingContact  = errors.New("at least one contact method is required")
	ErrMissingOrigin   = errors.New("owner origin is required")
	ErrInvalidOrigin   = errors.New("owner origin must be an African country")
)

type Location struct {
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Listing represents a directory entry or request.
type Listing struct {
	CreatedAt             time.Time  `json:"created_at" form:"created_at"`
	Deadline              time.Time  `json:"deadline" form:"deadline"`
	EventStart            time.Time  `json:"event_start" form:"event_start"`
	EventEnd              time.Time  `json:"event_end" form:"event_end"`
	JobStartDate          time.Time  `json:"job_start_date" form:"job_start_date"`
	ID                    string     `json:"id" form:"id"`
	OwnerID               string     `json:"owner_id" form:"owner_id"`
	OwnerOrigin           string     `json:"owner_origin" form:"owner_origin"`
	Anchor                string     `json:"anchor" form:"anchor"`
	Title                 string     `json:"title" form:"title"`
	Description           string     `json:"description" form:"description"`
	City                  string     `json:"city" form:"city"`
	State                 string     `json:"state" form:"state"`
	Country               string     `json:"country" form:"country"`
	Address               string     `json:"address" form:"address"`
	HoursOfOperation      string     `json:"hours_of_operation" form:"hours_of_operation"`
	ImageURL              string     `json:"image_url" form:"image_url"`
	ContactEmail          string     `json:"contact_email" form:"contact_email"`
	ContactPhone          string     `json:"contact_phone" form:"contact_phone"`
	ContactWhatsApp       string     `json:"contact_whatsapp" form:"contact_whatsapp"`
	WebsiteURL            string     `json:"website_url" form:"website_url"`
	Skills                string     `json:"skills" form:"skills"`
	JobApplyURL           string     `json:"job_apply_url" form:"job_apply_url"`
	Company               string     `json:"company" form:"company"`
	PayRange              string     `json:"pay_range" form:"pay_range"`
	RegionalSpecialty     string     `json:"regional_specialty" form:"regional_specialty"`
	TopDish               string     `json:"top_dish" form:"top_dish"`
	PaymentMethods        string     `json:"payment_methods" form:"payment_methods"`
	MenuURL               string     `json:"menu_url" form:"menu_url"`
	EnrichmentAttemptedAt *time.Time `json:"enrichment_attempted_at" form:"enrichment_attempted_at"`
	Type                  Category   `json:"type" form:"type"`

	Status    ListingStatus `json:"status" form:"status"`
	HeatLevel int           `json:"heat_level" form:"heat_level"`
	Latitude  float64       `json:"latitude" form:"latitude"`
	Longitude float64       `json:"longitude" form:"longitude"`
	IsActive  bool          `json:"is_active" form:"is_active"`
	Featured  bool          `json:"featured" form:"featured"`
}

// ListingStatus represents the moderation state of a listing.
type ListingStatus string

const (
	ListingStatusPending  ListingStatus = "Pending"
	ListingStatusApproved ListingStatus = "Approved"
	ListingStatusRejected ListingStatus = "Rejected"
)

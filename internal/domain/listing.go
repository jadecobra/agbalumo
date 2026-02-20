package domain

import (
	"errors"
	"time"
)

// Categories
type Category string

const (
	Business Category = "Business"
	Service  Category = "Service"
	Product  Category = "Product"
	Job      Category = "Job"
	Request  Category = "Request"
	Food     Category = "Food"
	Event    Category = "Event"
)

var (
	ErrInvalidDeadline = errors.New("request deadline cannot exceed 90 days")
	ErrMissingContact  = errors.New("at least one contact method is required")
	ErrMissingOrigin   = errors.New("owner origin is required")
	ErrInvalidOrigin   = errors.New("owner origin must be a West African country")
)

// Listing represents a directory entry or request.
type Listing struct {
	ID               string        `json:"id" form:"id"`
	OwnerID          string        `json:"owner_id" form:"owner_id"`         // Link to User.ID
	OwnerOrigin      string        `json:"owner_origin" form:"owner_origin"` // Required: Country of Origin
	Type             Category      `json:"type" form:"type"`
	Anchor           string        `json:"anchor" form:"anchor"` // Food, Professional, etc.
	Title            string        `json:"title" form:"title"`
	Description      string        `json:"description" form:"description"`
	City             string        `json:"city" form:"city"`
	Address          string        `json:"address" form:"address"`                       // New: Specific Business Address
	HoursOfOperation string        `json:"hours_of_operation" form:"hours_of_operation"` // New
	ImageURL         string        `json:"image_url" form:"image_url"`                   // New: Uploaded or Default Image
	ContactEmail     string        `json:"contact_email" form:"contact_email"`
	ContactPhone     string        `json:"contact_phone" form:"contact_phone"` // New: Validation alternative
	ContactWhatsApp  string        `json:"contact_whatsapp" form:"contact_whatsapp"`
	WebsiteURL       string        `json:"website_url" form:"website_url"` // New: Optional
	CreatedAt        time.Time     `json:"created_at" form:"created_at"`
	Deadline         time.Time     `json:"deadline" form:"deadline"` // Required for 'Request'
	EventStart       time.Time     `json:"event_start" form:"event_start"`
	EventEnd         time.Time     `json:"event_end" form:"event_end"`
	Skills           string        `json:"skills" form:"skills"`                 // New: For Job
	JobStartDate     time.Time     `json:"job_start_date" form:"job_start_date"` // New: For Job
	JobApplyURL      string        `json:"job_apply_url" form:"job_apply_url"`   // New: Optional
	Company          string        `json:"company" form:"company"`               // New: For Job
	PayRange         string        `json:"pay_range" form:"pay_range"`           // New: For Job
	IsActive         bool          `json:"is_active" form:"is_active"`
	Status           ListingStatus `json:"status" form:"status"` // New: Moderation Status
}

// ListingStatus represents the moderation state of a listing.
type ListingStatus string

const (
	ListingStatusPending  ListingStatus = "Pending"
	ListingStatusApproved ListingStatus = "Approved"
	ListingStatusRejected ListingStatus = "Rejected"
)

var ValidOrigins = map[string]bool{
	"Nigeria":       true,
	"Ghana":         true,
	"Senegal":       true,
	"Benin":         true,
	"Burkina Faso":  true,
	"Cape Verde":    true,
	"Cote d'Ivoire": true,
	"Gambia":        true,
	"Guinea":        true,
	"Guinea-Bissau": true,
	"Liberia":       true,
	"Mali":          true,
	"Niger":         true,
	"Sierra Leone":  true,
	"Togo":          true,
	"Other":         true,
}

// ClaimableTypes defines which listing categories can be claimed by users.
var ClaimableTypes = map[Category]bool{
	Business: true,
	Service:  true,
	Product:  true,
	Event:    true,
}

// Validate enforces domain rules for the Listing.
func (l *Listing) Validate() error {
	// Origin is required for ALL types
	if l.OwnerOrigin == "" {
		return ErrMissingOrigin
	}

	if !ValidOrigins[l.OwnerOrigin] {
		return ErrInvalidOrigin
	}

	// Address is required for Business and Food
	if l.Type == Business || l.Type == Food {
		if l.Address == "" {
			return errors.New("address is required for business and food listings")
		}
	}

	// Hours of Operation restricted to Business, Service, Food
	if l.HoursOfOperation != "" {
		allowed := l.Type == Business || l.Type == Service || l.Type == Food
		if !allowed {
			return errors.New("hours of operation not applicable for this listing type")
		}
	}

	if l.ContactEmail == "" && l.ContactWhatsApp == "" && l.ContactPhone == "" {
		return ErrMissingContact
	}

	if l.Type == Request {
		if err := l.validateRequest(); err != nil {
			return err
		}
	}

	if l.Type == Event {
		if err := l.validateEvent(); err != nil {
			return err
		}
	}

	if l.Type == Job {
		if err := l.validateJob(); err != nil {
			return err
		}
	}

	return l.validateLengths()
}

func (l *Listing) validateRequest() error {
	// Deadline cannot be in the past (allow for small clock skew/today)
	// Using a 24h buffer for "today" logic as implied by previous handler logic
	if !l.Deadline.IsZero() && l.Deadline.Before(time.Now().Add(-24*time.Hour)) {
		return errors.New("deadline cannot be in the past")
	}

	// Calculate duration between CreatedAt and Deadline
	// If CreatedAt is zero, use Now as a fallback for validation context
	start := l.CreatedAt
	if start.IsZero() {
		start = time.Now()
	}

	limit := start.Add(90 * 24 * time.Hour)
	if l.Deadline.After(limit) {
		return ErrInvalidDeadline
	}
	return nil
}

func (l *Listing) validateEvent() error {
	if l.EventStart.IsZero() {
		return errors.New("event start time is required")
	}
	if l.EventEnd.IsZero() {
		return errors.New("event end time is required")
	}
	if l.EventEnd.Before(l.EventStart) {
		return errors.New("event end time cannot be before start time")
	}
	return nil
}

func (l *Listing) validateJob() error {
	if l.Company == "" {
		return errors.New("company name is required for job listings")
	}
	if l.Description == "" {
		return errors.New("description is required")
	}
	if l.Skills == "" {
		return errors.New("skills are required for job listings")
	}
	if l.PayRange == "" {
		return errors.New("compensation/pay range is required")
	}
	if l.JobStartDate.IsZero() {
		return errors.New("job start date is required")
	}
	// Start date cannot be in the past (allow 24h buffer)
	if l.JobStartDate.Before(time.Now().Add(-24 * time.Hour)) {
		return errors.New("job start date cannot be in the past")
	}
	if l.City == "" && l.Address == "" {
		// We check City primarily as "Location" usually maps to City
		return errors.New("location (city) is required")
	}
	if l.JobApplyURL == "" {
		return errors.New("apply url is required")
	}
	return nil
}

func (l *Listing) validateLengths() error {
	if len(l.Title) > 100 {
		return errors.New("title cannot exceed 100 characters")
	}
	if len(l.Description) > 2000 {
		return errors.New("description cannot exceed 2000 characters")
	}
	if len(l.Company) > 100 {
		return errors.New("company name cannot exceed 100 characters")
	}
	if len(l.Address) > 200 {
		return errors.New("address cannot exceed 200 characters")
	}
	return nil
}

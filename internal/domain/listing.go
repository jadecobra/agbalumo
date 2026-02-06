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
	ID              string    `json:"id" form:"id"`
	OwnerID         string    `json:"owner_id" form:"owner_id"`         // Link to User.ID
	OwnerOrigin     string    `json:"owner_origin" form:"owner_origin"` // Required: Country of Origin
	Type            Category  `json:"type" form:"type"`
	Anchor          string    `json:"anchor" form:"anchor"` // Food, Professional, etc.
	Title           string    `json:"title" form:"title"`
	Description     string    `json:"description" form:"description"`
	City            string    `json:"city" form:"city"`
	Address         string    `json:"address" form:"address"`     // New: Specific Business Address
	ImageURL        string    `json:"image_url" form:"image_url"` // New: Uploaded or Default Image
	ContactEmail    string    `json:"contact_email" form:"contact_email"`
	ContactPhone    string    `json:"contact_phone" form:"contact_phone"` // New: Validation alternative
	ContactWhatsApp string    `json:"contact_whatsapp" form:"contact_whatsapp"`
	WebsiteURL      string    `json:"website_url" form:"website_url"` // New: Optional
	CreatedAt       time.Time `json:"created_at" form:"created_at"`
	Deadline        time.Time `json:"deadline" form:"deadline"` // Required for 'Request'
	IsActive        bool      `json:"is_active" form:"is_active"`
}

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

// Validate enforces domain rules for the Listing.
func (l *Listing) Validate() error {
	// Origin is required for ALL types
	if l.OwnerOrigin == "" {
		return ErrMissingOrigin
	}

	if !ValidOrigins[l.OwnerOrigin] {
		return ErrInvalidOrigin
	}

	if l.ContactEmail == "" && l.ContactWhatsApp == "" && l.ContactPhone == "" {
		return ErrMissingContact
	}

	if l.Type == Request {
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
	}

	return nil
}

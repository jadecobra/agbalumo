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
	Request  Category = "Request"
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
	OwnerOrigin     string    `json:"owner_origin" form:"owner_origin"` // Required: Country of Origin
	Type            Category  `json:"type" form:"type"`
	Anchor          string    `json:"anchor" form:"anchor"` // Food, Professional, etc.
	Title           string    `json:"title" form:"title"`
	Description     string    `json:"description" form:"description"`
	Neighborhood    string    `json:"neighborhood" form:"neighborhood"`
	ImageURL        string    `json:"image_url" form:"image_url"` // New: Uploaded or Default Image
	ContactEmail    string    `json:"contact_email" form:"contact_email"`
	ContactPhone    string    `json:"contact_phone" form:"contact_phone"` // New: Validation alternative
	ContactWhatsApp string    `json:"contact_whatsapp" form:"contact_whatsapp"`
	WebsiteURL      string    `json:"website_url" form:"website_url"` // New: Optional
	CreatedAt       time.Time `json:"created_at" form:"created_at"`
	Deadline        time.Time `json:"deadline" form:"deadline"` // Required for 'Request'
	IsActive        bool      `json:"is_active" form:"is_active"`
}

// Validate enforces domain rules for the Listing.
func (l *Listing) Validate() error {
	if l.OwnerOrigin == "" {
		return ErrMissingOrigin
	}

	// Validate Origin Whitelist (West African Countries)
	validOrigins := map[string]bool{
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

	if !validOrigins[l.OwnerOrigin] {
		return ErrInvalidOrigin
	}

	if l.ContactEmail == "" && l.ContactWhatsApp == "" && l.ContactPhone == "" {
		return ErrMissingContact
	}

	if l.Type == Request {
		// Calculate duration between CreatedAt and Deadline
		// If CreatedAt is zero (e.g. new struct), use Now()?
		// Ideally logic uses l.CreatedAt. If l.CreatedAt is not set, we might assume Now.
		// However, for strict validation, let's assume CreatedAt must be set or compare Deadline to strict 90 days from "now"
		// isn't precise if the object was created previously.
		// Spec says "within 90 days of CreatedAt".

		limit := l.CreatedAt.Add(90 * 24 * time.Hour)
		if l.Deadline.After(limit) {
			return ErrInvalidDeadline
		}
	}

	return nil
}

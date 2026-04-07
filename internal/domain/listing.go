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

type validationRule struct {
	condition func(*Listing) bool
	err       string
}

var validationRules = []validationRule{
	{condition: func(l *Listing) bool { return l.City == "" }, err: "city is required"},
}

var lengthRules = []struct {
	field func(*Listing) int
	err   string
	limit int
}{
	{field: func(l *Listing) int { return len(l.Title) }, limit: 100, err: "title cannot exceed 100 characters"},
	{field: func(l *Listing) int { return len(l.Description) }, limit: 2000, err: "description cannot exceed 2000 characters"},
	{field: func(l *Listing) int { return len(l.Company) }, limit: 100, err: "company name cannot exceed 100 characters"},
	{field: func(l *Listing) int { return len(l.Address) }, limit: 200, err: "address cannot exceed 200 characters"},
}

var jobFields = []struct {
	field func(*Listing) string
	err   string
}{
	{field: func(l *Listing) string { return l.Company }, err: "company name is required for job listings"},
	{field: func(l *Listing) string { return l.Description }, err: "description is required"},
	{field: func(l *Listing) string { return l.Skills }, err: "skills are required for job listings"},
	{field: func(l *Listing) string { return l.PayRange }, err: "compensation/pay range is required"},
	{field: func(l *Listing) string { return l.JobApplyURL }, err: "apply url is required"},
}

// Listing represents a directory entry or request.
type Listing struct {
	CreatedAt        time.Time     `json:"created_at" form:"created_at"`
	Deadline         time.Time     `json:"deadline" form:"deadline"`
	EventStart       time.Time     `json:"event_start" form:"event_start"`
	EventEnd         time.Time     `json:"event_end" form:"event_end"`
	JobStartDate     time.Time     `json:"job_start_date" form:"job_start_date"`
	ID               string        `json:"id" form:"id"`
	OwnerID          string        `json:"owner_id" form:"owner_id"`
	OwnerOrigin      string        `json:"owner_origin" form:"owner_origin"`
	Anchor           string        `json:"anchor" form:"anchor"`
	Title            string        `json:"title" form:"title"`
	Description      string        `json:"description" form:"description"`
	City             string        `json:"city" form:"city"`
	Address          string        `json:"address" form:"address"`
	HoursOfOperation string        `json:"hours_of_operation" form:"hours_of_operation"`
	ImageURL         string        `json:"image_url" form:"image_url"`
	ContactEmail     string        `json:"contact_email" form:"contact_email"`
	ContactPhone     string        `json:"contact_phone" form:"contact_phone"`
	ContactWhatsApp  string        `json:"contact_whatsapp" form:"contact_whatsapp"`
	WebsiteURL       string        `json:"website_url" form:"website_url"`
	Skills           string        `json:"skills" form:"skills"`
	JobApplyURL      string        `json:"job_apply_url" form:"job_apply_url"`
	Company          string        `json:"company" form:"company"`
	PayRange         string        `json:"pay_range" form:"pay_range"`
	Type             Category      `json:"type" form:"type"`
	Status           ListingStatus `json:"status" form:"status"`
	IsActive         bool          `json:"is_active" form:"is_active"`
	Featured         bool          `json:"featured" form:"featured"`
}

// ListingStatus represents the moderation state of a listing.
type ListingStatus string

const (
	ListingStatusPending  ListingStatus = "Pending"
	ListingStatusApproved ListingStatus = "Approved"
	ListingStatusRejected ListingStatus = "Rejected"
)

var ValidOrigins = map[string]bool{
	// West Africa
	"Benin":         true,
	"Burkina Faso":  true,
	"Cabo Verde":    true,
	"Cote d'Ivoire": true,
	"Gambia":        true,
	"Ghana":         true,
	"Guinea":        true,
	"Guinea-Bissau": true,
	"Liberia":       true,
	"Mali":          true,
	"Mauritania":    true,
	"Niger":         true,
	"Nigeria":       true,
	"Senegal":       true,
	"Sierra Leone":  true,
	"Togo":          true,
	// North Africa
	"Algeria":        true,
	"Egypt":          true,
	"Libya":          true,
	"Morocco":        true,
	"Sudan":          true,
	"Tunisia":        true,
	"Western Sahara": true,
	// East Africa
	"Burundi":     true,
	"Comoros":     true,
	"Djibouti":    true,
	"Eritrea":     true,
	"Ethiopia":    true,
	"Kenya":       true,
	"Madagascar":  true,
	"Malawi":      true,
	"Mauritius":   true,
	"Mozambique":  true,
	"Rwanda":      true,
	"Seychelles":  true,
	"Somalia":     true,
	"South Sudan": true,
	"Tanzania":    true,
	"Uganda":      true,
	"Zambia":      true,
	"Zimbabwe":    true,
	// Central Africa
	"Angola":                           true,
	"Cameroon":                         true,
	"Central African Republic":         true,
	"Chad":                             true,
	"Congo":                            true,
	"Democratic Republic of the Congo": true,
	"Equatorial Guinea":                true,
	"Gabon":                            true,
	"Sao Tome and Principe":            true,
	// Southern Africa
	"Botswana":     true,
	"Eswatini":     true,
	"Lesotho":      true,
	"Namibia":      true,
	"South Africa": true,
	// Other
	"Other": true,
}

// Validate enforces domain rules for the Listing.
func (l *Listing) Validate() error {
	if err := l.validateOrigin(); err != nil {
		return err
	}
	if err := l.validateTypeRequirements(); err != nil {
		return err
	}
	if err := l.validateContact(); err != nil {
		return err
	}
	if err := l.applyRules(); err != nil {
		return err
	}
	return l.validateTypeSpecific()
}

// applyRules runs the validationRules, lengthRules, and (for Job) jobFields in sequence.
func (l *Listing) applyRules() error {
	for _, rule := range validationRules {
		if rule.condition(l) {
			return errors.New(rule.err)
		}
	}
	for _, rule := range lengthRules {
		if rule.field(l) > rule.limit {
			return errors.New(rule.err)
		}
	}
	if l.Type != Job {
		return nil
	}
	for _, f := range jobFields {
		if f.field(l) == "" {
			return errors.New(f.err)
		}
	}
	return nil
}

func (l *Listing) validateOrigin() error {
	if l.OwnerOrigin == "" {
		return ErrMissingOrigin
	}
	if !ValidOrigins[l.OwnerOrigin] {
		return ErrInvalidOrigin
	}
	return nil
}

func (l *Listing) validateTypeRequirements() error {
	if (l.Type == Business || l.Type == Food) && l.Address == "" {
		return errors.New("address is required for business and food listings")
	}
	if l.HoursOfOperation != "" && !(l.Type == Business || l.Type == Service || l.Type == Food) {
		return errors.New("hours of operation not applicable for this listing type")
	}
	return nil
}

func (l *Listing) validateContact() error {
	if l.ContactEmail == "" && l.ContactWhatsApp == "" && l.ContactPhone == "" && l.WebsiteURL == "" {
		return ErrMissingContact
	}
	return nil
}

func (l *Listing) validateTypeSpecific() error {
	switch l.Type {
	case Request:
		return l.validateRequest()
	case Event:
		return l.validateEvent()
	case Job:
		return l.validateJob()
	}
	return nil
}

func (l *Listing) validateRequest() error {
	if !l.Deadline.IsZero() && l.Deadline.Before(time.Now().Add(-24*time.Hour)) {
		return errors.New("deadline cannot be in the past")
	}

	start := l.CreatedAt
	if start.IsZero() {
		start = time.Now()
	}

	if l.Deadline.After(start.Add(90 * 24 * time.Hour)) {
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
	if l.JobStartDate.IsZero() {
		return errors.New("job start date is required")
	}
	if l.JobStartDate.Before(time.Now().Add(-24 * time.Hour)) {
		return errors.New("job start date cannot be in the past")
	}
	return nil
}

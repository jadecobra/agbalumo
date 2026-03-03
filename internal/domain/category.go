package domain

import (
	"errors"
	"time"
)

// Category is no longer a strict enum, it's a dynamic string but we keep constants
// for the core types that have special validation logic (Job, Event, Request).
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
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryInactive = errors.New("category is inactive")
)

// CategoryData represents a category entity in the system.
type CategoryData struct {
	ID                        string    `json:"id"`
	Name                      string    `json:"name"`
	Claimable                 bool      `json:"claimable"`
	IsSystem                  bool      `json:"is_system"` // True for core categories (Job, Event, etc.)
	Active                    bool      `json:"active"`
	RequiresSpecialValidation bool      `json:"requires_special_validation"` // True for Job, Event, Request
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// CategoryFilter options for querying categories
type CategoryFilter struct {
	ActiveOnly bool
}

package domain

import (
	"sync"
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

// CategoryData represents a category entity in the system.
type CategoryData struct {
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	ID                        string    `json:"id"`
	Name                      string    `json:"name"`
	Claimable                 bool      `json:"claimable"`
	IsSystem                  bool      `json:"is_system"`
	Active                    bool      `json:"active"`
	RequiresSpecialValidation bool      `json:"requires_special_validation"`
}

// CategoryFilter options for querying categories
type CategoryFilter struct {
	ActiveOnly bool
}

// CategoryCache is a simple thread-safe cache for category data.
type CategoryCache struct {
	expiration time.Time
	categories []CategoryData
	mu         sync.RWMutex
}

func (c *CategoryCache) Get() ([]CategoryData, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if time.Now().After(c.expiration) {
		return nil, false
	}
	return c.categories, true
}

func (c *CategoryCache) Set(categories []CategoryData, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.categories = categories
	c.expiration = time.Now().Add(ttl)
}

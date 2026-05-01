package domain

import "time"

// SavedListing represents a user's saved/favorited listing.
type SavedListing struct {
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id"`
	ListingID string    `json:"listing_id"`
}

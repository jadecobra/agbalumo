package domain

import "time"

// ClaimStatus represents the state of a user's claim request on a listing.
type ClaimStatus string

const (
	ClaimStatusPending  ClaimStatus = "Pending"
	ClaimStatusApproved ClaimStatus = "Approved"
	ClaimStatusRejected ClaimStatus = "Rejected"
)

// ClaimRequest represents a user's request to claim ownership of an admin-created listing.
type ClaimRequest struct {
	ID           string      `json:"id"`
	ListingID    string      `json:"listing_id"`
	ListingTitle string      `json:"listing_title"` // denormalized for admin display
	UserID       string      `json:"user_id"`
	UserName     string      `json:"user_name"`  // denormalized for admin display
	UserEmail    string      `json:"user_email"` // denormalized for admin display
	Status       ClaimStatus `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
}

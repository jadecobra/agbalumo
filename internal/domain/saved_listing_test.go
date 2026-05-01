package domain

import (
	"testing"
	"time"
)

func TestSavedListing(t *testing.T) {
	now := time.Now()
	tests := []struct {
		listing  SavedListing
		expected SavedListing
		name     string
	}{
		{
			name: "initialization",
			listing: SavedListing{
				UserID:    "user-1",
				ListingID: "listing-1",
				CreatedAt: now,
			},
			expected: SavedListing{
				UserID:    "user-1",
				ListingID: "listing-1",
				CreatedAt: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.listing.UserID != tt.expected.UserID {
				t.Errorf("expected UserID %s, got %s", tt.expected.UserID, tt.listing.UserID)
			}
			if tt.listing.ListingID != tt.expected.ListingID {
				t.Errorf("expected ListingID %s, got %s", tt.expected.ListingID, tt.listing.ListingID)
			}
			if !tt.listing.CreatedAt.Equal(tt.expected.CreatedAt) {
				t.Errorf("expected CreatedAt %v, got %v", tt.expected.CreatedAt, tt.listing.CreatedAt)
			}
		})
	}
}

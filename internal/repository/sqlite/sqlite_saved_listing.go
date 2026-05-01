package sqlite

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// SaveListing implements domain.SavedListingStore.
func (r *SQLiteRepository) SaveListing(ctx context.Context, userID, listingID string) error {
	return nil
}

// UnsaveListing implements domain.SavedListingStore.
func (r *SQLiteRepository) UnsaveListing(ctx context.Context, userID, listingID string) error {
	return nil
}

// GetSavedListings implements domain.SavedListingStore.
func (r *SQLiteRepository) GetSavedListings(ctx context.Context, userID string) ([]domain.SavedListing, error) {
	return nil, nil
}

// IsListingSaved implements domain.SavedListingStore.
func (r *SQLiteRepository) IsListingSaved(ctx context.Context, userID, listingID string) (bool, error) {
	return false, nil
}

// GetSavedListingIDs implements domain.SavedListingStore.
func (r *SQLiteRepository) GetSavedListingIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}

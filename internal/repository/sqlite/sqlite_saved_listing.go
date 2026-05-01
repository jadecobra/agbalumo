package sqlite

import (
	"context"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// SaveListing implements domain.SavedListingStore.
func (r *SQLiteRepository) SaveListing(ctx context.Context, userID, listingID string) error {
	query := `INSERT OR IGNORE INTO saved_listings (user_id, listing_id, created_at) VALUES (?, ?, ?)`
	_, err := r.writeDB.ExecContext(ctx, query, userID, listingID, time.Now())
	return err
}

// UnsaveListing implements domain.SavedListingStore.
func (r *SQLiteRepository) UnsaveListing(ctx context.Context, userID, listingID string) error {
	query := `DELETE FROM saved_listings WHERE user_id = ? AND listing_id = ?`
	_, err := r.writeDB.ExecContext(ctx, query, userID, listingID)
	return err
}

// GetSavedListings implements domain.SavedListingStore.
func (r *SQLiteRepository) GetSavedListings(ctx context.Context, userID string) ([]domain.SavedListing, error) {
	query := `SELECT user_id, listing_id, created_at FROM saved_listings WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.readDB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []domain.SavedListing
	for rows.Next() {
		var sl domain.SavedListing
		if err := rows.Scan(&sl.UserID, &sl.ListingID, &sl.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, sl)
	}
	return results, rows.Err()
}

// IsListingSaved implements domain.SavedListingStore.
func (r *SQLiteRepository) IsListingSaved(ctx context.Context, userID, listingID string) (bool, error) {
	query := `SELECT COUNT(*) FROM saved_listings WHERE user_id = ? AND listing_id = ?`
	var count int
	err := r.readDB.QueryRowContext(ctx, query, userID, listingID).Scan(&count)
	return count > 0, err
}

// GetSavedListingIDs implements domain.SavedListingStore.
func (r *SQLiteRepository) GetSavedListingIDs(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT listing_id FROM saved_listings WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.readDB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		results = append(results, id)
	}
	return results, rows.Err()
}

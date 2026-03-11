package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// SaveClaimRequest inserts or updates a claim request.
func (r *SQLiteRepository) SaveClaimRequest(ctx context.Context, req domain.ClaimRequest) error {
	query := `
	INSERT INTO claim_requests (id, listing_id, listing_title, user_id, user_name, user_email, status, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		status = excluded.status;
	`
	_, err := r.db.ExecContext(ctx, query,
		req.ID, req.ListingID, req.ListingTitle, req.UserID, req.UserName, req.UserEmail, req.Status, req.CreatedAt,
	)
	return err
}

// GetPendingClaimRequests returns all claim requests with status=Pending.
func (r *SQLiteRepository) GetPendingClaimRequests(ctx context.Context) ([]domain.ClaimRequest, error) {
	query := `
		SELECT id, listing_id, COALESCE(listing_title,''), user_id, COALESCE(user_name,''), COALESCE(user_email,''), status, created_at
		FROM claim_requests
		WHERE status = 'Pending'
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []domain.ClaimRequest
	for rows.Next() {
		var cr domain.ClaimRequest
		if err := rows.Scan(&cr.ID, &cr.ListingID, &cr.ListingTitle, &cr.UserID, &cr.UserName, &cr.UserEmail, &cr.Status, &cr.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, cr)
	}
	return results, rows.Err()
}

// UpdateClaimRequestStatus updates a claim request's status.
func (r *SQLiteRepository) UpdateClaimRequestStatus(ctx context.Context, id string, status domain.ClaimStatus) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `UPDATE claim_requests SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("claim request not found")
	}

	if status == domain.ClaimStatusApproved {
		_, err = tx.ExecContext(ctx, `
			UPDATE listings SET owner_id = (
				SELECT user_id FROM claim_requests WHERE id = ?
			)
			WHERE id = (
				SELECT listing_id FROM claim_requests WHERE id = ?
			)`, id, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetClaimRequestByUserAndListing retrieves any existing claim request for a user/listing pair.
func (r *SQLiteRepository) GetClaimRequestByUserAndListing(ctx context.Context, userID, listingID string) (domain.ClaimRequest, error) {
	query := `
		SELECT id, listing_id, COALESCE(listing_title,''), user_id, COALESCE(user_name,''), COALESCE(user_email,''), status, created_at
		FROM claim_requests
		WHERE user_id = ? AND listing_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, userID, listingID)
	var cr domain.ClaimRequest
	err := row.Scan(&cr.ID, &cr.ListingID, &cr.ListingTitle, &cr.UserID, &cr.UserName, &cr.UserEmail, &cr.Status, &cr.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.ClaimRequest{}, errors.New("claim request not found")
	}
	return cr, err
}

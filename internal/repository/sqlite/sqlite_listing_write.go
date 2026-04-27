package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// Save inserts or updates a listing.
func (r *SQLiteRepository) Save(ctx context.Context, l domain.Listing) error {
	query := ListingUpsertSQL

	_, err := r.writeDB.ExecContext(ctx, query, r.listingArgs(l)...)
	return err
}

// SaveBatch inserts or updates multiple listings in a single transaction.
func (r *SQLiteRepository) SaveBatch(ctx context.Context, listings []domain.Listing) error {
	tx, err := r.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query := ListingUpsertSQL

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for _, l := range listings {

		_, err := stmt.ExecContext(ctx, r.listingArgs(l)...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// BulkInsertListings executes bulk INSERT statements chunked into batches of 500.
func (r *SQLiteRepository) BulkInsertListings(ctx context.Context, listings []domain.Listing) error {
	tx, err := r.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	const batchSize = 500
	for i := 0; i < len(listings); i += batchSize {
		end := i + batchSize
		if end > len(listings) {
			end = len(listings)
		}

		if err := r.insertBatch(ctx, tx, listings[i:end]); err != nil {
			return err
		}
		r.logBulkProgress(end, len(listings))
	}

	return tx.Commit()
}

func (r *SQLiteRepository) insertBatch(ctx context.Context, tx *sql.Tx, batch []domain.Listing) error {
	query, args := r.buildBulkInsertSQL(batch)
	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

func (r *SQLiteRepository) buildBulkInsertSQL(batch []domain.Listing) (string, []interface{}) {
	const numFields = 36
	const placeholders = "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	var sb strings.Builder
	// Pre-allocate approximate size: len(batch) * len(placeholders) + SQL header/footer
	sb.Grow(len(batch)*(len(placeholders)+2) + 512)

	sb.WriteString(`INSERT INTO listings ` + listingColumns + ` VALUES `)

	args := make([]interface{}, len(batch)*numFields)

	for i, l := range batch {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(placeholders)
		r.fillListingArgs(args, i*numFields, l)
	}

	sb.WriteString(` ` + listingUpsertUpdate + `;`)

	return sb.String(), args
}

func (r *SQLiteRepository) listingArgs(l domain.Listing) []interface{} {
	args := make([]interface{}, 36)
	r.fillListingArgs(args, 0, l)
	return args
}

func (r *SQLiteRepository) fillListingArgs(args []interface{}, offset int, l domain.Listing) {
	args[offset+0] = l.ID
	args[offset+1] = l.OwnerID
	args[offset+2] = l.Title
	args[offset+3] = l.Description
	args[offset+4] = l.Type
	args[offset+5] = l.OwnerOrigin
	args[offset+6] = l.City
	args[offset+7] = l.State
	args[offset+8] = l.Country
	args[offset+9] = l.Address
	args[offset+10] = l.HoursOfOperation
	args[offset+11] = l.IsActive
	args[offset+12] = l.CreatedAt
	args[offset+13] = l.ImageURL
	args[offset+14] = l.ContactEmail
	args[offset+15] = l.ContactPhone
	args[offset+16] = l.ContactWhatsApp
	args[offset+17] = l.WebsiteURL
	args[offset+18] = l.Deadline
	args[offset+19] = l.EventStart
	args[offset+20] = l.EventEnd
	args[offset+21] = l.Skills
	args[offset+22] = l.JobStartDate
	args[offset+23] = l.JobApplyURL
	args[offset+24] = l.Company
	args[offset+25] = l.PayRange
	args[offset+26] = r.ensureStatus(l.Status)
	args[offset+27] = l.Featured
	args[offset+28] = l.HeatLevel
	args[offset+29] = l.RegionalSpecialty
	args[offset+30] = l.TopDish
	args[offset+31] = l.PaymentMethods
	args[offset+32] = l.MenuURL
	args[offset+33] = l.Latitude
	args[offset+34] = l.Longitude
	args[offset+35] = l.EnrichmentAttemptedAt
}

func (r *SQLiteRepository) ensureStatus(s domain.ListingStatus) string {
	if s == "" {
		return string(domain.ListingStatusApproved)
	}
	return string(s)
}

func (r *SQLiteRepository) logBulkProgress(current, total int) {
	if total <= 0 {
		return
	}
	percentage := (current * 100) / total
	// Log progress at major milestones to keep logs clean
	if percentage%25 == 0 || current == total {
		slog.Info("Bulk insert progress", slog.Int("percentage", percentage), slog.Int("processed", current), slog.Int("total", total))
	}
}

func (r *SQLiteRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM listings WHERE id = ?`
	result, err := r.writeDB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrListingNotFound
	}
	return nil
}

func (r *SQLiteRepository) ExpireListings(ctx context.Context) (int64, error) {
	now := time.Now().UTC()
	var totalAffected int64
	batchSize := 100

	query := `
		UPDATE listings 
		SET is_active = false 
		WHERE rowid IN (
			SELECT rowid FROM listings
			WHERE is_active = true 
			AND (
				(type = 'Request' AND deadline < ?) 
				OR 
				(type = 'Event' AND event_end < ?)
				OR
				(type = 'Job' AND job_start_date < ?)
			)
			LIMIT ?
		)
	`

	for {
		result, err := r.writeDB.ExecContext(ctx, query, now, now, now.AddDate(0, 0, -90), batchSize)
		if err != nil {
			return totalAffected, err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return totalAffected, err
		}

		totalAffected += affected
		if affected < int64(batchSize) {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	return totalAffected, nil
}

// SetFeatured toggles the featured status of a listing.
func (r *SQLiteRepository) SetFeatured(ctx context.Context, id string, featured bool) error {
	query := `UPDATE listings SET featured = ? WHERE id = ?`
	_, err := r.writeDB.ExecContext(ctx, query, featured, id)
	return err
}

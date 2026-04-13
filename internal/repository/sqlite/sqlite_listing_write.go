package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// Save inserts or updates a listing.
func (r *SQLiteRepository) Save(ctx context.Context, l domain.Listing) error {
	query := ListingUpsertSQL


	_, err := r.writeDB.ExecContext(ctx, query,
		l.ID, l.OwnerID, l.Title, l.Description, l.Type, l.OwnerOrigin, l.City, l.Address, l.HoursOfOperation, l.IsActive, l.CreatedAt,
		l.ImageURL, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.Deadline, l.EventStart, l.EventEnd,
		l.Skills, l.JobStartDate, l.JobApplyURL, l.Company, l.PayRange, r.ensureStatus(l.Status), l.Featured,
		l.HeatLevel, l.RegionalSpecialty, l.TopDish,
	)
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

		_, err := stmt.ExecContext(ctx,
			l.ID, l.OwnerID, l.Title, l.Description, l.Type, l.OwnerOrigin, l.City, l.Address, l.HoursOfOperation, l.IsActive, l.CreatedAt,
			l.ImageURL, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.Deadline, l.EventStart, l.EventEnd,
			l.Skills, l.JobStartDate, l.JobApplyURL, l.Company, l.PayRange, r.ensureStatus(l.Status), l.Featured,
			l.HeatLevel, l.RegionalSpecialty, l.TopDish,
		)
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
	query := `INSERT INTO listings ` + listingColumns + ` VALUES `
	args := make([]interface{}, 0, len(batch)*29)

	for i, l := range batch {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		args = append(args,
			l.ID, l.OwnerID, l.Title, l.Description, l.Type, l.OwnerOrigin, l.City, l.Address, l.HoursOfOperation, l.IsActive, l.CreatedAt,
			l.ImageURL, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.Deadline, l.EventStart, l.EventEnd,
			l.Skills, l.JobStartDate, l.JobApplyURL, l.Company, l.PayRange, r.ensureStatus(l.Status), l.Featured,
			l.HeatLevel, l.RegionalSpecialty, l.TopDish,
		)
	}

	query += ` ` + listingUpsertUpdate + `;`

	return query, args
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

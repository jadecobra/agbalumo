package sqlite

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// Save inserts or updates a listing.
func (r *SQLiteRepository) Save(ctx context.Context, l domain.Listing) error {
	query := `
	INSERT INTO listings (id, owner_id, title, description, type, owner_origin, city, address, hours_of_operation, is_active, created_at, image_url, contact_email, contact_phone, contact_whatsapp, website_url, deadline, event_start, event_end, skills, job_start_date, job_apply_url, company, pay_range, status, featured)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		owner_id = excluded.owner_id,
		title = excluded.title,
		description = excluded.description,
		type = excluded.type,
		owner_origin = excluded.owner_origin,
		city = excluded.city,
		address = excluded.address,
		hours_of_operation = excluded.hours_of_operation,
		is_active = excluded.is_active,
		image_url = excluded.image_url,
		contact_email = excluded.contact_email,
		contact_phone = excluded.contact_phone,
		contact_whatsapp = excluded.contact_whatsapp,
		website_url = excluded.website_url,
		deadline = excluded.deadline,
		event_start = excluded.event_start,
		event_end = excluded.event_end,
		skills = excluded.skills,
		job_start_date = excluded.job_start_date,
		job_apply_url = excluded.job_apply_url,
		company = excluded.company,
		pay_range = excluded.pay_range,
		status = excluded.status,
		featured = excluded.featured;
	`

	status := string(l.Status)
	if status == "" {
		status = string(domain.ListingStatusApproved)
	}

	_, err := r.writeDB.ExecContext(ctx, query,
		l.ID, l.OwnerID, l.Title, l.Description, l.Type, l.OwnerOrigin, l.City, l.Address, l.HoursOfOperation, l.IsActive, l.CreatedAt,
		l.ImageURL, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.Deadline, l.EventStart, l.EventEnd,
		l.Skills, l.JobStartDate, l.JobApplyURL, l.Company, l.PayRange, status, l.Featured,
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

	query := `
	INSERT INTO listings (id, owner_id, title, description, type, owner_origin, city, address, hours_of_operation, is_active, created_at, image_url, contact_email, contact_phone, contact_whatsapp, website_url, deadline, event_start, event_end, skills, job_start_date, job_apply_url, company, pay_range, status, featured)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		owner_id = excluded.owner_id,
		title = excluded.title,
		description = excluded.description,
		type = excluded.type,
		owner_origin = excluded.owner_origin,
		city = excluded.city,
		address = excluded.address,
		hours_of_operation = excluded.hours_of_operation,
		is_active = excluded.is_active,
		image_url = excluded.image_url,
		contact_email = excluded.contact_email,
		contact_phone = excluded.contact_phone,
		contact_whatsapp = excluded.contact_whatsapp,
		website_url = excluded.website_url,
		deadline = excluded.deadline,
		event_start = excluded.event_start,
		event_end = excluded.event_end,
		skills = excluded.skills,
		job_start_date = excluded.job_start_date,
		job_apply_url = excluded.job_apply_url,
		company = excluded.company,
		pay_range = excluded.pay_range,
		status = excluded.status,
		featured = excluded.featured;
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for _, l := range listings {
		status := string(l.Status)
		if status == "" {
			status = string(domain.ListingStatusApproved)
		}

		_, err := stmt.ExecContext(ctx,
			l.ID, l.OwnerID, l.Title, l.Description, l.Type, l.OwnerOrigin, l.City, l.Address, l.HoursOfOperation, l.IsActive, l.CreatedAt,
			l.ImageURL, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.Deadline, l.EventStart, l.EventEnd,
			l.Skills, l.JobStartDate, l.JobApplyURL, l.Company, l.PayRange, status, l.Featured,
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

	total := len(listings)
	nextThreshold := 10

	const batchSize = 500
	for i := 0; i < len(listings); i += batchSize {
		end := i + batchSize
		if end > len(listings) {
			end = len(listings)
		}
		batch := listings[i:end]

		query := `
		INSERT INTO listings (id, owner_id, title, description, type, owner_origin, city, address, hours_of_operation, is_active, created_at, image_url, contact_email, contact_phone, contact_whatsapp, website_url, deadline, event_start, event_end, skills, job_start_date, job_apply_url, company, pay_range, status, featured)
		VALUES `
		var args []interface{}
		for j, l := range batch {
			if j > 0 {
				query += ", "
			}
			query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

			status := string(l.Status)
			if status == "" {
				status = string(domain.ListingStatusApproved)
			}

			args = append(args,
				l.ID, l.OwnerID, l.Title, l.Description, l.Type, l.OwnerOrigin, l.City, l.Address, l.HoursOfOperation, l.IsActive, l.CreatedAt,
				l.ImageURL, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.Deadline, l.EventStart, l.EventEnd,
				l.Skills, l.JobStartDate, l.JobApplyURL, l.Company, l.PayRange, status, l.Featured,
			)
		}

		query += `
		ON CONFLICT(id) DO UPDATE SET
			owner_id = excluded.owner_id,
			title = excluded.title,
			description = excluded.description,
			type = excluded.type,
			owner_origin = excluded.owner_origin,
			city = excluded.city,
			address = excluded.address,
			hours_of_operation = excluded.hours_of_operation,
			is_active = excluded.is_active,
			image_url = excluded.image_url,
			contact_email = excluded.contact_email,
			contact_phone = excluded.contact_phone,
			contact_whatsapp = excluded.contact_whatsapp,
			website_url = excluded.website_url,
			deadline = excluded.deadline,
			event_start = excluded.event_start,
			event_end = excluded.event_end,
			skills = excluded.skills,
			job_start_date = excluded.job_start_date,
			job_apply_url = excluded.job_apply_url,
			company = excluded.company,
			pay_range = excluded.pay_range,
			status = excluded.status,
			featured = excluded.featured;
		`

		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		if total > 0 {
			percentage := (end * 100) / total
			if percentage >= nextThreshold {
				slog.Info("Bulk insert progress", slog.Int("percentage", percentage), slog.Int("processed", end), slog.Int("total", total))
				nextThreshold = ((percentage / 10) + 1) * 10
			}
		}
	}

	return tx.Commit()
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
		return errors.New("listing not found")
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

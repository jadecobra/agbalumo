package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

const listingSelections = `
	id, COALESCE(owner_id, ''), owner_origin, type, title, description,
	COALESCE(city, ''), COALESCE(address, ''), COALESCE(hours_of_operation, ''), 
	COALESCE(contact_email, ''), COALESCE(contact_phone, ''), COALESCE(contact_whatsapp, ''),
	COALESCE(website_url, ''), COALESCE(image_url, ''), created_at, deadline, is_active,
	event_start, event_end,
	COALESCE(skills, ''), job_start_date, COALESCE(job_apply_url, ''),
	COALESCE(company, ''), COALESCE(pay_range, ''), COALESCE(status, 'Approved'), featured
`

func scanListing(s Scanner) (domain.Listing, error) {
	var l domain.Listing
	var deadline, eventStart, eventEnd, jobStart sql.NullTime

	err := s.Scan(
		&l.ID, &l.OwnerID, &l.OwnerOrigin, &l.Type, &l.Title, &l.Description,
		&l.City, &l.Address, &l.HoursOfOperation, &l.ContactEmail, &l.ContactPhone, &l.ContactWhatsApp,
		&l.WebsiteURL, &l.ImageURL, &l.CreatedAt, &deadline, &l.IsActive,
		&eventStart, &eventEnd,
		&l.Skills, &jobStart, &l.JobApplyURL,
		&l.Company, &l.PayRange, &l.Status, &l.Featured,
	)
	if err != nil {
		return domain.Listing{}, err
	}

	if deadline.Valid {
		l.Deadline = deadline.Time
	}
	if eventStart.Valid {
		l.EventStart = eventStart.Time
	}
	if eventEnd.Valid {
		l.EventEnd = eventEnd.Time
	}
	if jobStart.Valid {
		l.JobStartDate = jobStart.Time
	}
	return l, nil
}

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

	_, err := r.db.ExecContext(ctx, query,
		l.ID, l.OwnerID, l.Title, l.Description, l.Type, l.OwnerOrigin, l.City, l.Address, l.HoursOfOperation, l.IsActive, l.CreatedAt,
		l.ImageURL, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.Deadline, l.EventStart, l.EventEnd,
		l.Skills, l.JobStartDate, l.JobApplyURL, l.Company, l.PayRange, status, l.Featured,
	)
	return err
}

// SaveBatch inserts or updates multiple listings in a single transaction.
func (r *SQLiteRepository) SaveBatch(ctx context.Context, listings []domain.Listing) error {
	tx, err := r.db.BeginTx(ctx, nil)
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
	tx, err := r.db.BeginTx(ctx, nil)
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
	}

	return tx.Commit()
}

func (r *SQLiteRepository) FindAll(ctx context.Context, filterType string, queryText string, sortField string, sortOrder string, includeInactive bool, limit int, offset int) ([]domain.Listing, int, error) {
	start := time.Now()
	defer func() {
		if duration := time.Since(start); duration > 50*time.Millisecond {
			slog.Info("Slow query detected", slog.String("query", "FindAll"), slog.Int64("duration_ms", duration.Milliseconds()))
		}
	}()

	whereClause := " WHERE 1=1"
	var args []interface{}

	if !includeInactive {
		whereClause += ` AND is_active = true AND status = 'Approved'`
	}

	if filterType != "" {
		whereClause += ` AND type = ?`
		args = append(args, filterType)
	}

	if queryText != "" {
		whereClause += ` AND rowid IN (SELECT rowid FROM listings_fts WHERE listings_fts MATCH ?)`
		args = append(args, queryText)
	}

	// Get total count first
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM listings` + whereClause
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	orderClause := "created_at DESC"
	if sortField != "" {
		field := "created_at"
		switch sortField {
		case "title":
			field = "title"
		case "status":
			field = "status"
		case "featured":
			field = "featured"
		case "type":
			field = "type"
		}

		order := "DESC"
		if sortOrder == "ASC" || sortOrder == "asc" {
			order = "ASC"
		}
		orderClause = field + " " + order
	}

	query := `SELECT ` + listingSelections + ` FROM listings 
	          WHERE rowid IN (SELECT rowid FROM listings` + whereClause + ` ORDER BY ` + orderClause + ` LIMIT ? OFFSET ?)
	          ORDER BY ` + orderClause
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var listings []domain.Listing
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, 0, err
		}
		listings = append(listings, l)
	}
	return listings, totalCount, rows.Err()
}

func (r *SQLiteRepository) FindByID(ctx context.Context, id string) (domain.Listing, error) {
	start := time.Now()
	defer func() {
		if duration := time.Since(start); duration > 50*time.Millisecond {
			slog.Info("Slow query detected", slog.String("query", "FindByID"), slog.Int64("duration_ms", duration.Milliseconds()))
		}
	}()

	query := `
		SELECT ` + listingSelections + `
		FROM listings
		WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	l, err := scanListing(row)
	if err == sql.ErrNoRows {
		return domain.Listing{}, errors.New("listing not found")
	}
	return l, err
}

func (r *SQLiteRepository) FindByTitle(ctx context.Context, title string) ([]domain.Listing, error) {
	query := `
		SELECT ` + listingSelections + `
		FROM listings
		WHERE title = ?
	`
	rows, err := r.db.QueryContext(ctx, query, title)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var listings []domain.Listing
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, err
		}
		listings = append(listings, l)
	}
	return listings, rows.Err()
}

func (r *SQLiteRepository) FindAllByOwner(ctx context.Context, ownerID string, limit int, offset int) ([]domain.Listing, int, error) {
	var totalCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM listings WHERE owner_id = ?", ownerID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT ` + listingSelections + `
              FROM listings 
              WHERE owner_id = ? 
              ORDER BY created_at DESC
              LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var listings []domain.Listing
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, 0, err
		}
		listings = append(listings, l)
	}
	return listings, totalCount, rows.Err()
}

// TitleExists checks if a listing with the given title exists using an efficient EXISTS query.
func (r *SQLiteRepository) TitleExists(ctx context.Context, title string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM listings WHERE title = ?)`
	err := r.db.QueryRowContext(ctx, query, title).Scan(&exists)
	return exists, err
}

func (r *SQLiteRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM listings WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *SQLiteRepository) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	start := time.Now()
	defer func() {
		if duration := time.Since(start); duration > 50*time.Millisecond {
			slog.Info("Slow query detected", slog.String("query", "GetCounts"), slog.Int64("duration_ms", duration.Milliseconds()))
		}
	}()

	query := `SELECT type, COUNT(*) FROM listings WHERE is_active = true AND status = 'Approved' GROUP BY type`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	counts := make(map[domain.Category]int)
	for rows.Next() {
		var cat domain.Category
		var count int
		if err := rows.Scan(&cat, &count); err != nil {
			return nil, err
		}
		counts[cat] = count
	}
	return counts, rows.Err()
}

func (r *SQLiteRepository) GetLocations(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT city FROM listings WHERE is_active = true AND status = 'Approved' AND city != '' ORDER BY city ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var locations []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			return nil, err
		}
		locations = append(locations, city)
	}
	return locations, rows.Err()
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
		result, err := r.db.ExecContext(ctx, query, now, now, now.AddDate(0, 0, -90), batchSize)
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

// GetFeaturedListings returns featured listings set by admin.
func (r *SQLiteRepository) GetFeaturedListings(ctx context.Context) ([]domain.Listing, error) {
	query := `
		SELECT ` + listingSelections + `
		FROM listings 
		WHERE featured = 1 
		AND is_active = 1 
		ORDER BY created_at DESC 
		LIMIT 5
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var listings []domain.Listing
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, err
		}
		listings = append(listings, l)
	}
	return listings, rows.Err()
}

// SetFeatured toggles the featured status of a listing.
func (r *SQLiteRepository) SetFeatured(ctx context.Context, id string, featured bool) error {
	query := `UPDATE listings SET featured = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, featured, id)
	return err
}

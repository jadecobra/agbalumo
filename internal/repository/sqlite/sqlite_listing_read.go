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

func (r *SQLiteRepository) FindAll(ctx context.Context, filterType string, queryText string, sortField string, sortOrder string, includeInactive bool, limit int, offset int) ([]domain.Listing, int, error) {
	start := time.Now()
	defer func() {
		if duration := time.Since(start); duration > r.slowQueryThreshold {
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
		if duration := time.Since(start); duration > r.slowQueryThreshold {
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

func (r *SQLiteRepository) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	start := time.Now()
	defer func() {
		if duration := time.Since(start); duration > r.slowQueryThreshold {
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

// GetFeaturedListings returns featured listings set by admin, optionally filtered by category.
func (r *SQLiteRepository) GetFeaturedListings(ctx context.Context, category string) ([]domain.Listing, error) {
	whereClause := "WHERE featured = 1 AND is_active = 1"
	var args []interface{}

	if category != "" {
		whereClause += " AND type = ?"
		args = append(args, category)
	}

	query := `
		SELECT ` + listingSelections + `
		FROM listings 
		` + whereClause + `
		ORDER BY created_at DESC 
		LIMIT 3
	`
	rows, err := r.db.QueryContext(ctx, query, args...)
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

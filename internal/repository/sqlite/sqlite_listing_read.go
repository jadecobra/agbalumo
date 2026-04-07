package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

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
	defer r.logSlowQuery("FindAll", start)

	where, args := r.buildFindAllWhere(filterType, queryText, includeInactive)

	totalCount, err := r.getCount(ctx, "listings", where, args)
	if err != nil {
		return nil, 0, err
	}

	order := r.buildOrderClause(sortField, sortOrder)

	// #nosec G202 - Dynamic query construction with trusted internal fragments
	query := `SELECT ` + ListingSelectionsSQL + ` FROM listings 
	          WHERE rowid IN (SELECT rowid FROM listings ` + where + ` ORDER BY ` + order + ` LIMIT ? OFFSET ?)
	          ORDER BY ` + order
	args = append(args, limit, offset)

	rows, err := r.readDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	listings, err := scanListings(rows)
	if err != nil {
		return nil, 0, err
	}
	return listings, totalCount, nil
}

func (r *SQLiteRepository) buildFindAllWhere(filterType, queryText string, includeInactive bool) (string, []interface{}) {
	where := " WHERE 1=1"
	var args []interface{}

	if !includeInactive {
		where += ` AND ` + ListingActiveApprovedSQL
	}

	if filterType != "" {
		where += ListingFilterTypeSQL
		args = append(args, filterType)
	}

	if queryText != "" {
		where += ` AND rowid IN (SELECT rowid FROM listings_fts WHERE listings_fts MATCH ?)`
		args = append(args, queryText)
	}

	return where, args
}

func (r *SQLiteRepository) buildOrderClause(sortField, sortOrder string) string {
	if sortField == "" {
		return "featured DESC, created_at DESC"
	}

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
	if strings.ToLower(sortOrder) == "asc" {
		order = "ASC"
	}

	if field == "featured" {
		return "featured " + order + ", created_at DESC"
	}
	return "featured DESC, " + field + " " + order
}

func (r *SQLiteRepository) getCount(ctx context.Context, table, where string, args []interface{}) (int, error) {
	var count int
	// #nosec G202 - Dynamic query construction with trusted internal fragments
	query := "SELECT COUNT(*) FROM " + table + where
	err := r.readDB.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *SQLiteRepository) logSlowQuery(name string, start time.Time) {
	if duration := time.Since(start); duration > r.slowQueryThreshold {
		slog.Info("Slow query detected", slog.String("query", name), slog.Int64("duration_ms", duration.Milliseconds()))
	}
}

func scanListings(rows *sql.Rows) ([]domain.Listing, error) {
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

func (r *SQLiteRepository) FindByID(ctx context.Context, id string) (domain.Listing, error) {
	start := time.Now()
	defer func() {
		if duration := time.Since(start); duration > r.slowQueryThreshold {
			slog.Info("Slow query detected", slog.String("query", "FindByID"), slog.Int64("duration_ms", duration.Milliseconds()))
		}
	}()

	query := `
		SELECT ` + ListingSelectionsSQL + `
		FROM listings
		WHERE id = ?
	`
	row := r.readDB.QueryRowContext(ctx, query, id)

	l, err := scanListing(row)
	if err == sql.ErrNoRows {
		return domain.Listing{}, errors.New("listing not found")
	}
	return l, err
}

func (r *SQLiteRepository) FindByTitle(ctx context.Context, title string) ([]domain.Listing, error) {
	query := `
		SELECT ` + ListingSelectionsSQL + `
		FROM listings
		WHERE title = ?
	`
	rows, err := r.readDB.QueryContext(ctx, query, title)
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
	err := r.readDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM listings WHERE owner_id = ?", ownerID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT ` + ListingSelectionsSQL + `
              FROM listings 
              WHERE owner_id = ? 
              ORDER BY created_at DESC
              LIMIT ? OFFSET ?`

	rows, err := r.readDB.QueryContext(ctx, query, ownerID, limit, offset)
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
	err := r.readDB.QueryRowContext(ctx, ListingTitleExistsSQL, title).Scan(&exists)
	return exists, err
}

func (r *SQLiteRepository) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	start := time.Now()
	defer func() {
		if duration := time.Since(start); duration > r.slowQueryThreshold {
			slog.Info("Slow query detected", slog.String("query", "GetCounts"), slog.Int64("duration_ms", duration.Milliseconds()))
		}
	}()

	rows, err := r.readDB.QueryContext(ctx, ListingGetCountsSQL)
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
	rows, err := r.readDB.QueryContext(ctx, ListingGetLocationsSQL)
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
		whereClause += ListingFilterTypeSQL
		args = append(args, category)
	}

	// #nosec G202 - Dynamic query construction with trusted internal fragments
	query := `
		SELECT ` + ListingSelectionsSQL + `
		FROM listings 
		` + whereClause + `
		ORDER BY created_at DESC 
		LIMIT 3
	`
	rows, err := r.readDB.QueryContext(ctx, query, args...)
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

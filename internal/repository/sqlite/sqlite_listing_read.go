package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

type ListingFilters struct {
	Type            string
	QueryText       string
	OwnerID         string
	City            string
	WebsiteURL      string
	ListingStatus   domain.ListingStatus
	IncludedLat     float64
	IncludedLng     float64
	Radius          float64
	IncludeInactive bool
	FeaturedOnly    bool
}

func scanListing(s Scanner) (domain.Listing, error) {
	var l domain.Listing
	var deadline, eventStart, eventEnd, jobStart sql.NullTime
	var enrichmentAttemptedAtStr, ratingUpdatedAtStr sql.NullString

	err := s.Scan(
		&l.ID, &l.OwnerID, &l.OwnerOrigin, &l.Type, &l.Title, &l.Description,
		&l.City, &l.State, &l.Country, &l.Address, &l.HoursOfOperation, &l.ContactEmail, &l.ContactPhone, &l.ContactWhatsApp,
		&l.WebsiteURL, &l.ImageURL, &l.CreatedAt, &deadline, &l.IsActive,
		&eventStart, &eventEnd,
		&l.Skills, &jobStart, &l.JobApplyURL,
		&l.Company, &l.PayRange, &l.Status, &l.Featured,
		&l.HeatLevel, &l.RegionalSpecialty, &l.TopDish,
		&l.PaymentMethods, &l.MenuURL,
		&l.Latitude, &l.Longitude,
		&l.DeliveryPlatforms,
		&enrichmentAttemptedAtStr,
		&l.Rating, &l.ReviewCount,
		&ratingUpdatedAtStr,
		&l.StructuredHours,
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
	if enrichmentAttemptedAtStr.Valid {
		l.EnrichmentAttemptedAt = parseNullableTime(enrichmentAttemptedAtStr.String)
	}
	if ratingUpdatedAtStr.Valid {
		l.RatingUpdatedAt = parseNullableTime(ratingUpdatedAtStr.String)
	}
	return l, nil
}

func parseNullableTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	if idx := strings.Index(s, " m="); idx != -1 {
		s = s[:idx]
	}
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
	}
	for _, fmt := range formats {
		if t, err := time.Parse(fmt, s); err == nil {
			return &t
		}
	}
	slog.Warn("Failed to parse enrichment_attempted_at", slog.String("value", s))
	return nil
}

func (r *SQLiteRepository) FindAll(ctx context.Context, filterType string, queryText string, city string, lat float64, lng float64, radius float64, sortField string, sortOrder string, includeInactive bool, limit int, offset int) ([]domain.Listing, int, error) {
	start := time.Now()
	defer r.logSlowQuery("FindAll", start)

	filters := ListingFilters{
		Type:            filterType,
		QueryText:       queryText,
		City:            city,
		IncludedLat:     lat,
		IncludedLng:     lng,
		Radius:          radius,
		IncludeInactive: includeInactive,
	}
	where, args := r.buildListingWhere(filters)

	totalCount, err := r.getCount(ctx, "listings", where, args)
	if err != nil {
		return nil, 0, err
	}

	order := r.buildOrderClause(sortField, sortOrder)

	listings, err := r.queryListingsPaginated(ctx, where, order, args, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return listings, totalCount, nil
}

func (r *SQLiteRepository) queryListingsPaginated(ctx context.Context, where, order string, baseArgs []interface{}, limit, offset int) ([]domain.Listing, error) {
	args := make([]interface{}, 0, len(baseArgs)+2)
	args = append(args, baseArgs...)
	args = append(args, limit, offset)

	// #nosec G202 - Dynamic query construction with trusted internal fragments
	query := `SELECT ` + r.buildListingColumns() + ` FROM listings 
	          WHERE rowid IN (SELECT rowid FROM listings ` + where + ` ORDER BY ` + order + ` LIMIT ? OFFSET ?)
	          ORDER BY ` + order

	rows, err := r.readDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return scanListings(rows)
}

func (r *SQLiteRepository) buildListingColumns() string {
	return ListingSelectionsSQL
}

func (r *SQLiteRepository) buildListingWhere(filters ListingFilters) (string, []interface{}) {
	where := " WHERE 1=1"
	var args []interface{}

	if !filters.IncludeInactive {
		where += ` AND ` + ListingActiveApprovedSQL
	}

	if filters.Type != "" {
		where += ListingFilterTypeSQL
		args = append(args, filters.Type)
	}

	if filters.OwnerID != "" {
		where += ` AND owner_id = ?`
		args = append(args, filters.OwnerID)
	}

	if filters.Radius > 0 && filters.IncludedLat != 0 && filters.IncludedLng != 0 {
		// Bounding Box Optimization (Roughly 1 degree = 69 miles)
		latDelta := filters.Radius / 69.0
		lngDelta := filters.Radius / (69.0 * 0.707) // Approximation for mid-latitudes

		where += ` AND latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?`
		args = append(args, filters.IncludedLat-latDelta, filters.IncludedLat+latDelta, filters.IncludedLng-lngDelta, filters.IncludedLng+lngDelta)

		// Haversine formula for exact radius filtering
		// 3959 is the Earth's radius in miles
		where += ` AND (3959 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) <= ?`
		args = append(args, filters.IncludedLat, filters.IncludedLng, filters.IncludedLat, filters.Radius)
	} else if filters.City != "" {
		where += ` AND (city = ? OR address LIKE ?)`
		args = append(args, filters.City, "%"+filters.City+"%")
	}

	if filters.FeaturedOnly {
		where += ` AND featured = 1`
	}

	if filters.WebsiteURL != "" {
		where += ` AND website_url != ''`
	}

	if filters.QueryText != "" {
		where += ` AND rowid IN (SELECT rowid FROM listings_fts WHERE listings_fts MATCH ?)`
		args = append(args, filters.QueryText)
	}

	return where, args
}

func (r *SQLiteRepository) buildOrderClause(sortField, sortOrder string) string {
	if sortField == "" {
		return "featured DESC, heat_level DESC, rating DESC, created_at DESC"
	}

	field := "created_at"
	switch sortField {
	case "title":
		field = "title"
	case domain.FieldStatus:
		field = domain.FieldStatus
	case domain.FieldFeatured:
		field = domain.FieldFeatured
	case domain.FieldType:
		field = domain.FieldType
	}

	order := "DESC"
	if strings.ToLower(sortOrder) == "asc" {
		order = "ASC"
	}

	if field == domain.FieldFeatured {
		return domain.FieldFeatured + " " + order + ", created_at DESC"
	}
	return domain.FieldFeatured + " DESC, " + field + " " + order
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
	return scanAll(rows, scanListing)
}

func (r *SQLiteRepository) FindByID(ctx context.Context, id string) (domain.Listing, error) {
	start := time.Now()
	defer r.logSlowQuery("FindByID", start)

	// #nosec G202 - Dynamic query construction with trusted internal fragments
	query := `SELECT ` + r.buildListingColumns() + ` FROM listings WHERE id = ?`
	row := r.readDB.QueryRowContext(ctx, query, id)

	l, err := scanListing(row)
	if err == sql.ErrNoRows {
		return domain.Listing{}, ErrListingNotFound
	}
	return l, err
}

func (r *SQLiteRepository) FindByTitle(ctx context.Context, title string) ([]domain.Listing, error) {
	return r.queryListingsSimple(ctx, "WHERE title = ?", title)
}

func (r *SQLiteRepository) FindAllByOwner(ctx context.Context, ownerID string, limit int, offset int) ([]domain.Listing, int, error) {
	filters := ListingFilters{OwnerID: ownerID, IncludeInactive: true}
	where, args := r.buildListingWhere(filters)

	totalCount, err := r.getCount(ctx, "listings", where, args)
	if err != nil {
		return nil, 0, err
	}

	listings, err := r.queryListingsSimple(ctx, where+" ORDER BY created_at DESC LIMIT ? OFFSET ?", append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	return listings, totalCount, nil
}

func (r *SQLiteRepository) queryListingsSimple(ctx context.Context, where string, args ...interface{}) ([]domain.Listing, error) {
	// #nosec G202 - Dynamic query construction with trusted internal fragments
	query := `SELECT ` + r.buildListingColumns() + ` FROM listings ` + where
	rows, err := r.readDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return scanListings(rows)
}

// TitleExists checks if a listing with the given title exists using an efficient EXISTS query.
func (r *SQLiteRepository) TitleExists(ctx context.Context, title string) (bool, error) {
	var exists bool
	err := r.readDB.QueryRowContext(ctx, ListingTitleExistsSQL, title).Scan(&exists)
	return exists, err
}

func (r *SQLiteRepository) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	start := time.Now()
	defer r.logSlowQuery("GetCounts", start)

	rows, err := r.readDB.QueryContext(ctx, ListingGetCountsSQL)
	if err != nil {
		return nil, err
	}
	return scanCounts[domain.Category](rows)
}

func (r *SQLiteRepository) GetLocations(ctx context.Context) ([]domain.Location, error) {
	rows, err := r.readDB.QueryContext(ctx, ListingGetLocationsSQL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var locations []domain.Location
	for rows.Next() {
		var loc domain.Location
		if err := rows.Scan(&loc.City, &loc.State, &loc.Country); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, rows.Err()
}

// GetFeaturedListings returns featured listings set by admin, optionally filtered by category and city.
func (r *SQLiteRepository) GetFeaturedListings(ctx context.Context, category string, city string) ([]domain.Listing, error) {
	filters := ListingFilters{
		Type:         category,
		City:         city,
		FeaturedOnly: true,
	}
	where, args := r.buildListingWhere(filters)
	return r.queryListingsSimple(ctx, where+" ORDER BY created_at DESC LIMIT 3", args...)
}

func (r *SQLiteRepository) FindEnrichmentTargets(ctx context.Context, limit int) ([]domain.Listing, error) {
	// Custom WHERE for enrichment, still uses queryListingsSimple for scan logic
	where := "WHERE website_url != '' AND (heat_level = 0 OR menu_url = '' OR payment_methods = '') AND (enrichment_attempted_at IS NULL OR enrichment_attempted_at < datetime('now', '-7 days')) LIMIT ?"
	return r.queryListingsSimple(ctx, where, limit)
}

func (r *SQLiteRepository) FindRatingBackfillTargets(ctx context.Context, limit int) ([]domain.Listing, error) {
	where := "WHERE is_active = 1 AND status = 'Approved' AND type = 'Food' AND (rating_updated_at IS NULL OR rating_updated_at < datetime('now', '-30 days')) LIMIT ?"
	return r.queryListingsSimple(ctx, where, limit)
}

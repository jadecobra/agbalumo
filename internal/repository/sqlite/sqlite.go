package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	_ "modernc.org/sqlite" // register driver
)

const listingSelections = `
	id, owner_id, owner_origin, type, title, description,
	city, COALESCE(address, ''), hours_of_operation, contact_email, contact_phone, contact_whatsapp,
	website_url, image_url, created_at, deadline, is_active,
	event_start, event_end,
	COALESCE(skills, ''), job_start_date, COALESCE(job_apply_url, ''),
	COALESCE(company, ''), COALESCE(pay_range, '')
`

type Scanner interface {
	Scan(dest ...interface{}) error
}

func scanListing(s Scanner) (domain.Listing, error) {
	var l domain.Listing
	var deadline, eventStart, eventEnd, jobStart sql.NullTime

	err := s.Scan(
		&l.ID, &l.OwnerID, &l.OwnerOrigin, &l.Type, &l.Title, &l.Description,
		&l.City, &l.Address, &l.HoursOfOperation, &l.ContactEmail, &l.ContactPhone, &l.ContactWhatsApp,
		&l.WebsiteURL, &l.ImageURL, &l.CreatedAt, &deadline, &l.IsActive,
		&eventStart, &eventEnd,
		&l.Skills, &jobStart, &l.JobApplyURL,
		&l.Company, &l.PayRange,
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

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.migrate(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteRepository) migrate() error {
	// Create Listings Table
	createListingsTable := `
	CREATE TABLE IF NOT EXISTS listings (
		id TEXT PRIMARY KEY,
		owner_id TEXT,
		title TEXT,
		description TEXT,
		type TEXT,
		owner_origin TEXT,

		city TEXT,
		address TEXT,
		is_active BOOLEAN,
		created_at DATETIME,
		image_url TEXT,
		contact_email TEXT,
		contact_phone TEXT,
		contact_whatsapp TEXT,
		website_url TEXT,
		deadline DATETIME,
		skills TEXT,
		job_start_date DATETIME,
		job_apply_url TEXT,
		company TEXT,
		pay_range TEXT
	);`

	if _, err := r.db.ExecContext(context.Background(), createListingsTable); err != nil {
		return err
	}

	// Create Users Table
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        google_id TEXT UNIQUE,
        email TEXT,
        name TEXT,
        avatar_url TEXT,
        created_at DATETIME
    );`

	if _, err := r.db.ExecContext(context.Background(), createUsersTable); err != nil {
		return err
	}

	// Migration: Add address column if missing (simple check)
	// We ignore error if column exists (naive but works for dev SQLite)
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN address TEXT;")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN city TEXT;")
	// Add Hours of Operation
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN hours_of_operation TEXT DEFAULT '';")

	// Add Event Columns
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN event_start DATETIME;")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN event_end DATETIME;")

	// Add Job Columns
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN skills TEXT DEFAULT '';")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN job_start_date DATETIME;")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN job_apply_url TEXT DEFAULT '';")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN company TEXT DEFAULT '';")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN pay_range TEXT DEFAULT '';")

	return nil
}

// Save inserts or updates a listing.
func (r *SQLiteRepository) Save(ctx context.Context, l domain.Listing) error {
	query := `
	INSERT INTO listings (id, owner_id, title, description, type, owner_origin, city, address, hours_of_operation, is_active, created_at, image_url, contact_email, contact_phone, contact_whatsapp, website_url, deadline, event_start, event_end, skills, job_start_date, job_apply_url, company, pay_range)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		pay_range = excluded.pay_range;
	`

	_, err := r.db.ExecContext(ctx, query,
		l.ID, l.OwnerID, l.Title, l.Description, l.Type, l.OwnerOrigin, l.City, l.Address, l.HoursOfOperation, l.IsActive, l.CreatedAt,
		l.ImageURL, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.Deadline, l.EventStart, l.EventEnd,
		l.Skills, l.JobStartDate, l.JobApplyURL, l.Company, l.PayRange,
	)
	return err
}

func (r *SQLiteRepository) FindAll(ctx context.Context, filterType string, queryText string, includeInactive bool) ([]domain.Listing, error) {
	query := `SELECT ` + listingSelections + ` FROM listings WHERE 1=1`
	var args []interface{}

	if !includeInactive {
		query += ` AND is_active = true`
	}

	if filterType != "" {
		query += ` AND type = ?`
		args = append(args, filterType)
	}

	if queryText != "" {
		query += ` AND (title LIKE ? OR description LIKE ? OR city LIKE ?)`
		likeQuery := "%" + queryText + "%"
		args = append(args, likeQuery, likeQuery, likeQuery)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []domain.Listing
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, err
		}
		listings = append(listings, l)
	}
	return listings, nil
}

func (r *SQLiteRepository) FindByID(ctx context.Context, id string) (domain.Listing, error) {
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

// SaveUser inserts or updates a user.
func (r *SQLiteRepository) SaveUser(ctx context.Context, u domain.User) error {
	query := `
	INSERT INTO users (id, google_id, email, name, avatar_url, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(google_id) DO UPDATE SET
		email = excluded.email,
		name = excluded.name,
		avatar_url = excluded.avatar_url;
	`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.GoogleID, u.Email, u.Name, u.AvatarURL, u.CreatedAt,
	)
	return err
}

// FindUserByGoogleID retrieves a user by their Google ID.
func (r *SQLiteRepository) FindUserByGoogleID(ctx context.Context, googleID string) (domain.User, error) {
	query := `SELECT id, google_id, email, name, avatar_url, created_at FROM users WHERE google_id = ?`
	row := r.db.QueryRowContext(ctx, query, googleID)

	var u domain.User
	err := row.Scan(&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("user not found")
	}
	return u, err
}

// FindUserByID retrieves a user by their ID.
func (r *SQLiteRepository) FindUserByID(ctx context.Context, id string) (domain.User, error) {
	query := `SELECT id, google_id, email, name, avatar_url, created_at FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var u domain.User
	err := row.Scan(&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("user not found")
	}
	return u, err
}

func (r *SQLiteRepository) FindAllByOwner(ctx context.Context, ownerID string) ([]domain.Listing, error) {
	query := `SELECT ` + listingSelections + `
              FROM listings 
              WHERE owner_id = ? 
              ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []domain.Listing
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, err
		}
		listings = append(listings, l)
	}
	return listings, nil
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
	query := `SELECT type, COUNT(*) FROM listings WHERE is_active = true GROUP BY type`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[domain.Category]int)
	for rows.Next() {
		var cat domain.Category
		var count int
		if err := rows.Scan(&cat, &count); err != nil {
			return nil, err
		}
		counts[cat] = count
	}
	return counts, nil
}



func (r *SQLiteRepository) ExpireListings(ctx context.Context) (int64, error) {
	// Use Go's time to ensure driver handles serialization correctly and we control the timezone (UTC)
	now := time.Now().UTC()
	
	// Expire Requests past deadline AND Events past end time
	query := `
		UPDATE listings 
		SET is_active = false 
		WHERE is_active = true 
		AND (
			(type = 'Request' AND deadline < ?) 
			OR 
			(type = 'Event' AND event_end < ?)
			OR
			(type = 'Job' AND job_start_date < ?)
		)
	`
	// Added Job expiration rule for consistency
	// Passed 'now' 3 times for the 3 placeholders
	
	result, err := r.db.ExecContext(ctx, query, now, now, now.AddDate(0, 0, -90)) 
	// Note: Job rule was < now - 90 days. So we pass now.Add(-90 days).
	
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

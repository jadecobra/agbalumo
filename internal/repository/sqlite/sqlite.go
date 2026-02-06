package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jadecobra/agbalumo/internal/domain"
	_ "modernc.org/sqlite" // register driver
)

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
		deadline DATETIME
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

	return nil
}

// Save inserts or updates a listing.
func (r *SQLiteRepository) Save(ctx context.Context, l domain.Listing) error {
	query := `
	INSERT INTO listings (id, owner_id, title, description, type, owner_origin, city, address, is_active, created_at, image_url, contact_email, contact_phone, contact_whatsapp, website_url, deadline)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		owner_id = excluded.owner_id,
		title = excluded.title,
		description = excluded.description,
		type = excluded.type,
		owner_origin = excluded.owner_origin,
		city = excluded.city,
		address = excluded.address,
		is_active = excluded.is_active,
		image_url = excluded.image_url,
		contact_email = excluded.contact_email,
		contact_phone = excluded.contact_phone,
		contact_whatsapp = excluded.contact_whatsapp,
		website_url = excluded.website_url,
		deadline = excluded.deadline;
	`

	_, err := r.db.ExecContext(ctx, query,
		l.ID, l.OwnerID, l.Title, l.Description, l.Type, l.OwnerOrigin, l.City, l.Address, l.IsActive, l.CreatedAt,
		l.ImageURL, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.Deadline,
	)
	return err
}

func (r *SQLiteRepository) FindAll(ctx context.Context, filterType string, queryText string, includeInactive bool) ([]domain.Listing, error) {
	query := `SELECT id, owner_id, owner_origin, type, title, description, city, COALESCE(address, ''), contact_email, contact_phone, contact_whatsapp, website_url, image_url, created_at, deadline, is_active FROM listings WHERE 1=1`
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
		var l domain.Listing
		if err := rows.Scan(
			&l.ID, &l.OwnerID, &l.OwnerOrigin, &l.Type, &l.Title, &l.Description,
			&l.City, &l.Address, &l.ContactEmail, &l.ContactPhone, &l.ContactWhatsApp,
			&l.WebsiteURL, &l.ImageURL, &l.CreatedAt, &l.Deadline, &l.IsActive,
		); err != nil {
			return nil, err
		}
		listings = append(listings, l)
	}
	return listings, nil
}

func (r *SQLiteRepository) FindByID(ctx context.Context, id string) (domain.Listing, error) {
	query := `
		SELECT id, owner_id, title, description, type, owner_origin, city, COALESCE(address, ''), is_active, created_at, image_url, contact_email, contact_phone, contact_whatsapp, website_url, deadline
		FROM listings
		WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var l domain.Listing
	err := row.Scan(
		&l.ID, &l.OwnerID, &l.Title, &l.Description, &l.Type, &l.OwnerOrigin,
		&l.City, &l.Address, &l.IsActive, &l.CreatedAt,
		&l.ImageURL, &l.ContactEmail, &l.ContactPhone, &l.ContactWhatsApp, &l.WebsiteURL, &l.Deadline,
	)
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
	query := `SELECT id, owner_id, owner_origin, type, title, description, city, COALESCE(address, ''), contact_email, contact_phone, contact_whatsapp, website_url, image_url, created_at, deadline, is_active 
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
		var l domain.Listing
		if err := rows.Scan(
			&l.ID, &l.OwnerID, &l.OwnerOrigin, &l.Type, &l.Title, &l.Description,
			&l.City, &l.Address, &l.ContactEmail, &l.ContactPhone, &l.ContactWhatsApp,
			&l.WebsiteURL, &l.ImageURL, &l.CreatedAt, &l.Deadline, &l.IsActive,
		); err != nil {
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

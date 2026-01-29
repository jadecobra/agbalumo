package sqlite

import (
	"context"
	"database/sql"

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
	query := `
	CREATE TABLE IF NOT EXISTS listings (
		id TEXT PRIMARY KEY,
		owner_origin TEXT,
		type TEXT,
		title TEXT,
		description TEXT,
		neighborhood TEXT,
		contact_email TEXT,
		contact_phone TEXT,
		contact_whatsapp TEXT,
		website_url TEXT,
		image_url TEXT,
		created_at DATETIME,
		deadline DATETIME,
		is_active BOOLEAN
	);`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Save(ctx context.Context, l domain.Listing) error {
	query := `
	INSERT INTO listings (id, owner_origin, type, title, description, neighborhood, contact_email, contact_phone, contact_whatsapp, website_url, image_url, created_at, deadline, is_active)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		owner_origin=excluded.owner_origin,
		type=excluded.type,
		title=excluded.title,
		description=excluded.description,
		neighborhood=excluded.neighborhood,
		contact_email=excluded.contact_email,
		contact_phone=excluded.contact_phone,
		contact_whatsapp=excluded.contact_whatsapp,
		website_url=excluded.website_url,
		image_url=excluded.image_url,
		deadline=excluded.deadline,
		is_active=excluded.is_active;
	`
	_, err := r.db.ExecContext(ctx, query,
		l.ID, l.OwnerOrigin, l.Type, l.Title, l.Description, l.Neighborhood,
		l.ContactEmail, l.ContactPhone, l.ContactWhatsApp, l.WebsiteURL, l.ImageURL, l.CreatedAt, l.Deadline, l.IsActive,
	)
	return err
}

func (r *SQLiteRepository) FindAll(ctx context.Context, filterType string) ([]domain.Listing, error) {
	query := `SELECT id, owner_origin, type, title, description, neighborhood, contact_email, contact_phone, contact_whatsapp, website_url, image_url, created_at, deadline, is_active FROM listings`
	var args []interface{}

	if filterType != "" {
		query += ` WHERE type = ?`
		args = append(args, filterType)
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
		var typeStr string
		err := rows.Scan(
			&l.ID,
			&l.OwnerOrigin,
			&typeStr,
			&l.Title,
			&l.Description,
			&l.Neighborhood,
			&l.ContactEmail,
			&l.ContactPhone,
			&l.ContactWhatsApp,
			&l.WebsiteURL,
			&l.ImageURL,
			&l.CreatedAt,
			&l.Deadline,
			&l.IsActive,
		)
		if err != nil {
			return nil, err
		}
		l.Type = domain.Category(typeStr)
		listings = append(listings, l)
	}
	return listings, nil
}

func (r *SQLiteRepository) FindByID(ctx context.Context, id string) (domain.Listing, error) {
	query := `SELECT id, owner_origin, type, title, description, neighborhood, contact_email, contact_phone, contact_whatsapp, website_url, image_url, created_at, deadline, is_active FROM listings WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var l domain.Listing
	var typeStr string
	err := row.Scan(
		&l.ID, &l.OwnerOrigin, &typeStr, &l.Title, &l.Description, &l.Neighborhood,
		&l.ContactEmail, &l.ContactPhone, &l.ContactWhatsApp, &l.WebsiteURL, &l.ImageURL, &l.CreatedAt, &l.Deadline, &l.IsActive,
	)
	if err != nil {
		return domain.Listing{}, err
	}
	l.Type = domain.Category(typeStr)
	return l, nil
}

package sqlite

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite" // register driver
)

type Scanner interface {
	Scan(dest ...interface{}) error
}

type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepositoryFromDB creates a new repository using an existing DB connection.
func NewSQLiteRepositoryFromDB(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout=5000;"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		return nil, err
	}

	if dbPath == ":memory:" {
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	} else {
		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(100)
	}
	db.SetConnMaxLifetime(0)

	repo := &SQLiteRepository{db: db}
	if err := repo.migrate(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteRepository) migrate() error {
	createListingsTable := `
	CREATE TABLE IF NOT EXISTS listings (
		id TEXT PRIMARY KEY,
		owner_id TEXT,
		title TEXT,
		description TEXT,
		type TEXT,
		owner_origin TEXT,
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

	createFeedbackTable := `
	CREATE TABLE IF NOT EXISTS feedback (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		type TEXT,
		content TEXT,
		created_at DATETIME
	);`

	if _, err := r.db.ExecContext(context.Background(), createFeedbackTable); err != nil {
		return err
	}

	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN address TEXT;")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN city TEXT;")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN hours_of_operation TEXT DEFAULT '';")
	_, _ = r.db.ExecContext(context.Background(), "UPDATE listings SET city = '' WHERE city IS NULL OR city = '';")
	_, _ = r.db.ExecContext(context.Background(), "UPDATE listings SET city = '' WHERE city = 'Unknown';")

	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN event_start DATETIME;")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN event_end DATETIME;")

	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN skills TEXT DEFAULT '';")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN job_start_date DATETIME;")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN job_apply_url TEXT DEFAULT '';")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN company TEXT DEFAULT '';")
	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN pay_range TEXT DEFAULT '';")

	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN status TEXT DEFAULT 'Approved';")

	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE users ADD COLUMN role TEXT DEFAULT 'User';")

	_, _ = r.db.ExecContext(context.Background(), "ALTER TABLE listings ADD COLUMN featured BOOLEAN DEFAULT 0;")

	_, _ = r.db.ExecContext(context.Background(), "CREATE UNIQUE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);")

	_, _ = r.db.ExecContext(context.Background(), "CREATE INDEX IF NOT EXISTS idx_listings_owner_id ON listings(owner_id);")
	_, _ = r.db.ExecContext(context.Background(), "CREATE INDEX IF NOT EXISTS idx_listings_filter_sort ON listings(is_active, status, type, created_at DESC);")

	_, _ = r.db.ExecContext(context.Background(), `
		CREATE VIRTUAL TABLE IF NOT EXISTS listings_fts USING fts5(
			title, description, city,
			content=listings,
			content_rowid=rowid,
			tokenize='trigram'
		);
	`)

	_, _ = r.db.ExecContext(context.Background(), `
		CREATE TRIGGER IF NOT EXISTS listings_ai AFTER INSERT ON listings BEGIN
			INSERT INTO listings_fts(rowid, title, description, city)
			VALUES (new.rowid, new.title, new.description, new.city);
		END;
	`)
	_, _ = r.db.ExecContext(context.Background(), `
		CREATE TRIGGER IF NOT EXISTS listings_ad AFTER DELETE ON listings BEGIN
			INSERT INTO listings_fts(listings_fts, rowid, title, description, city)
			VALUES ('delete', old.rowid, old.title, old.description, old.city);
		END;
	`)
	_, _ = r.db.ExecContext(context.Background(), `
		CREATE TRIGGER IF NOT EXISTS listings_au AFTER UPDATE ON listings BEGIN
			INSERT INTO listings_fts(listings_fts, rowid, title, description, city)
			VALUES ('delete', old.rowid, old.title, old.description, old.city);
			INSERT INTO listings_fts(rowid, title, description, city)
			VALUES (new.rowid, new.title, new.description, new.city);
		END;
	`)

	_, _ = r.db.ExecContext(context.Background(), "INSERT INTO listings_fts(listings_fts) VALUES('rebuild');")

	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id TEXT PRIMARY KEY,
		name TEXT,
		claimable BOOLEAN,
		is_system BOOLEAN,
		active BOOLEAN,
		requires_special_validation BOOLEAN,
		created_at DATETIME,
		updated_at DATETIME
	);`

	if _, err := r.db.ExecContext(context.Background(), createCategoriesTable); err != nil {
		return err
	}

	createClaimRequestsTable := `
	CREATE TABLE IF NOT EXISTS claim_requests (
		id TEXT PRIMARY KEY,
		listing_id TEXT NOT NULL,
		listing_title TEXT,
		user_id TEXT NOT NULL,
		user_name TEXT,
		user_email TEXT,
		status TEXT NOT NULL DEFAULT 'Pending',
		created_at DATETIME
	);`

	if _, err := r.db.ExecContext(context.Background(), createClaimRequestsTable); err != nil {
		return err
	}

	_, _ = r.db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_claim_requests_user_listing ON claim_requests(user_id, listing_id);`)
	_, _ = r.db.ExecContext(context.Background(), `CREATE INDEX IF NOT EXISTS idx_claim_requests_status ON claim_requests(status);`)

	_, err := r.db.ExecContext(context.Background(), "ALTER TABLE categories ADD COLUMN active_fixed BOOLEAN DEFAULT 0;")
	if err == nil {
		_, _ = r.db.ExecContext(context.Background(), "UPDATE categories SET active = 1, active_fixed = 1;")
	}

	return nil
}

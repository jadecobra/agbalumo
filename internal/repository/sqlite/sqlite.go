package sqlite

import (
	"database/sql"
	"embed"
	"strings"
	"time"

	_ "modernc.org/sqlite" // register driver
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type Scanner interface {
	Scan(dest ...interface{}) error
}

type SQLiteRepository struct {
	db                 *sql.DB
	slowQueryThreshold time.Duration
}

// NewSQLiteRepositoryFromDB creates a new repository using an existing DB connection.
func NewSQLiteRepositoryFromDB(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db:                 db,
		slowQueryThreshold: 50 * time.Millisecond,
	}
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

	repo := &SQLiteRepository{
		db:                 db,
		slowQueryThreshold: 50 * time.Millisecond,
	}
	if err := repo.migrate(); err != nil {
		return nil, err
	}

	return repo, nil
}

// SetSlowQueryThreshold updates the threshold for logging slow queries.
func (r *SQLiteRepository) SetSlowQueryThreshold(d time.Duration) {
	r.slowQueryThreshold = d
}

func (r *SQLiteRepository) migrate() error {
	// Create schema_migrations table if it doesn't exist
	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);`)
	if err != nil {
		return err
	}

	// Read migration files from embedded FS
	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()

		// Check if this migration has already been applied
		var exists int
		err := r.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", name).Scan(&exists)
		if err != nil {
			return err
		}
		if exists > 0 {
			continue
		}

		// Read and execute the migration
		contentRaw, err := migrationFS.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}
		content := string(contentRaw)

		// Execute the migration content statement by statement.
		// We use a custom separator to avoid breaking SQLite triggers with internal semicolons.
		statements := strings.Split(content, "-- STATEMENT")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			_, err = r.db.Exec(stmt)
			if err != nil {
				// For the initial legacy migration (001), we ignore errors (like "duplicate column")
				// to ensure backward compatibility with partially migrated databases.
				if name == "001_initial_schema.sql" {
					continue
				}
				return err
			}
		}

		// Record the migration as applied
		_, err = r.db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close closes the underlying database connection.
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

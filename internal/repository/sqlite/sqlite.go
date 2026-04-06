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
	writeDB            *sql.DB
	readDB             *sql.DB
	slowQueryThreshold time.Duration
}

// NewSQLiteRepositoryFromDB creates a new repository using an existing DB connection for both pools.
func NewSQLiteRepositoryFromDB(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		writeDB:            db,
		readDB:             db,
		slowQueryThreshold: 50 * time.Millisecond,
	}
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	// 1. Open the Write Database (MaxOpenConns=1)
	writeDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := writeDB.Ping(); err != nil {
		return nil, err
	}

	// Applying Pragmas (most important for journal/sync)
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA foreign_keys=ON;",
	}

	for _, p := range pragmas {
		if _, err := writeDB.Exec(p); err != nil {
			return nil, err
		}
	}

	if dbPath == ":memory:" {
		writeDB.SetMaxOpenConns(1)
		writeDB.SetMaxIdleConns(1)
		writeDB.SetConnMaxLifetime(0)
		repo := &SQLiteRepository{
			writeDB:            writeDB,
			readDB:             writeDB, // Same pool for memory
			slowQueryThreshold: 50 * time.Millisecond,
		}
		if err := repo.migrate(); err != nil {
			return nil, err
		}
		return repo, nil
	}

	// 2. Open the Read Database (MaxOpenConns=100)
	readDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := readDB.Ping(); err != nil {
		return nil, err
	}

	for _, p := range pragmas {
		if _, err := readDB.Exec(p); err != nil {
			return nil, err
		}
	}

	// Configure pools for file-based database
	writeDB.SetMaxOpenConns(1)
	writeDB.SetMaxIdleConns(1)
	readDB.SetMaxOpenConns(100)
	readDB.SetMaxIdleConns(100)
	writeDB.SetConnMaxLifetime(0)
	readDB.SetConnMaxLifetime(0)

	repo := &SQLiteRepository{
		writeDB:            writeDB,
		readDB:             readDB,
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
	_, err := r.writeDB.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);`)
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
		err := r.writeDB.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", name).Scan(&exists)
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
			_, err = r.writeDB.Exec(stmt)
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
		_, err = r.writeDB.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close closes both underlying database connections.
func (r *SQLiteRepository) Close() error {
	if r.writeDB == r.readDB {
		return r.writeDB.Close()
	}
	err1 := r.writeDB.Close()
	err2 := r.readDB.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

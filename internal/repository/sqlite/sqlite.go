package sqlite

import (
	"database/sql"
	"embed"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite" // register driver

	"github.com/jadecobra/agbalumo/internal/domain"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

const defaultSlowQueryThreshold = 50 * time.Millisecond

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
		slowQueryThreshold: defaultSlowQueryThreshold,
	}
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	writeDB, err := sql.Open(domain.SQLiteDriver, dbPath)
	if err != nil {
		return nil, err
	}
	if err := applyPragmas(writeDB); err != nil {
		return nil, err
	}

	readDB := writeDB
	if dbPath != domain.SQLiteMemory {
		readDB, err = sql.Open(domain.SQLiteDriver, dbPath)
		if err != nil {
			return nil, err
		}
		if err := applyPragmas(readDB); err != nil {
			return nil, err
		}
	}

	configurePools(writeDB, readDB, dbPath == domain.SQLiteMemory)

	repo := &SQLiteRepository{
		writeDB:            writeDB,
		readDB:             readDB,
		slowQueryThreshold: defaultSlowQueryThreshold,
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

func applyPragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA foreign_keys=ON;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return err
		}
	}
	return nil
}

func configurePools(writeDB, readDB *sql.DB, isMemory bool) {
	writeDB.SetMaxOpenConns(1)
	writeDB.SetMaxIdleConns(1)
	writeDB.SetConnMaxLifetime(0)
	if !isMemory {
		readDB.SetMaxOpenConns(100)
		readDB.SetMaxIdleConns(100)
		readDB.SetConnMaxLifetime(0)
	}
}

func (r *SQLiteRepository) migrate() error {
	_, err := r.writeDB.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);`)
	if err != nil {
		return err
	}

	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if err := r.applyMigrationEntry(entry); err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLiteRepository) applyMigrationEntry(entry os.DirEntry) error {
	if entry.IsDir() {
		return nil // skip
	}
	name := entry.Name()

	var exists int
	err := r.writeDB.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", name).Scan(&exists)
	if err != nil || exists > 0 {
		return err
	}

	content, err := migrationFS.ReadFile("migrations/" + name)
	if err != nil {
		return err
	}

	if err := r.executeMigration(name, string(content)); err != nil {
		return err
	}

	_, err = r.writeDB.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name)
	return err
}

func (r *SQLiteRepository) executeMigration(name, content string) error {
	statements := strings.Split(content, "-- STATEMENT")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := r.writeDB.Exec(stmt); err != nil {
			if name == "001_initial_schema.sql" {
				continue // ignore legacy issues
			}
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

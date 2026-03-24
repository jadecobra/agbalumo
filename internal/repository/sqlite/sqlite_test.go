package sqlite_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func newTestRepo(t *testing.T) (*sqlite.SQLiteRepository, string) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	dsn := dbPath + "?_time_format=sqlite"
	repo, err := sqlite.NewSQLiteRepository(dsn)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	return repo, dbPath
}

func TestNewSQLiteRepositoryFromDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}
	defer func() { _ = db.Close() }()

	repo := sqlite.NewSQLiteRepositoryFromDB(db)
	if repo == nil {
		t.Fatal("Expected repository, got nil")
	}

	// Verify we can use it
	ctx := context.Background()
	_, _, err = repo.FindAll(ctx, "All", "", "", "", false, 20, 0)
	if err == nil {
		t.Error("Expected error due to missing tables, got nil")
	}
}

func TestSQLiteRepository_InitializationErrors(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file")
	_ = os.WriteFile(filePath, []byte("content"), 0644)

	// Try to use that file as a directory for the DB
	_, err := sqlite.NewSQLiteRepository(filePath)
	if err == nil {
		t.Error("Expected error initializing repo on a plain file without write access as DB, got nil")
	}
}

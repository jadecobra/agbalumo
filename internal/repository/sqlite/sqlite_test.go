package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	_ "modernc.org/sqlite"
	"sync"
)

var (
	dbCounter int64
	counterMu sync.Mutex
)

func newTestRepo(t *testing.T) (*sqlite.SQLiteRepository, string) {
	counterMu.Lock()
	dbCounter++
	id := dbCounter
	counterMu.Unlock()

	// Use a unique name per test to ensure isolation in shared-cache mode
	dbName := fmt.Sprintf("file:test_%s_%d?mode=memory&cache=shared&_time_format=sqlite", t.Name(), id)
	repo, err := sqlite.NewSQLiteRepository(dbName)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	return repo, dbName
}

func TestNewSQLiteRepositoryFromDB(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_time_format=sqlite")
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
	_ = os.WriteFile( /*nolint:gosec*/ filePath, []byte("content"), 0600)

	// Try to use that file as a directory for the DB
	_, err := sqlite.NewSQLiteRepository(filePath)
	if err == nil {
		t.Error("Expected error initializing repo on a plain file without write access as DB, got nil")
	}
}

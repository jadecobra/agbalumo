package sqlite_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/testutil"
	_ "modernc.org/sqlite"
)

// Local counter logic moved to testutil.

func TestNewSQLiteRepositoryFromDB(t *testing.T) {
	t.Parallel()
	db := testutil.SetupTestDB(t)
	defer func() { _ = db.Close() }()

	repo := sqlite.NewSQLiteRepositoryFromDB(db)
	if repo == nil {
		t.Fatal("Expected repository, got nil")
	}

	// Verify we can use it
	ctx := context.Background()
	_, _, err := repo.FindAll(ctx, "All", "", "", "", false, 20, 0)
	if err == nil {
		t.Error("Expected error due to missing tables, got nil")
	}
}

func TestSQLiteRepository_InitializationErrors(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file")
	_ = os.WriteFile( /*nolint:gosec*/ filePath, []byte("content"), 0600)

	// Try to use that file as a directory for the DB
	_, err := sqlite.NewSQLiteRepository(filePath)
	if err == nil {
		t.Error("Expected error initializing repo on a plain file without write access as DB, got nil")
	}
}

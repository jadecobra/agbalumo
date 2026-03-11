package handler

import (
	"testing"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

// SetupTestRepository initializes an in-memory sqlite repository for integration tests.
func SetupTestRepository(t *testing.T) *sqlite.SQLiteRepository {
	t.Helper()
	// Using ":memory:" creates a new, empty, in-memory database every time.
	// Since we set MaxOpenConns(1) in the repository, this works well for a single instance.
	repo, err := sqlite.NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory sqlite repository: %v", err)
	}
	return repo
}

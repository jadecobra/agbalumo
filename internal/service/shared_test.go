package service

import (
	"testing"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func setupTestRepo(t *testing.T) *sqlite.SQLiteRepository {
	repo, err := sqlite.NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory sqlite repository: %v", err)
	}
	return repo
}

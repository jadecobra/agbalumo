package sqlite_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func TestDeliveryPlatformsDiskParity(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create file-backed repo: %v", err)
	}
	defer func() {
		_ = repo.Close()
	}()

	ctx := context.Background()
	l := domain.Listing{
		ID:                "dp-disk",
		Title:             "Disk Parity Test",
		DeliveryPlatforms: `["UberEats"]`,
		IsActive:          true,
	}

	err = repo.Save(ctx, l)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	found, err := repo.FindByID(ctx, "dp-disk")
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}
	if found.DeliveryPlatforms != l.DeliveryPlatforms {
		t.Errorf("Expected %q, got %q", l.DeliveryPlatforms, found.DeliveryPlatforms)
	}
}

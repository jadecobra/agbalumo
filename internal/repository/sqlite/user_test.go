package sqlite_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func TestUserOperations(t *testing.T) {
	// Setup temporary DB
	tmpFile, err := os.CreateTemp("", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	repo, err := sqlite.NewSQLiteRepository(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Create User
	user := domain.User{
		ID:        "user-1",
		GoogleID:  "google-1",
		Email:     "test@example.com",
		Name:      "Test User",
		AvatarURL: "http://avatar.com",
		CreatedAt: time.Now(),
	}

	if err := repo.SaveUser(ctx, user); err != nil {
		t.Fatalf("SaveUser failed: %v", err)
	}

	// Find By ID
	found, err := repo.FindUserByID(ctx, "user-1")
	if err != nil {
		t.Fatalf("FindUserByID failed: %v", err)
	}
	if found.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, found.Email)
	}

	// Find By Google ID
	foundG, err := repo.FindUserByGoogleID(ctx, "google-1")
	if err != nil {
		t.Fatalf("FindUserByGoogleID failed: %v", err)
	}
	if foundG.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, foundG.Name)
	}

	// Update User
	user.Name = "Updated Name"
	if err := repo.SaveUser(ctx, user); err != nil {
		t.Fatalf("SaveUser (update) failed: %v", err)
	}

	foundUpdated, err := repo.FindUserByID(ctx, "user-1")
	if err != nil {
		t.Fatal(err)
	}
	if foundUpdated.Name != "Updated Name" {
		t.Errorf("Expected updated name, got %s", foundUpdated.Name)
	}
}

func TestFindUser_NotFound(t *testing.T) {
	// Setup
	tmpFile, _ := os.CreateTemp("", "test_nf.db")
	defer os.Remove(tmpFile.Name())
	repo, _ := sqlite.NewSQLiteRepository(tmpFile.Name())
	ctx := context.Background()

	// By ID
	_, err := repo.FindUserByID(ctx, "missing")
	if err == nil {
		t.Error("Expected error for missing user ID")
	}

	// By Google ID
	_, err = repo.FindUserByGoogleID(ctx, "missing-g")
	if err == nil {
		t.Error("Expected error for missing google ID")
	}
}

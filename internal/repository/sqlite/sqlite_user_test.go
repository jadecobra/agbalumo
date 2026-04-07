package sqlite_test

import (
	"context"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func TestGetUserCount(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Add Users
	_ = repo.SaveUser(ctx, domain.User{ID: "u1", GoogleID: "g1", Email: "e1", CreatedAt: time.Now()})
	_ = repo.SaveUser(ctx, domain.User{ID: "u2", GoogleID: "g2", Email: "e2", CreatedAt: time.Now()})

	c, err := repo.GetUserCount(ctx)
	if err != nil {
		t.Fatalf("GetUserCount failed: %v", err)
	}
	if c != 2 {
		t.Errorf("Expected 2 users, got %d", c)
	}
}

func TestGetAllUsers(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Seed
	_ = repo.SaveUser(ctx, domain.User{ID: "u1", GoogleID: "g1", Email: "e1", CreatedAt: time.Now()})
	_ = repo.SaveUser(ctx, domain.User{ID: "u2", GoogleID: "g2", Email: "e2", CreatedAt: time.Now().Add(time.Hour)})

	users, err := repo.GetAllUsers(ctx, 100, 0)
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
	// Should be ordered by created_at DESC
	if users[0].ID != "u2" {
		t.Errorf("Expected u2 first (newest), got %s", users[0].ID)
	}
}

func TestGetUserGrowth(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()
	now := time.Now().UTC()

	// 2. Users Growth
	_ = repo.SaveUser(ctx, domain.User{ID: "u1", GoogleID: "g1", Email: "e1", CreatedAt: now.Add(-24 * time.Hour)})
	_ = repo.SaveUser(ctx, domain.User{ID: "u2", GoogleID: "g2", Email: "e2", CreatedAt: now})

	metrics, err := repo.GetUserGrowth(ctx)
	if err != nil {
		t.Fatalf("GetUserGrowth failed: %v", err)
	}
	if len(metrics) != 2 {
		t.Errorf("Expected 2 days of user metrics, got %d", len(metrics))
	}
}

func TestUserOperations(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

	err := repo.SaveUser(ctx, user)
	if err != nil {
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
	err = repo.SaveUser(ctx, user)
	if err != nil {
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
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

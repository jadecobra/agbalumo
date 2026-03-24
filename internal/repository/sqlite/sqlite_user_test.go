package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func TestGetUserCount(t *testing.T) {
	repo, _ := newTestRepo(t)
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
	repo, _ := newTestRepo(t)
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
	repo, _ := newTestRepo(t)
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

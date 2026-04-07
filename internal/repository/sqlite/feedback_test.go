package sqlite_test

import (
	"context"
	"database/sql"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	_ "modernc.org/sqlite" // registers sqlite driver
)

func TestSaveFeedback(t *testing.T) {
	repo, dbName := testutil.SetupTestRepositoryUnique(t)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	t.Run("Successfully save feedback", func(t *testing.T) {
		ctx := context.Background()
		feedback := domain.Feedback{
			ID:        uuid.New().String(),
			UserID:    "user-123",
			Type:      domain.FeedbackTypeIssue,
			Content:   "This is a test issue.",
			CreatedAt: time.Now(),
		}

		if err := repo.SaveFeedback(ctx, feedback); err != nil {
			t.Fatalf("SaveFeedback failed: %v", err)
		}

		var count int
		_ = db.QueryRow("SELECT COUNT(*) FROM feedback WHERE id = ?", feedback.ID).Scan(&count)
		if count != 1 {
			t.Errorf("Expected 1 feedback, got %d", count)
		}
	})

	t.Run("Save multiple feedbacks", func(t *testing.T) {
		ctx := context.Background()
		f := domain.Feedback{
			ID:        uuid.New().String(),
			UserID:    "user-1",
			Type:      domain.FeedbackTypeFeature,
			Content:   "Feature 1",
			CreatedAt: time.Now(),
		}

		if err := repo.SaveFeedback(ctx, f); err != nil {
			t.Errorf("Failed to save f: %v", err)
		}

		var count int
		_ = db.QueryRow("SELECT COUNT(*) FROM feedback").Scan(&count)
		if count < 2 {
			t.Errorf("Expected at least 2 feedbacks, got %d", count)
		}
	})
}

func TestGetAllFeedback(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Seed Data
	feedbacks := []domain.Feedback{
		{ID: "f1", UserID: "user1", Type: domain.FeedbackTypeIssue, Content: "Oldest", CreatedAt: time.Now().Add(-2 * time.Hour)},
		{ID: "f2", UserID: "user2", Type: domain.FeedbackTypeFeature, Content: "Newest", CreatedAt: time.Now()},
	}

	for _, f := range feedbacks {
		if err := repo.SaveFeedback(ctx, f); err != nil {
			t.Fatalf("Failed to save feedback: %v", err)
		}
	}

	retrieved, err := repo.GetAllFeedback(ctx)
	if err != nil {
		t.Fatalf("Failed to get all feedback: %v", err)
	}

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 items, got %d", len(retrieved))
	}

	if len(retrieved) > 0 && retrieved[0].ID != "f2" {
		t.Errorf("Expected first item to be f2 (Newest), got %s", retrieved[0].ID)
	}
}

func TestGetFeedbackCounts(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Seed Feedback
	f1 := domain.Feedback{ID: uuid.New().String(), UserID: "u1", Type: domain.FeedbackTypeIssue, Content: "Bug 1", CreatedAt: time.Now()}
	f2 := domain.Feedback{ID: uuid.New().String(), UserID: "u2", Type: domain.FeedbackTypeIssue, Content: "Bug 2", CreatedAt: time.Now()}
	f3 := domain.Feedback{ID: uuid.New().String(), UserID: "u3", Type: domain.FeedbackTypeFeature, Content: "Feature 1", CreatedAt: time.Now()}

	if err := repo.SaveFeedback(ctx, f1); err != nil {
		t.Fatalf("Failed to save f1: %v", err)
	}
	if err := repo.SaveFeedback(ctx, f2); err != nil {
		t.Fatalf("Failed to save f2: %v", err)
	}
	if err := repo.SaveFeedback(ctx, f3); err != nil {
		t.Fatalf("Failed to save f3: %v", err)
	}

	counts, err := repo.GetFeedbackCounts(ctx)
	if err != nil {
		t.Fatalf("GetFeedbackCounts failed: %v", err)
	}

	if counts[domain.FeedbackTypeIssue] != 2 {
		t.Errorf("Expected 2 bug reports, got %d", counts[domain.FeedbackTypeIssue])
	}
	if counts[domain.FeedbackTypeFeature] != 1 {
		t.Errorf("Expected 1 feature request, got %d", counts[domain.FeedbackTypeFeature])
	}
	if counts[domain.FeedbackTypeOther] != 0 {
		t.Errorf("Expected 0 other, got %d", counts[domain.FeedbackTypeOther])
	}
}

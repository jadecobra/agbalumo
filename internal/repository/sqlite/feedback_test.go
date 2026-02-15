package sqlite_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	_ "modernc.org/sqlite" // registers sqlite driver
)


func TestSaveFeedback(t *testing.T) {
	// Setup temporary DB
	tmpFile, err := os.CreateTemp("", "test_feedback.db")
	if err != nil {
		t.Fatal(err)
	}
	dbPath := tmpFile.Name()
	defer os.Remove(dbPath)

	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	
	// Open raw DB specifically for verification
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	t.Run("Successfully save feedback", func(t *testing.T) {
		ctx := context.Background()
		feedback := domain.Feedback{
			ID:        uuid.New().String(),
			UserID:    "user-123",
			Type:      domain.FeedbackTypeIssue,
			Content:   "This is a test issue.",
			CreatedAt: time.Now(),
		}

		err := repo.SaveFeedback(ctx, feedback)
		if err != nil {
			t.Fatalf("SaveFeedback failed: %v", err)
		}

		// Verify it was saved
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM feedback WHERE id = ?", feedback.ID).Scan(&count)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 feedback, got %d", count)
		}
		
		var content string
		var fbType string
		err = db.QueryRow("SELECT content, type FROM feedback WHERE id = ?", feedback.ID).Scan(&content, &fbType)
		if err != nil {
			t.Fatalf("Query details failed: %v", err)
		}
		
		if content != "This is a test issue." {
			t.Errorf("Expected content 'This is a test issue.', got '%s'", content)
		}
		if fbType != "Issue" {
			t.Errorf("Expected type 'Issue', got '%s'", fbType)
		}
	})
	
	t.Run("Save multiple feedbacks", func(t *testing.T) {
		ctx := context.Background()
		f1 := domain.Feedback{
			ID:        uuid.New().String(),
			UserID:    "user-1",
			Type:      domain.FeedbackTypeFeature,
			Content:   "Feature 1",
			CreatedAt: time.Now(),
		}
		
		if err := repo.SaveFeedback(ctx, f1); err != nil {
			t.Errorf("Failed to save f1: %v", err)
		}
		
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM feedback").Scan(&count)
		if err != nil {
			t.Errorf("Count query failed: %v", err)
		}
		// 1 from previous test + 1 new one = 2
		if count != 2 {
			t.Errorf("Expected 2 feedbacks, got %d", count)
		}
	})
}

func TestGetAllFeedback(t *testing.T) {
	// Setup temporary DB (copy paste setup for now or refactor helper if available in package)
	tmpFile, err := os.CreateTemp("", "test_feedback_get.db")
	if err != nil {
		t.Fatal(err)
	}
	dbPath := tmpFile.Name()
	defer os.Remove(dbPath)

	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	
	ctx := context.Background()

	// Seed Data
	feedbacks := []domain.Feedback{
		{
			ID:        "f1",
			UserID:    "user1",
			Type:      domain.FeedbackTypeIssue,
			Content:   "Oldest",
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        "f2",
			UserID:    "user2",
			Type:      domain.FeedbackTypeFeature,
			Content:   "Newest",
			CreatedAt: time.Now(),
		},
	}

	for _, f := range feedbacks {
		if err := repo.SaveFeedback(ctx, f); err != nil {
			t.Fatalf("Failed to save feedback: %v", err)
		}
	}

	// Test GetAll
	retrieved, err := repo.GetAllFeedback(ctx)
	if err != nil {
		t.Fatalf("Failed to get all feedback: %v", err)
	}

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 feedback items, got %d", len(retrieved))
	}

	// Verify order (newest first)
	if len(retrieved) > 0 {
		if retrieved[0].ID != "f2" { // f2 is newest
			t.Errorf("Expected first item to be f2 (Newest), got %s", retrieved[0].ID)
		}
		if retrieved[1].ID != "f1" {
			t.Errorf("Expected second item to be f1 (Oldest), got %s", retrieved[1].ID)
		}
	}
}

func TestGetFeedbackCounts(t *testing.T) {
	repo, _ := newTestRepo(t)
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

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

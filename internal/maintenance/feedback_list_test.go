package maintenance

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

type feedbackTestCase struct {
	verifyFields func(t *testing.T, feedbacks []domain.Feedback)
	name         string
	dbPath       string
	wantCount    int
	wantErr      bool
}

func setupPopulatedDb(t *testing.T, dbPath string) {
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create populated db: %v", err)
	}

	fb1 := domain.Feedback{
		ID:          "fb-1",
		UserID:      "user-1",
		Type:        domain.FeedbackTypeIssue,
		Content:     "First feedback content",
		Fingerprint: "fingerprint-1",
		Resolved:    false,
		CreatedAt:   time.Now().Add(-1 * time.Hour),
	}
	fb2 := domain.Feedback{
		ID:          "fb-2",
		UserID:      "",
		Type:        domain.FeedbackTypeFeature,
		Content:     "Second feedback content",
		Fingerprint: "fingerprint-2",
		Resolved:    true,
		CreatedAt:   time.Now(),
	}

	if err := repo.SaveFeedback(context.Background(), fb1); err != nil {
		t.Fatalf("failed to save feedback 1: %v", err)
	}
	if err := repo.SaveFeedback(context.Background(), fb2); err != nil {
		t.Fatalf("failed to save feedback 2: %v", err)
	}
	_ = repo.Close()
}

func verifyPopulatedDbFields(t *testing.T, feedbacks []domain.Feedback) {
	if feedbacks[0].ID != "fb-2" {
		t.Errorf("expected first feedback to be fb-2, got %s", feedbacks[0].ID)
	}
	if !feedbacks[0].Resolved {
		t.Errorf("expected fb-2 to be resolved")
	}
	if feedbacks[1].ID != "fb-1" {
		t.Errorf("expected second feedback to be fb-1, got %s", feedbacks[1].ID)
	}
	if feedbacks[1].Resolved {
		t.Errorf("expected fb-1 to not be resolved")
	}
	if feedbacks[1].Fingerprint != "fingerprint-1" {
		t.Errorf("expected fb-1 fingerprint to be fingerprint-1, got %s", feedbacks[1].Fingerprint)
	}
}

func runFeedbackTestCase(t *testing.T, tt feedbackTestCase) {
	got, err := CheckFeedbackList(tt.dbPath)
	if (err != nil) != tt.wantErr {
		t.Fatalf("CheckFeedbackList() error = %v, wantErr %v", err, tt.wantErr)
	}
	if len(got) != tt.wantCount {
		t.Errorf("CheckFeedbackList() got %d feedbacks, want %d", len(got), tt.wantCount)
	}
	if tt.verifyFields != nil {
		tt.verifyFields(t, got)
	}
}

func TestCheckFeedbackList(t *testing.T) {
	tmpDir := t.TempDir()

	missingDbPath := filepath.Join(tmpDir, "missing.db")

	emptyDbPath := filepath.Join(tmpDir, "empty.db")
	repoEmpty, err := sqlite.NewSQLiteRepository(emptyDbPath)
	if err != nil {
		t.Fatalf("failed to create empty db: %v", err)
	}
	_ = repoEmpty.Close()

	populatedDbPath := filepath.Join(tmpDir, "populated.db")
	setupPopulatedDb(t, populatedDbPath)

	tests := []feedbackTestCase{
		{
			name:      "missing db returns error",
			dbPath:    missingDbPath,
			wantErr:   true,
			wantCount: 0,
		},
		{
			name:      "empty db returns no feedback",
			dbPath:    emptyDbPath,
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:         "populated db returns feedbacks in correct order",
			dbPath:       populatedDbPath,
			wantErr:      false,
			wantCount:    2,
			verifyFields: verifyPopulatedDbFields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runFeedbackTestCase(t, tt)
		})
	}
}

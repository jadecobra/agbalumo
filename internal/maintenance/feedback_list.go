package maintenance

import (
	"context"
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

// CheckFeedbackList retrieves all feedback submissions from the database path.
func CheckFeedbackList(dbPath string) ([]domain.Feedback, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("database file does not exist: %s", dbPath)
	}

	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}
	defer func() { _ = repo.Close() }()

	feedbacks, err := repo.GetAllFeedback(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve feedback: %w", err)
	}

	return feedbacks, nil
}

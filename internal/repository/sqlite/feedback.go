package sqlite

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// SaveFeedback saves a user feedback entry.
func (r *SQLiteRepository) SaveFeedback(ctx context.Context, f domain.Feedback) error {
	query := `
	INSERT INTO feedback (id, user_id, type, content, created_at)
	VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		f.ID, f.UserID, f.Type, f.Content, f.CreatedAt,
	)
	return err
}

// GetAllFeedback retrieves all feedback entries ordered by creation time (newest first).
func (r *SQLiteRepository) GetAllFeedback(ctx context.Context) ([]domain.Feedback, error) {
	query := `SELECT id, user_id, type, content, created_at FROM feedback ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []domain.Feedback
	for rows.Next() {
		var f domain.Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.Type, &f.Content, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, rows.Err()
}

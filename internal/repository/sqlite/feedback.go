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

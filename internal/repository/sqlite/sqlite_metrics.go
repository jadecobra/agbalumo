package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
)

func (r *SQLiteRepository) SaveMetric(ctx context.Context, m domain.Metric) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	query := `INSERT INTO metrics (id, event_type, value, metadata, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.writeDB.ExecContext(ctx, query, m.ID, m.EventType, m.Value, m.Metadata, m.CreatedAt)
	return err
}

func (r *SQLiteRepository) GetAverageValue(ctx context.Context, eventType string, since time.Time) (float64, error) {
	query := `SELECT AVG(value) FROM metrics WHERE event_type = ? AND created_at >= ?`
	var avg sql.NullFloat64
	err := r.readDB.QueryRowContext(ctx, query, eventType, since).Scan(&avg)
	if err != nil {
		return 0, err
	}
	if !avg.Valid {
		return 0, nil
	}
	return avg.Float64, nil
}

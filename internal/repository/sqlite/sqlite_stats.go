package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func (r *SQLiteRepository) GetFeedbackCounts(ctx context.Context) (map[domain.FeedbackType]int, error) {
	query := `SELECT type, COUNT(*) FROM feedback GROUP BY type`
	rows, err := r.readDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	counts := make(map[domain.FeedbackType]int)
	for rows.Next() {
		var t domain.FeedbackType
		var count int
		if err := rows.Scan(&t, &count); err != nil {
			return nil, err
		}
		counts[t] = count
	}
	return counts, rows.Err()
}

func (r *SQLiteRepository) queryDailyMetrics(ctx context.Context, query string) ([]domain.DailyMetric, error) {
	rows, err := r.readDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	metrics := []domain.DailyMetric{}
	for rows.Next() {
		var m domain.DailyMetric
		var day sql.NullString
		if err := rows.Scan(&day, &m.Count); err != nil {
			return nil, err
		}
		if day.Valid {
			m.Date = day.String
			metrics = append(metrics, m)
		}
	}
	return metrics, nil
}

// dailyGrowthSQL returns the SQL query for counting new rows per day over the last 30 days.
// table is always a compile-time constant ("listings" or "users") — never user input.
func dailyGrowthSQL(table string) string { //nolint:gosec // table is always a compile-time constant ("listings" or "users"), never user input
	return fmt.Sprintf(`
		SELECT date(created_at) as day, COUNT(*) as count
		FROM %s
		WHERE created_at IS NOT NULL AND created_at != '' AND created_at >= date('now', '-30 days')
		GROUP BY day
		ORDER BY day ASC
	`, table)
}

// GetListingGrowth returns the count of new listings per day for the last 30 days.
func (r *SQLiteRepository) GetListingGrowth(ctx context.Context) ([]domain.DailyMetric, error) {
	return r.queryDailyMetrics(ctx, dailyGrowthSQL("listings"))
}

// GetUserGrowth returns the count of new users per day for the last 30 days.
func (r *SQLiteRepository) GetUserGrowth(ctx context.Context) ([]domain.DailyMetric, error) {
	return r.queryDailyMetrics(ctx, dailyGrowthSQL("users"))
}

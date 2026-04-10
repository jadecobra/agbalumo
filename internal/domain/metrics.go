package domain

import (
	"context"
	"time"
)

type Metric struct {
	ID        string    `json:"id"`
	EventType string    `json:"event_type"`
	Value     float64   `json:"value"`
	Metadata  string    `json:"metadata"` // JSON blob
	CreatedAt time.Time `json:"created_at"`
}

type MetricsRepository interface {
	SaveMetric(ctx context.Context, m Metric) error
	GetAverageValue(ctx context.Context, eventType string, since time.Time) (float64, error)
	GetDailyMetrics(ctx context.Context, eventType string, days int) ([]DailyMetric, error)
}

type MetricsService interface {
	LogAndSave(ctx context.Context, eventType string, value float64, metadata map[string]interface{})
}

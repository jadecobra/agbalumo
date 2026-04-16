package domain

import (
	"context"
	"time"
)

type Metric struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	EventType string    `json:"event_type"`
	Metadata  string    `json:"metadata"`
	Value     float64   `json:"value"`
}

type MetricsRepository interface {
	SaveMetric(ctx context.Context, m Metric) error
	GetAverageValue(ctx context.Context, eventType string, since time.Time) (float64, error)
	GetDailyMetrics(ctx context.Context, eventType string, days int) ([]DailyMetric, error)
}

type MetricsService interface {
	LogAndSave(ctx context.Context, eventType string, value float64, metadata map[string]interface{})
}

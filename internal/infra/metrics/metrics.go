package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jadecobra/agbalumo/internal/domain"
)

const MetricPrefix = "[ADA-METRIC]"

type Service struct {
	repo   domain.ListingRepository
	logger *slog.Logger
}

func NewService(repo domain.ListingRepository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// LogAndSave captures a metric, persists it to the DB, and logs it to stdout.
func (s *Service) LogAndSave(ctx context.Context, eventType string, value float64, metadata map[string]interface{}) {
	m := domain.Metric{
		EventType: eventType,
		Value:     value,
	}

	if metadata != nil {
		if detailJSON, err := json.Marshal(metadata); err == nil {
			m.Metadata = string(detailJSON)
		}
	}

	// Persistent Save
	if err := s.repo.SaveMetric(ctx, m); err != nil {
		s.logger.Error("failed to save metric to DB", "error", err, "event", eventType)
	}

	// Log to stdout for real-time visibility and external log aggregation
	s.logger.Info(fmt.Sprintf("%s %s Value: %.2f Extras: %s", MetricPrefix, eventType, value, m.Metadata))
}

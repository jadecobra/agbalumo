package sqlite

import (
	"context"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMetricsLifecycle(t *testing.T) {
	t.Parallel()
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}
	defer func() { _ = repo.Close() }()

	ctx := context.Background()

	// 1. Save Metric
	m := domain.Metric{
		ID:        "m1",
		EventType: "test_event",
		Value:     10.5,
		Metadata:  `{"id": "test-1"}`,
	}
	err = repo.SaveMetric(ctx, m)
	assert.NoError(t, err)

	// 2. Save Another
	m2 := domain.Metric{
		ID:        "m2",
		EventType: "test_event",
		Value:     20.5,
		Metadata:  `{"id": "test-2"}`,
	}
	_ = repo.SaveMetric(ctx, m2)

	// 3. Get Average
	avg, err := repo.GetAverageValue(ctx, "test_event", time.Now().Add(-1*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, 15.5, avg)

	// 4. Get Average - No metrics
	avg2, err := repo.GetAverageValue(ctx, "non_existent", time.Now().Add(-1*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, 0.0, avg2)
}

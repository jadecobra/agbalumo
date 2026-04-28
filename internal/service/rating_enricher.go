package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

type RatingEnricherJob struct {
	repo         domain.ListingRepository
	placesClient *GooglePlacesClient
}

func NewRatingEnricherJob(repo domain.ListingRepository, client *GooglePlacesClient) *RatingEnricherJob {
	return &RatingEnricherJob{
		repo:         repo,
		placesClient: client,
	}
}

func (j *RatingEnricherJob) EnrichRatings(ctx context.Context, limit int) (int, error) {
	targets, err := j.repo.FindRatingBackfillTargets(ctx, limit)
	if err != nil {
		return 0, err
	}

	successCount := 0
	for _, l := range targets {
		metrics, err := j.placesClient.FetchMetrics(ctx, l.Title, l.City)
		now := time.Now()
		l.RatingUpdatedAt = &now

		if err != nil {
			slog.Warn("[RatingEnricherJob] Failed to fetch Places API metrics", slog.String("id", l.ID), slog.String("title", l.Title), slog.Any("error", err))
			// Still save to set RatingUpdatedAt so we don't spam retry until the 30-day window passes
			_ = j.repo.Save(ctx, l)
			continue
		}

		l.Rating = metrics.Rating
		l.ReviewCount = metrics.ReviewCount

		if err := j.repo.Save(ctx, l); err != nil {
			slog.Error("[RatingEnricherJob] Failed to save updated rating", slog.String("id", l.ID), slog.Any("error", err))
			continue
		}
		successCount++

		// Brief sleep to avoid rapid API depletion
		time.Sleep(500 * time.Millisecond)
	}

	return successCount, nil
}

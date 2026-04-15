package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

type BackgroundService struct {
	Repo     domain.ListingExpirer
	Scraper  *ScraperJob
	Interval time.Duration
}

func NewBackgroundService(repo domain.ListingExpirer, scraper *ScraperJob) *BackgroundService {
	return &BackgroundService{
		Repo:     repo,
		Scraper:  scraper,
		Interval: 1 * time.Hour, // Default
	}
}

// StartTicker runs the expiration logic periodically.
// It blocks, so it should be run in a goroutine.
func (s *BackgroundService) StartTicker(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	slog.Info("[Background] Service started. Ticking every 1 hour.")

	// Run once immediately on start
	s.expireListings(ctx)
	s.enrichListings(ctx)

	for {
		select {
		case <-ticker.C:
			s.expireListings(ctx)
			s.enrichListings(ctx)
		case <-ctx.Done():
			slog.Info("[Background] Service stopping...")
			return
		}
	}
}

func (s *BackgroundService) expireListings(ctx context.Context) {
	count, err := s.Repo.ExpireListings(ctx)
	if err != nil {
		slog.Error("[Background] Error expiring listings", "error", err)
		return
	}
	if count > 0 {
		slog.Info("[Background] Expired listings", "count", count)
	}
}

func (s *BackgroundService) enrichListings(ctx context.Context) {
	if s.Scraper == nil {
		return
	}
	// Enrich up to 5 listings per tick to avoid rate limiting
	count, err := s.Scraper.EnrichListings(ctx, 5)
	if err != nil {
		slog.Error("[Background] Error enriching listings", "error", err)
		return
	}
	if count > 0 {
		slog.Info("[Background] Enriched listings", "count", count)
	}
}

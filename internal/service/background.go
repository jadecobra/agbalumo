package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

type BackgroundService struct {
	Repo domain.ListingExpirer
}

func NewBackgroundService(repo domain.ListingExpirer) *BackgroundService {
	return &BackgroundService{Repo: repo}
}

// StartTicker runs the expiration logic periodically.
// It blocks, so it should be run in a goroutine.
func (s *BackgroundService) StartTicker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour) // Run every hour
	defer ticker.Stop()

	slog.Info("[Background] Service started. Ticking every 1 hour.")

	// Run once immediately on start
	s.expireListings(ctx)

	for {
		select {
		case <-ticker.C:
			s.expireListings(ctx)
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

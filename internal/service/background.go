package service

import (
	"context"
	"log"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

type BackgroundService struct {
	Repo domain.ListingRepository
}

func NewBackgroundService(repo domain.ListingRepository) *BackgroundService {
	return &BackgroundService{Repo: repo}
}

// StartTicker runs the expiration logic periodically.
// It blocks, so it should be run in a goroutine.
func (s *BackgroundService) StartTicker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour) // Run every hour
	defer ticker.Stop()

	log.Println("[Background] Service started. Ticking every 1 hour.")

	// Run once immediately on start
	s.expireListings(ctx)

	for {
		select {
		case <-ticker.C:
			s.expireListings(ctx)
		case <-ctx.Done():
			log.Println("[Background] Service stopping...")
			return
		}
	}
}

func (s *BackgroundService) expireListings(ctx context.Context) {
	count, err := s.Repo.ExpireListings(ctx)
	if err != nil {
		log.Printf("[Background] Error expiring listings: %v", err)
		return
	}
	if count > 0 {
		log.Printf("[Background] Expired %d listings", count)
	}
}

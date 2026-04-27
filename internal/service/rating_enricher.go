package service

import (
	"context"
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
	return 0, nil
}

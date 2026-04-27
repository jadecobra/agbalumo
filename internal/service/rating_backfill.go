package service

import (
	"context"
	"github.com/jadecobra/agbalumo/internal/domain"
)

type RatingBackfiller struct {
	repo         domain.ListingRepository
	placesClient *GooglePlacesClient
}

func NewRatingBackfiller(repo domain.ListingRepository, client *GooglePlacesClient) *RatingBackfiller {
	return &RatingBackfiller{repo: repo, placesClient: client}
}

func (b *RatingBackfiller) Backfill(ctx context.Context) (int, error) {
	return 0, nil
}

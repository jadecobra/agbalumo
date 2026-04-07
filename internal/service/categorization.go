package service

import (
	"context"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

type categorizationService struct {
	db    domain.ListingRepository
	cache *domain.CategoryCache
}

func NewCategorizationService(db domain.ListingRepository, cache *domain.CategoryCache) domain.CategorizationService {
	return &categorizationService{
		db:    db,
		cache: cache,
	}
}

func (s *categorizationService) GetActiveCategories(ctx context.Context) ([]domain.CategoryData, error) {
	return s.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
}

func (s *categorizationService) GetCategories(ctx context.Context, filter domain.CategoryFilter) ([]domain.CategoryData, error) {
	// Only cache active categories for now (as per original logic in handlers)
	if filter.ActiveOnly && s.cache != nil {
		if cats, ok := s.cache.Get(); ok {
			return cats, nil
		}
	}

	cats, err := s.db.GetCategories(ctx, filter)
	if err != nil {
		return nil, err
	}

	if filter.ActiveOnly && s.cache != nil {
		s.cache.Set(cats, 5*time.Minute)
	}

	return cats, nil
}


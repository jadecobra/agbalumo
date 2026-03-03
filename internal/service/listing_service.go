package service

import (
	"context"
	"errors"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// ListingServiceRepo is the required repository interface for the listing service.
type ListingServiceRepo interface {
	domain.ListingStore
	domain.CategoryStore
}

// ListingService encapsulates business logic for listing operations.
type ListingService struct {
	Repo ListingServiceRepo
}

// NewListingService creates a new ListingService.
func NewListingService(repo ListingServiceRepo) *ListingService {
	return &ListingService{Repo: repo}
}

// ClaimListing assigns ownership of an unclaimed listing to the given user.
// It validates that the listing exists, is unclaimed, and is a claimable type.
func (s *ListingService) ClaimListing(ctx context.Context, userID, listingID string) (domain.Listing, error) {
	if userID == "" {
		return domain.Listing{}, errors.New("user ID is required")
	}

	listing, err := s.Repo.FindByID(ctx, listingID)
	if err != nil {
		return domain.Listing{}, errors.New("listing not found")
	}

	if listing.OwnerID != "" {
		return domain.Listing{}, errors.New("listing is already owned")
	}

	categoryInfo, err := s.Repo.GetCategory(ctx, string(listing.Type))
	if err != nil {
		return domain.Listing{}, errors.New("invalid category type")
	}

	if !categoryInfo.Claimable {
		return domain.Listing{}, errors.New("listing type cannot be claimed")
	}

	listing.OwnerID = userID
	if err := s.Repo.Save(ctx, listing); err != nil {
		return domain.Listing{}, errors.New("failed to save listing")
	}

	return listing, nil
}

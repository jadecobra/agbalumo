package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
)

// ListingService encapsulates business logic for listing operations.
type ListingService struct {
	ListingStore      domain.ListingStore
	CategoryStore     domain.CategoryStore
	ClaimRequestStore domain.ClaimRequestStore
}

// NewListingService creates a new ListingService.
func NewListingService(
	listingStore domain.ListingStore,
	categoryStore domain.CategoryStore,
	claimRequestStore domain.ClaimRequestStore,
) *ListingService {
	return &ListingService{
		ListingStore:      listingStore,
		CategoryStore:     categoryStore,
		ClaimRequestStore: claimRequestStore,
	}
}

// ClaimListing creates a pending claim request for an unclaimed, claimable listing.
// It validates that the listing exists, is unclaimed, is a claimable type, and that
// the user does not already have a pending claim for this listing.
func (s *ListingService) ClaimListing(ctx context.Context, user domain.User, listingID string) (domain.ClaimRequest, error) {
	if user.ID == "" {
		return domain.ClaimRequest{}, errors.New("user ID is required")
	}

	listing, err := s.ListingStore.FindByID(ctx, listingID)
	if err != nil {
		return domain.ClaimRequest{}, errors.New("listing not found")
	}

	if listing.OwnerID != "" {
		return domain.ClaimRequest{}, errors.New("listing is already owned")
	}

	categoryInfo, err := s.CategoryStore.GetCategory(ctx, string(listing.Type))
	if err != nil {
		return domain.ClaimRequest{}, errors.New("invalid category type")
	}

	if !categoryInfo.Claimable {
		return domain.ClaimRequest{}, errors.New("listing type cannot be claimed")
	}

	// Check for an existing pending claim from this user
	existing, err := s.ClaimRequestStore.GetClaimRequestByUserAndListing(ctx, user.ID, listingID)
	if err == nil && existing.Status == domain.ClaimStatusPending {
		return domain.ClaimRequest{}, errors.New("you already have a pending claim for this listing")
	}

	cr := domain.ClaimRequest{
		ID:           uuid.New().String(),
		ListingID:    listingID,
		ListingTitle: listing.Title,
		UserID:       user.ID,
		UserName:     user.Name,
		UserEmail:    user.Email,
		Status:       domain.ClaimStatusPending,
		CreatedAt:    time.Now(),
	}

	if err := s.ClaimRequestStore.SaveClaimRequest(ctx, cr); err != nil {
		return domain.ClaimRequest{}, errors.New("failed to save claim request")
	}

	return cr, nil
}

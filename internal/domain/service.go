package domain

import "context"

// ListingService encapsulates business logic associated with Listings across modules.
type ListingService interface {
	ClaimListing(ctx context.Context, user User, listingID string) (ClaimRequest, error)
}

// CategorizationService handles category management and caching.
type CategorizationService interface {
	GetActiveCategories(ctx context.Context) ([]CategoryData, error)
	GetCategories(ctx context.Context, filter CategoryFilter) ([]CategoryData, error)
}



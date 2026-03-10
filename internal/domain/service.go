package domain

import "context"

// ListingService encapsulates business logic associated with Listings across modules.
type ListingService interface {
	ClaimListing(ctx context.Context, user User, listingID string) (ClaimRequest, error)
}

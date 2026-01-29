package domain

import "context"

// ListingRepository defines the interface for persisting and retrieving listings.
type ListingRepository interface {
	Save(ctx context.Context, listing Listing) error
	FindAll(ctx context.Context, filterType string) ([]Listing, error)
	FindByID(ctx context.Context, id string) (Listing, error)
}

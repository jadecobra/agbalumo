package domain

import "context"

// ListingRepository defines the interface for persisting and retrieving listings.
type ListingRepository interface {
	Save(ctx context.Context, listing Listing) error
	FindAll(ctx context.Context, filterType string, queryText string, includeInactive bool) ([]Listing, error)
	FindByID(ctx context.Context, id string) (Listing, error)

	// User Methods
	SaveUser(ctx context.Context, user User) error
	FindUserByGoogleID(ctx context.Context, googleID string) (User, error)
	FindUserByID(ctx context.Context, id string) (User, error)
}

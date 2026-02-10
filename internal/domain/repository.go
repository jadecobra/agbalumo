package domain

import "context"

// ListingRepository defines the interface for persisting and retrieving listings.
type ListingRepository interface {
	Save(ctx context.Context, listing Listing) error
	FindAll(ctx context.Context, filterType string, queryText string, includeInactive bool) ([]Listing, error)
	FindByID(ctx context.Context, id string) (Listing, error)
	FindAllByOwner(ctx context.Context, ownerID string) ([]Listing, error)
	Delete(ctx context.Context, id string) error
	GetCounts(ctx context.Context) (map[Category]int, error)

	// User Methods
	SaveUser(ctx context.Context, user User) error
	FindUserByGoogleID(ctx context.Context, googleID string) (User, error)
	FindUserByID(ctx context.Context, id string) (User, error)
	
	// Maintenance
	ExpireListings(ctx context.Context) (int64, error)

	// Feedback
	SaveFeedback(ctx context.Context, feedback Feedback) error
}

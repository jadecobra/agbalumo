package domain

import "context"

// --- Focused Interfaces ---

// ListingStore handles core listing CRUD and query operations.
type ListingStore interface {
	Save(ctx context.Context, listing Listing) error
	FindAll(ctx context.Context, filterType string, queryText string, includeInactive bool, limit int, offset int) ([]Listing, error)
	FindByID(ctx context.Context, id string) (Listing, error)
	FindAllByOwner(ctx context.Context, ownerID string) ([]Listing, error)
	Delete(ctx context.Context, id string) error
	GetCounts(ctx context.Context) (map[Category]int, error)
	GetFeaturedListings(ctx context.Context) ([]Listing, error)
}

// ListingSaver is the minimal interface for saving a listing.
type ListingSaver interface {
	Save(ctx context.Context, listing Listing) error
}

// ListingExpirer handles expiration of stale listings.
type ListingExpirer interface {
	ExpireListings(ctx context.Context) (int64, error)
}

// UserStore handles user persistence and lookup.
type UserStore interface {
	SaveUser(ctx context.Context, user User) error
	FindUserByGoogleID(ctx context.Context, googleID string) (User, error)
	FindUserByID(ctx context.Context, id string) (User, error)
}

// FeedbackStore handles feedback persistence.
type FeedbackStore interface {
	SaveFeedback(ctx context.Context, feedback Feedback) error
	GetAllFeedback(ctx context.Context) ([]Feedback, error)
	GetFeedbackCounts(ctx context.Context) (map[FeedbackType]int, error)
}

// AdminStore handles admin-specific queries.
type AdminStore interface {
	GetPendingListings(ctx context.Context, limit int, offset int) ([]Listing, error)
	GetUserCount(ctx context.Context) (int, error)
	GetAllUsers(ctx context.Context) ([]User, error)
}

// AnalyticsStore handles growth/analytics queries.
type AnalyticsStore interface {
	GetListingGrowth(ctx context.Context) ([]DailyMetric, error)
	GetUserGrowth(ctx context.Context) ([]DailyMetric, error)
}

// --- Composed Super-Interface (Backward Compatible) ---

// ListingRepository composes all store interfaces into a single contract.
// Consumers should prefer focused interfaces where possible.
type ListingRepository interface {
	ListingStore
	ListingExpirer
	UserStore
	FeedbackStore
	AdminStore
	AnalyticsStore
}

// DailyMetric represents a daily count of an entity.
type DailyMetric struct {
	Date  string
	Count int
}

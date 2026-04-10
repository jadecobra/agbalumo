package domain

import (
	"context"
	"time"
)

// --- Focused Interfaces ---

// ListingReader handles read-only queries for listings.
type ListingReader interface {
	FindAll(ctx context.Context, filterType string, queryText string, sortField string, sortOrder string, includeInactive bool, limit int, offset int) ([]Listing, int, error)
	FindByID(ctx context.Context, id string) (Listing, error)
	FindByTitle(ctx context.Context, title string) ([]Listing, error)
	TitleExists(ctx context.Context, title string) (bool, error)
	FindAllByOwner(ctx context.Context, ownerID string, limit int, offset int) ([]Listing, int, error)
	GetCounts(ctx context.Context) (map[Category]int, error)
	GetLocations(ctx context.Context) ([]string, error)
	GetFeaturedListings(ctx context.Context, category string) ([]Listing, error)
}

// ListingWriter handles write operations for listings.
type ListingWriter interface {
	Save(ctx context.Context, listing Listing) error
	Delete(ctx context.Context, id string) error
	SetFeatured(ctx context.Context, id string, featured bool) error
}

// ListingStore handles core listing CRUD and query operations by composing
// reader and writer interfaces (CQRS).
type ListingStore interface {
	ListingReader
	ListingWriter
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
	GetUserCount(ctx context.Context) (int, error)
	GetAllUsers(ctx context.Context, limit int, offset int) ([]User, error)
}

// ClaimRequestStore handles claim request persistence.
type ClaimRequestStore interface {
	SaveClaimRequest(ctx context.Context, r ClaimRequest) error
	GetPendingClaimRequests(ctx context.Context) ([]ClaimRequest, error)
	UpdateClaimRequestStatus(ctx context.Context, id string, status ClaimStatus) error
	GetClaimRequestByUserAndListing(ctx context.Context, userID, listingID string) (ClaimRequest, error)
}

// AnalyticsStore handles growth/analytics queries.
type AnalyticsStore interface {
	GetListingGrowth(ctx context.Context) ([]DailyMetric, error)
	GetUserGrowth(ctx context.Context) ([]DailyMetric, error)
	SaveMetric(ctx context.Context, m Metric) error
	GetAverageValue(ctx context.Context, eventType string, since time.Time) (float64, error)
}

// CategoryStore handles category persistence and retrieval.
type CategoryStore interface {
	GetCategories(ctx context.Context, filter CategoryFilter) ([]CategoryData, error)
	GetCategory(ctx context.Context, name string) (CategoryData, error)
	SaveCategory(ctx context.Context, c CategoryData) error
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
	CategoryStore
	ClaimRequestStore
}

// DailyMetric represents a daily count of an entity.
type DailyMetric struct {
	Date  string
	Count int
}

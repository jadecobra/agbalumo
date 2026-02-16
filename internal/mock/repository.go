package mock

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// MockListingRepository is a mock implementation of domain.ListingRepository
type MockListingRepository struct {
	SaveFn     func(ctx context.Context, listing domain.Listing) error
	FindAllFn  func(ctx context.Context, filterType, queryText string, includeInactive bool) ([]domain.Listing, error) // Updated signature
	FindByIDFn func(ctx context.Context, id string) (domain.Listing, error)

	// User Mocks
	SaveUserFn           func(ctx context.Context, u domain.User) error
	FindUserByGoogleIDFn func(ctx context.Context, googleID string) (domain.User, error)
	FindUserByIDFn       func(ctx context.Context, id string) (domain.User, error)
	FindAllByOwnerFn     func(ctx context.Context, ownerID string) ([]domain.Listing, error)
	DeleteFn             func(ctx context.Context, id string) error
	GetCountsFn          func(ctx context.Context) (map[domain.Category]int, error)
	ExpireListingsFn     func(ctx context.Context) (int64, error)
	SaveFeedbackFn       func(ctx context.Context, feedback domain.Feedback) error
	GetAllFeedbackFn     func(ctx context.Context) ([]domain.Feedback, error)

	// New Admin Methods
	GetPendingListingsFn func(ctx context.Context) ([]domain.Listing, error)
	GetUserCountFn       func(ctx context.Context) (int, error)
	GetFeedbackCountsFn  func(ctx context.Context) (map[domain.FeedbackType]int, error)
}

func (m *MockListingRepository) Save(ctx context.Context, l domain.Listing) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, l)
	}
	return nil
}

func (m *MockListingRepository) FindAll(ctx context.Context, filterType, queryText string, includeInactive bool) ([]domain.Listing, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx, filterType, queryText, includeInactive)
	}
	return nil, nil
}

func (m *MockListingRepository) FindByID(ctx context.Context, id string) (domain.Listing, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return domain.Listing{}, nil
}

func (m *MockListingRepository) SaveUser(ctx context.Context, u domain.User) error {
	if m.SaveUserFn != nil {
		return m.SaveUserFn(ctx, u)
	}
	return nil
}

func (m *MockListingRepository) FindUserByGoogleID(ctx context.Context, googleID string) (domain.User, error) {
	if m.FindUserByGoogleIDFn != nil {
		return m.FindUserByGoogleIDFn(ctx, googleID)
	}
	return domain.User{}, nil
}

func (m *MockListingRepository) FindUserByID(ctx context.Context, id string) (domain.User, error) {
	if m.FindUserByIDFn != nil {
		return m.FindUserByIDFn(ctx, id)
	}
	return domain.User{}, nil
}

func (m *MockListingRepository) FindAllByOwner(ctx context.Context, ownerID string) ([]domain.Listing, error) {
	if m.FindAllByOwnerFn != nil {
		return m.FindAllByOwnerFn(ctx, ownerID)
	}
	return nil, nil
}

func (m *MockListingRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockListingRepository) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	if m.GetCountsFn != nil {
		return m.GetCountsFn(ctx)
	}
	return nil, nil
}

func (m *MockListingRepository) ExpireListings(ctx context.Context) (int64, error) {
	if m.ExpireListingsFn != nil {
		return m.ExpireListingsFn(ctx)
	}
	return 0, nil
}

func (m *MockListingRepository) SaveFeedback(ctx context.Context, f domain.Feedback) error {
	if m.SaveFeedbackFn != nil {
		return m.SaveFeedbackFn(ctx, f)
	}
	return nil
}

func (m *MockListingRepository) GetAllFeedback(ctx context.Context) ([]domain.Feedback, error) {
	if m.GetAllFeedbackFn != nil {
		return m.GetAllFeedbackFn(ctx)
	}
	return nil, nil
}

func (m *MockListingRepository) GetFeedbackCounts(ctx context.Context) (map[domain.FeedbackType]int, error) {
	if m.GetFeedbackCountsFn != nil {
		return m.GetFeedbackCountsFn(ctx)
	}
	return nil, nil
}

func (m *MockListingRepository) GetPendingListings(ctx context.Context) ([]domain.Listing, error) {
	if m.GetPendingListingsFn != nil {
		return m.GetPendingListingsFn(ctx)
	}
	return nil, nil
}

func (m *MockListingRepository) GetUserCount(ctx context.Context) (int, error) {
	if m.GetUserCountFn != nil {
		return m.GetUserCountFn(ctx)
	}
	return 0, nil
}

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


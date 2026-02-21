package mock

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/domain"
	testifyMock "github.com/stretchr/testify/mock"
)

// MockListingRepository is a mock implementation of domain.ListingRepository
// MockListingRepository is a mock implementation of domain.ListingRepository
type MockListingRepository struct {
	testifyMock.Mock
}

func (m *MockListingRepository) Save(ctx context.Context, l domain.Listing) error {
	args := m.Called(ctx, l)
	return args.Error(0)
}

func (m *MockListingRepository) FindAll(ctx context.Context, filterType, queryText string, includeInactive bool, limit int, offset int) ([]domain.Listing, error) {
	args := m.Called(ctx, filterType, queryText, includeInactive, limit, offset)
	return args.Get(0).([]domain.Listing), args.Error(1)
}

func (m *MockListingRepository) FindByID(ctx context.Context, id string) (domain.Listing, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Listing), args.Error(1)
}

func (m *MockListingRepository) FindByTitle(ctx context.Context, title string) ([]domain.Listing, error) {
	args := m.Called(ctx, title)
	return args.Get(0).([]domain.Listing), args.Error(1)
}

func (m *MockListingRepository) SaveUser(ctx context.Context, u domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockListingRepository) FindUserByGoogleID(ctx context.Context, googleID string) (domain.User, error) {
	args := m.Called(ctx, googleID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockListingRepository) FindUserByID(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockListingRepository) FindAllByOwner(ctx context.Context, ownerID string) ([]domain.Listing, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).([]domain.Listing), args.Error(1)
}

func (m *MockListingRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockListingRepository) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[domain.Category]int), args.Error(1)
}

func (m *MockListingRepository) ExpireListings(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockListingRepository) SaveFeedback(ctx context.Context, f domain.Feedback) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *MockListingRepository) GetAllFeedback(ctx context.Context) ([]domain.Feedback, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Feedback), args.Error(1)
}

func (m *MockListingRepository) GetFeedbackCounts(ctx context.Context) (map[domain.FeedbackType]int, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[domain.FeedbackType]int), args.Error(1)
}

func (m *MockListingRepository) GetPendingListings(ctx context.Context, limit int, offset int) ([]domain.Listing, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Listing), args.Error(1)
}

func (m *MockListingRepository) GetUserCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockListingRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockListingRepository) GetFeaturedListings(ctx context.Context) ([]domain.Listing, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Listing), args.Error(1)
}

func (m *MockListingRepository) GetListingGrowth(ctx context.Context) ([]domain.DailyMetric, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.DailyMetric), args.Error(1)
}

func (m *MockListingRepository) GetUserGrowth(ctx context.Context) ([]domain.DailyMetric, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.DailyMetric), args.Error(1)
}

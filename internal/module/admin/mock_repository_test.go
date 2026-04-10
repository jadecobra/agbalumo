package admin_test

import (
	"context"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// MockListingRepository is a minimalist mock for testing error paths.
type MockListingRepository struct {
	ErrorOn map[string]error
}

func NewMockRepository() *MockListingRepository {
	return &MockListingRepository{
		ErrorOn: make(map[string]error),
	}
}

func (m *MockListingRepository) Save(ctx context.Context, listing domain.Listing) error {
	return m.ErrorOn["Save"]
}

func (m *MockListingRepository) FindAll(ctx context.Context, filterType string, queryText string, sortField string, sortOrder string, includeInactive bool, limit int, offset int) ([]domain.Listing, int, error) {
	return nil, 0, m.ErrorOn["FindAll"]
}

func (m *MockListingRepository) FindByID(ctx context.Context, id string) (domain.Listing, error) {
	return domain.Listing{}, m.ErrorOn["FindByID"]
}

func (m *MockListingRepository) FindByTitle(ctx context.Context, title string) ([]domain.Listing, error) {
	return nil, m.ErrorOn["FindByTitle"]
}

func (m *MockListingRepository) TitleExists(ctx context.Context, title string) (bool, error) {
	return false, m.ErrorOn["TitleExists"]
}

func (m *MockListingRepository) FindAllByOwner(ctx context.Context, ownerID string, limit int, offset int) ([]domain.Listing, int, error) {
	return nil, 0, m.ErrorOn["FindAllByOwner"]
}

func (m *MockListingRepository) Delete(ctx context.Context, id string) error {
	return m.ErrorOn["Delete"]
}

func (m *MockListingRepository) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	return nil, m.ErrorOn["GetCounts"]
}

func (m *MockListingRepository) GetLocations(ctx context.Context) ([]string, error) {
	return nil, m.ErrorOn["GetLocations"]
}

func (m *MockListingRepository) GetFeaturedListings(ctx context.Context, category string) ([]domain.Listing, error) {
	return nil, m.ErrorOn["GetFeaturedListings"]
}

func (m *MockListingRepository) SetFeatured(ctx context.Context, id string, featured bool) error {
	return m.ErrorOn["SetFeatured"]
}

func (m *MockListingRepository) SaveUser(ctx context.Context, user domain.User) error {
	return m.ErrorOn["SaveUser"]
}

func (m *MockListingRepository) FindUserByGoogleID(ctx context.Context, googleID string) (domain.User, error) {
	return domain.User{}, m.ErrorOn["FindUserByGoogleID"]
}

func (m *MockListingRepository) FindUserByID(ctx context.Context, id string) (domain.User, error) {
	return domain.User{}, m.ErrorOn["FindUserByID"]
}

func (m *MockListingRepository) SaveFeedback(ctx context.Context, feedback domain.Feedback) error {
	return m.ErrorOn["SaveFeedback"]
}

func (m *MockListingRepository) GetAllFeedback(ctx context.Context) ([]domain.Feedback, error) {
	return nil, m.ErrorOn["GetAllFeedback"]
}

func (m *MockListingRepository) GetFeedbackCounts(ctx context.Context) (map[domain.FeedbackType]int, error) {
	return nil, m.ErrorOn["GetFeedbackCounts"]
}

func (m *MockListingRepository) GetUserCount(ctx context.Context) (int, error) {
	return 0, m.ErrorOn["GetUserCount"]
}

func (m *MockListingRepository) GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error) {
	return nil, m.ErrorOn["GetAllUsers"]
}

func (m *MockListingRepository) GetListingGrowth(ctx context.Context) ([]domain.DailyMetric, error) {
	return nil, m.ErrorOn["GetListingGrowth"]
}

func (m *MockListingRepository) GetUserGrowth(ctx context.Context) ([]domain.DailyMetric, error) {
	return nil, m.ErrorOn["GetUserGrowth"]
}

func (m *MockListingRepository) GetCategories(ctx context.Context, filter domain.CategoryFilter) ([]domain.CategoryData, error) {
	return nil, m.ErrorOn["GetCategories"]
}

func (m *MockListingRepository) GetCategory(ctx context.Context, name string) (domain.CategoryData, error) {
	return domain.CategoryData{}, m.ErrorOn["GetCategory"]
}

func (m *MockListingRepository) SaveCategory(ctx context.Context, c domain.CategoryData) error {
	return m.ErrorOn["SaveCategory"]
}

func (m *MockListingRepository) SaveClaimRequest(ctx context.Context, r domain.ClaimRequest) error {
	return m.ErrorOn["SaveClaimRequest"]
}

func (m *MockListingRepository) GetPendingClaimRequests(ctx context.Context) ([]domain.ClaimRequest, error) {
	return nil, m.ErrorOn["GetPendingClaimRequests"]
}

func (m *MockListingRepository) UpdateClaimRequestStatus(ctx context.Context, id string, status domain.ClaimStatus) error {
	return m.ErrorOn["UpdateClaimRequestStatus"]
}

func (m *MockListingRepository) GetClaimRequestByUserAndListing(ctx context.Context, userID, listingID string) (domain.ClaimRequest, error) {
	return domain.ClaimRequest{}, m.ErrorOn["GetClaimRequestByUserAndListing"]
}

func (m *MockListingRepository) ExpireListings(ctx context.Context) (int64, error) {
	return 0, m.ErrorOn["ExpireListings"]
}

func (m *MockListingRepository) SaveMetric(ctx context.Context, mt domain.Metric) error {
	return m.ErrorOn["SaveMetric"]
}

func (m *MockListingRepository) GetAverageValue(ctx context.Context, eventType string, since time.Time) (float64, error) {
	return 0, m.ErrorOn["GetAverageValue"]
}

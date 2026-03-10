package mock

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/domain"
	testifyMock "github.com/stretchr/testify/mock"
)

// MockListingService is a mock implementation of domain.ListingService
type MockListingService struct {
	testifyMock.Mock
}

func (m *MockListingService) ClaimListing(ctx context.Context, user domain.User, listingID string) (domain.ClaimRequest, error) {
	args := m.Called(ctx, user, listingID)
	return args.Get(0).(domain.ClaimRequest), args.Error(1)
}

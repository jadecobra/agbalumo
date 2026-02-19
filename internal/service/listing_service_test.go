package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/jadecobra/agbalumo/internal/service"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListingService_ClaimListing(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "loc-123").Return(domain.Listing{ID: "loc-123", Type: domain.Business}, nil)
		mockRepo.On("Save", ctx, testifyMock.MatchedBy(func(l domain.Listing) bool {
			return l.OwnerID == "user-123"
		})).Return(nil)

		listing, err := svc.ClaimListing(ctx, "user-123", "loc-123")
		require.NoError(t, err)
		require.Equal(t, "user-123", listing.OwnerID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("missing user id", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		_, err := svc.ClaimListing(ctx, "", "loc-123")
		require.Error(t, err)
		require.Equal(t, "user ID is required", err.Error())
	})

	t.Run("listing not found", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "bad-id").Return(domain.Listing{}, errors.New("listing not found"))

		_, err := svc.ClaimListing(ctx, "user-123", "bad-id")
		require.Error(t, err)
		require.Equal(t, "listing not found", err.Error())
	})

	t.Run("already owned", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "loc-123").Return(domain.Listing{ID: "loc-123", OwnerID: "someone-else", Type: domain.Business}, nil)

		_, err := svc.ClaimListing(ctx, "user-123", "loc-123")
		require.Error(t, err)
		require.Equal(t, "listing is already owned", err.Error())
	})

	t.Run("unclaimable type", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "loc-123").Return(domain.Listing{ID: "loc-123", Type: domain.Job}, nil)

		_, err := svc.ClaimListing(ctx, "user-123", "loc-123")
		require.Error(t, err)
		require.Equal(t, "listing type cannot be claimed", err.Error())
	})

	t.Run("save fails", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "loc-123").Return(domain.Listing{ID: "loc-123", Type: domain.Business}, nil)
		mockRepo.On("Save", ctx, testifyMock.Anything).Return(errors.New("db error"))

		_, err := svc.ClaimListing(ctx, "user-123", "loc-123")
		require.Error(t, err)
		require.Equal(t, "failed to save listing", err.Error())
	})
}

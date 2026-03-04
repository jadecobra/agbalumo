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

var testUser = domain.User{ID: "user-123", Name: "Test User", Email: "test@example.com"}

func TestListingService_ClaimListing(t *testing.T) {
	ctx := context.Background()

	t.Run("success creates pending claim request", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "loc-123").Return(domain.Listing{ID: "loc-123", Title: "Test Listing", Type: domain.Business}, nil)
		mockRepo.On("GetCategory", ctx, string(domain.Business)).Return(domain.CategoryData{ID: string(domain.Business), Claimable: true}, nil)
		mockRepo.On("GetClaimRequestByUserAndListing", ctx, testUser.ID, "loc-123").Return(domain.ClaimRequest{}, errors.New("not found"))
		mockRepo.On("SaveClaimRequest", ctx, testifyMock.MatchedBy(func(cr domain.ClaimRequest) bool {
			return cr.UserID == testUser.ID &&
				cr.ListingID == "loc-123" &&
				cr.Status == domain.ClaimStatusPending &&
				cr.ListingTitle == "Test Listing" &&
				cr.UserName == testUser.Name &&
				cr.UserEmail == testUser.Email
		})).Return(nil)

		cr, err := svc.ClaimListing(ctx, testUser, "loc-123")
		require.NoError(t, err)
		require.Equal(t, domain.ClaimStatusPending, cr.Status)
		require.Equal(t, testUser.ID, cr.UserID)
		require.NotEmpty(t, cr.ID) // uuid generated

		mockRepo.AssertExpectations(t)
	})

	t.Run("missing user id", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		_, err := svc.ClaimListing(ctx, domain.User{}, "loc-123")
		require.Error(t, err)
		require.Equal(t, "user ID is required", err.Error())
	})

	t.Run("listing not found", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "bad-id").Return(domain.Listing{}, errors.New("listing not found"))

		_, err := svc.ClaimListing(ctx, testUser, "bad-id")
		require.Error(t, err)
		require.Equal(t, "listing not found", err.Error())
	})

	t.Run("already owned", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "loc-123").Return(domain.Listing{ID: "loc-123", OwnerID: "someone-else", Type: domain.Business}, nil)

		_, err := svc.ClaimListing(ctx, testUser, "loc-123")
		require.Error(t, err)
		require.Equal(t, "listing is already owned", err.Error())
	})

	t.Run("unclaimable type", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "loc-123").Return(domain.Listing{ID: "loc-123", Type: domain.Job}, nil)
		mockRepo.On("GetCategory", ctx, string(domain.Job)).Return(domain.CategoryData{ID: string(domain.Job), Claimable: false}, nil)

		_, err := svc.ClaimListing(ctx, testUser, "loc-123")
		require.Error(t, err)
		require.Equal(t, "listing type cannot be claimed", err.Error())
	})

	t.Run("duplicate pending claim rejected", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "loc-123").Return(domain.Listing{ID: "loc-123", Title: "Test", Type: domain.Business}, nil)
		mockRepo.On("GetCategory", ctx, string(domain.Business)).Return(domain.CategoryData{ID: string(domain.Business), Claimable: true}, nil)
		mockRepo.On("GetClaimRequestByUserAndListing", ctx, testUser.ID, "loc-123").Return(
			domain.ClaimRequest{ID: "existing", Status: domain.ClaimStatusPending}, nil,
		)

		_, err := svc.ClaimListing(ctx, testUser, "loc-123")
		require.Error(t, err)
		require.Equal(t, "you already have a pending claim for this listing", err.Error())
	})

	t.Run("save fails", func(t *testing.T) {
		mockRepo := new(mock.MockListingRepository)
		svc := service.NewListingService(mockRepo)

		mockRepo.On("FindByID", ctx, "loc-123").Return(domain.Listing{ID: "loc-123", Title: "Test", Type: domain.Business}, nil)
		mockRepo.On("GetCategory", ctx, string(domain.Business)).Return(domain.CategoryData{ID: string(domain.Business), Claimable: true}, nil)
		mockRepo.On("GetClaimRequestByUserAndListing", ctx, testUser.ID, "loc-123").Return(domain.ClaimRequest{}, errors.New("not found"))
		mockRepo.On("SaveClaimRequest", ctx, testifyMock.Anything).Return(errors.New("db error"))

		_, err := svc.ClaimListing(ctx, testUser, "loc-123")
		require.Error(t, err)
		require.Equal(t, "failed to save claim request", err.Error())
	})
}

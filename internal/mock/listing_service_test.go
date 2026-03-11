package mock_test

import (
	"context"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestMockListingService_ClaimListing(t *testing.T) {
	service := new(mock.MockListingService)
	ctx := context.Background()
	user := domain.User{ID: "u1"}

	service.On("ClaimListing", ctx, user, "l1").Return(domain.ClaimRequest{}, nil)

	_, err := service.ClaimListing(ctx, user, "l1")
	assert.NoError(t, err)
	service.AssertExpectations(t)
}

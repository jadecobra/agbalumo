package handler_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAdminHandler_HandleDashboard_ErrorPaths(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*mock.MockListingRepository)
		expectCode int
	}{
		{
			name: "PendingClaimRequestsError",
			setupMock: func(r *mock.MockListingRepository) {
				r.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, assert.AnError)
			},
			expectCode: http.StatusInternalServerError,
		},
		{
			name: "UserCountError",
			setupMock: func(r *mock.MockListingRepository) {
				r.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
				r.On("GetUserCount", testifyMock.Anything).Return(0, assert.AnError)
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
			c.Set("User", domain.User{Role: domain.UserRoleAdmin})

			mockRepo := &mock.MockListingRepository{}
			mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
			mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, nil).Maybe()
			mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil).Maybe()
			tt.setupMock(mockRepo)

			h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
			_ = h.HandleDashboard(c)
			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

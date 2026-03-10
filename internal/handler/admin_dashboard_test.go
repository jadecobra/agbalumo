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
		{
			name: "FeedbackCountsError",
			setupMock: func(r *mock.MockListingRepository) {
				r.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
				r.On("GetUserCount", testifyMock.Anything).Return(5, nil)
				r.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, assert.AnError)
			},
			expectCode: http.StatusInternalServerError,
		},
		{
			name: "ListingGrowthError",
			setupMock: func(r *mock.MockListingRepository) {
				r.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
				r.On("GetUserCount", testifyMock.Anything).Return(5, nil)
				r.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
				r.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, assert.AnError)
			},
			expectCode: http.StatusInternalServerError,
		},
		{
			name: "UserGrowthError",
			setupMock: func(r *mock.MockListingRepository) {
				r.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
				r.On("GetUserCount", testifyMock.Anything).Return(5, nil)
				r.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
				r.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
				r.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, assert.AnError)
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

func TestAdminHandler_HandleDashboard_HappyPath(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(10, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil)
	mockRepo.On("GetAllUsers", testifyMock.Anything, 10, 0).Return([]domain.User{}, nil)

	h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
	err := h.HandleDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

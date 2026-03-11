package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAdminHandler_HandleAllListings_Extended(t *testing.T) {
	repo := handler.SetupTestRepository(t)

	// Seed listings for filtering and counts
	_ = repo.Save(context.Background(), domain.Listing{ID: "l1", Title: "Test Business", Type: "business", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria", Address: "123 Test St", City: "Lagos"})
	_ = repo.Save(context.Background(), domain.Listing{ID: "l2", Title: "Test Event", Type: "events", Status: domain.ListingStatusApproved, OwnerOrigin: "Nigeria"})

	tests := []struct {
		name       string
		query      string
		expectCode int
	}{
		{
			name:       "HappyPath_WithCategoryFilter",
			query:      "?category=business&sort=title&order=asc",
			expectCode: http.StatusOK,
		},
		{
			name:       "HappyPath_NoFilters",
			query:      "",
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings"+tt.query, nil)
			c.Set("User", domain.User{Role: domain.UserRoleAdmin})

			h := handler.NewAdminHandler(repo, nil, config.LoadConfig())
			_ = h.HandleAllListings(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

func TestAdminHandler_HandleAllListings_Extended_Errors(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		setupMock  func(*mock.MockListingRepository)
		expectCode int
	}{
		{
			name:  "FindAllError_ReturnsError",
			query: "",
			setupMock: func(r *mock.MockListingRepository) {
				r.On("FindAll", testifyMock.Anything, "", "", "", "", true, 50, 0).
					Return([]domain.Listing{}, assert.AnError)
			},
			expectCode: http.StatusInternalServerError,
		},
		{
			name:  "GetCountsError_GracefulFallback",
			query: "",
			setupMock: func(r *mock.MockListingRepository) {
				r.On("FindAll", testifyMock.Anything, "", "", "", "", true, 50, 0).
					Return([]domain.Listing{}, nil)
				r.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, assert.AnError)
				r.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil)
			},
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings"+tt.query, nil)
			c.Set("User", domain.User{Role: domain.UserRoleAdmin})

			mockRepo := &mock.MockListingRepository{}
			tt.setupMock(mockRepo)

			h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
			_ = h.HandleAllListings(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

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

func TestAdminHandler_HandleAllListings_Extended(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		setupMock  func(*mock.MockListingRepository)
		expectCode int
	}{
		{
			name:  "HappyPath_WithCategoryFilter",
			query: "?category=business&sort=title&order=asc",
			setupMock: func(r *mock.MockListingRepository) {
				r.On("FindAll", testifyMock.Anything, "business", "", "title", "ASC", true, 50, 0).
					Return([]domain.Listing{{ID: "l1", Title: "Test"}}, nil)
				r.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{"business": 1}, nil)
				r.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil)
			},
			expectCode: http.StatusOK,
		},
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
			name:  "HappyPath_NoFilters",
			query: "",
			setupMock: func(r *mock.MockListingRepository) {
				r.On("FindAll", testifyMock.Anything, "", "", "", "", true, 50, 0).
					Return([]domain.Listing{}, nil)
				r.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
				r.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil)
			},
			expectCode: http.StatusOK,
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

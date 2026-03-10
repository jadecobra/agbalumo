package handler_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAdminHandler_HandleAllListings(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{{ID: "1"}}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

	h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
	_ = h.HandleAllListings(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_HandleToggleFeatured(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		featured   string
		setupMock  func(*mock.MockListingRepository)
		expectCode int
	}{
		{
			name:     "Success",
			id:       "123",
			featured: "true",
			setupMock: func(m *mock.MockListingRepository) {
				m.On("SetFeatured", testifyMock.Anything, "123", true).Return(nil)
			},
			expectCode: http.StatusOK,
		},
		{
			name:       "MissingID",
			id:         "",
			featured:   "true",
			setupMock:  func(m *mock.MockListingRepository) {},
			expectCode: http.StatusBadRequest,
		},
		{
			name:     "RepoError",
			id:       "123",
			featured: "true",
			setupMock: func(m *mock.MockListingRepository) {
				m.On("SetFeatured", testifyMock.Anything, "123", true).Return(assert.AnError)
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("featured", tt.featured)
			urlPath := "/admin/listings/" + tt.id + "/featured"
			if tt.id == "" {
				urlPath = "/admin/listings/featured" // Simulate routing without ID
			}
			c, rec := setupAdminTestContext(http.MethodPost, urlPath, strings.NewReader(formData.Encode()))
			if tt.id != "" {
				c.SetParamNames("id")
				c.SetParamValues(tt.id)
			}
			c.Set("User", domain.User{Role: domain.UserRoleAdmin})

			mockRepo := &mock.MockListingRepository{}
			tt.setupMock(mockRepo)

			h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
			_ = h.HandleToggleFeatured(c)
			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

func TestAdminHandler_HandleApproveClaim(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/cr1/approve", nil)
	c.SetParamNames("id")
	c.SetParamValues("cr1")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("UpdateClaimRequestStatus", testifyMock.Anything, "cr1", domain.ClaimStatusApproved).Return(nil)

	h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
	_ = h.HandleApproveClaim(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

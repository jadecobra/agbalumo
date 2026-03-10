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
	formData := url.Values{}
	formData.Set("featured", "true")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/123/featured", strings.NewReader(formData.Encode()))
	c.SetParamNames("id")
	c.SetParamValues("123")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("SetFeatured", testifyMock.Anything, "123", true).Return(nil)

	h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
	_ = h.HandleToggleFeatured(c)
	assert.Equal(t, http.StatusOK, rec.Code)
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

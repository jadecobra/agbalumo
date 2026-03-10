package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestListingHandler_HandleClaim(t *testing.T) {
	tests := []struct {
		name       string
		user       *domain.User
		listingID  string
		setupMock  func(*mock.MockListingService)
		expectCode int
	}{
		{
			name:       "NoUser_RedirectsToLogin",
			user:       nil,
			listingID:  "listing1",
			expectCode: http.StatusFound,
		},
		{
			name:      "Success_ReturnsHTML",
			user:      &domain.User{ID: "u1", Name: "Test"},
			listingID: "listing1",
			setupMock: func(s *mock.MockListingService) {
				s.On("ClaimListing", testifyMock.Anything, testifyMock.Anything, "listing1").
					Return(domain.ClaimRequest{}, nil)
			},
			expectCode: http.StatusOK,
		},
		{
			name:      "NotFound_Returns404",
			user:      &domain.User{ID: "u1", Name: "Test"},
			listingID: "listing1",
			setupMock: func(s *mock.MockListingService) {
				s.On("ClaimListing", testifyMock.Anything, testifyMock.Anything, "listing1").
					Return(domain.ClaimRequest{}, errors.New("listing not found"))
			},
			expectCode: http.StatusNotFound,
		},
		{
			name:      "AlreadyOwned_Returns403",
			user:      &domain.User{ID: "u1", Name: "Test"},
			listingID: "listing1",
			setupMock: func(s *mock.MockListingService) {
				s.On("ClaimListing", testifyMock.Anything, testifyMock.Anything, "listing1").
					Return(domain.ClaimRequest{}, errors.New("listing is already owned"))
			},
			expectCode: http.StatusForbidden,
		},
		{
			name:      "NotClaimable_Returns403",
			user:      &domain.User{ID: "u1", Name: "Test"},
			listingID: "listing1",
			setupMock: func(s *mock.MockListingService) {
				s.On("ClaimListing", testifyMock.Anything, testifyMock.Anything, "listing1").
					Return(domain.ClaimRequest{}, errors.New("listing type cannot be claimed"))
			},
			expectCode: http.StatusForbidden,
		},
		{
			name:      "DuplicateClaim_Returns409",
			user:      &domain.User{ID: "u1", Name: "Test"},
			listingID: "listing1",
			setupMock: func(s *mock.MockListingService) {
				s.On("ClaimListing", testifyMock.Anything, testifyMock.Anything, "listing1").
					Return(domain.ClaimRequest{}, errors.New("you already have a pending claim for this listing"))
			},
			expectCode: http.StatusConflict,
		},
		{
			name:      "GenericError_Returns500",
			user:      &domain.User{ID: "u1", Name: "Test"},
			listingID: "listing1",
			setupMock: func(s *mock.MockListingService) {
				s.On("ClaimListing", testifyMock.Anything, testifyMock.Anything, "listing1").
					Return(domain.ClaimRequest{}, errors.New("database error"))
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Renderer = &AdminMockRenderer{}

			req := httptest.NewRequest(http.MethodPost, "/listings/listing1/claim", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.listingID)

			if tt.user != nil {
				c.Set("User", *tt.user)
			}

			mockRepo := &mock.MockListingRepository{}
			mockSvc := &mock.MockListingService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := &handler.ListingHandler{
				Repo:       mockRepo,
				ListingSvc: mockSvc,
			}
			_ = h.HandleClaim(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

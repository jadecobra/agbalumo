package handler_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
)

func TestSubmitFeedback(t *testing.T) {
	tests := []struct {
		name           string
		user           *domain.User
		body           string
		mockSetup      func() *mock.MockListingRepository
		expectedStatus int
	}{
		{
			name: "Success",
			user: &domain.User{ID: "user-1"},
			body: "type=Issue&content=Bug+Report",
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					SaveFeedbackFn: func(ctx context.Context, f domain.Feedback) error {
						if f.Type != domain.FeedbackTypeIssue {
							return errors.New("wrong type")
						}
						if f.Content != "Bug Report" {
							return errors.New("wrong content")
						}
						// Check UserID
						if f.UserID != "user-1" {
							return errors.New("wrong user id")
						}
						return nil
					},
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Unauthenticated",
			user: nil,
			body: "type=Issue&content=Bug",
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{}
			},
			expectedStatus: http.StatusUnauthorized, // Or redirect, checking logic
		},
		{
			name: "MissingContent",
			user: &domain.User{ID: "user-1"},
			body: "type=Issue&content=",
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					SaveFeedbackFn: func(ctx context.Context, f domain.Feedback) error {
						return errors.New("should not be called")
					},
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "RepoError",
			user: &domain.User{ID: "user-1"},
			body: "type=Other&content=Stuff",
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					SaveFeedbackFn: func(ctx context.Context, f domain.Feedback) error {
						return errors.New("db failed")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodPost, "/feedback", strings.NewReader(tt.body))
			
			if tt.user != nil {
				c.Set("User", *tt.user)
			}

			h := handler.NewFeedbackHandler(tt.mockSetup())
			
			// We assume HandleSubmit will be the method name
			err := h.HandleSubmit(c)
			
			// Handle errors that Echo might return (e.g. 400/500)
			if err != nil {
				// In Echo, returning an error often means it's processed by error handler.
				// However, if we return c.String(status, msg), err is nil.
				// If we return echo.NewHTTPError, err is not nil.
				// We'll check rec.Code mostly.
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

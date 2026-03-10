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

func TestAdminHandler_HandleLoginView(t *testing.T) {
	tests := []struct {
		name       string
		user       interface{}
		expectCode int
		expectLoc  string
	}{
		{
			name:       "AlreadyAdmin_RedirectsToDashboard",
			user:       domain.User{ID: "u1", Role: domain.UserRoleAdmin},
			expectCode: http.StatusTemporaryRedirect,
			expectLoc:  "/admin",
		},
		{
			name:       "NonAdmin_RendersLoginForm",
			user:       domain.User{ID: "u2", Role: "user"},
			expectCode: http.StatusOK,
		},
		{
			name:       "NoUser_RendersLoginForm",
			user:       nil,
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupAdminTestContext(http.MethodGet, "/admin/login", nil)
			if tt.user != nil {
				c.Set("User", tt.user)
			}

			mockRepo := &mock.MockListingRepository{}
			h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
			_ = h.HandleLoginView(c)

			assert.Equal(t, tt.expectCode, rec.Code)
			if tt.expectLoc != "" {
				assert.Equal(t, tt.expectLoc, rec.Header().Get("Location"))
			}
		})
	}
}

func TestAdminHandler_HandleLoginAction(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		adminCode  string
		user       interface{}
		setupMock  func(*mock.MockListingRepository)
		expectCode int
		expectLoc  string
	}{
		{
			name:       "WrongCode_RendersError",
			code:       "wrong",
			adminCode:  "correct",
			user:       domain.User{ID: "u1", Role: "user"},
			expectCode: http.StatusOK,
		},
		{
			name:       "NoUser_RedirectsToLogin",
			code:       "correct",
			adminCode:  "correct",
			user:       nil,
			expectCode: http.StatusTemporaryRedirect,
			expectLoc:  "/auth/google/login",
		},
		{
			name:      "ValidCode_PromotesUser",
			code:      "secret",
			adminCode: "secret",
			user:      domain.User{ID: "u1", Role: "user"},
			setupMock: func(r *mock.MockListingRepository) {
				r.On("SaveUser", testifyMock.Anything, testifyMock.Anything).Return(nil)
			},
			expectCode: http.StatusFound,
			expectLoc:  "/admin",
		},
		{
			name:      "ValidCode_SaveUserError",
			code:      "secret",
			adminCode: "secret",
			user:      domain.User{ID: "u1", Role: "user"},
			setupMock: func(r *mock.MockListingRepository) {
				r.On("SaveUser", testifyMock.Anything, testifyMock.Anything).Return(assert.AnError)
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("code", tt.code)
			c, rec := setupAdminTestContext(http.MethodPost, "/admin/login", strings.NewReader(formData.Encode()))
			if tt.user != nil {
				c.Set("User", tt.user)
			}

			cfg := config.LoadConfig()
			cfg.AdminCode = tt.adminCode

			mockRepo := &mock.MockListingRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			h := handler.NewAdminHandler(mockRepo, nil, cfg)
			_ = h.HandleLoginAction(c)

			assert.Equal(t, tt.expectCode, rec.Code)
			if tt.expectLoc != "" {
				assert.Equal(t, tt.expectLoc, rec.Header().Get("Location"))
			}
		})
	}
}

package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminMiddleware(t *testing.T) {
	e := echo.New()
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo)

	t.Run("Redirects when no user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.AdminMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
			assert.Equal(t, "/auth/google/login", rec.Header().Get("Location"))
		}
	})

	t.Run("Redirects to login when user is not admin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("User", domain.User{Role: domain.UserRoleUser})

		err := h.AdminMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
			assert.Equal(t, "/admin/login", rec.Header().Get("Location"))
		}
	})

	t.Run("Allows access when user is admin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("User", domain.User{Role: domain.UserRoleAdmin})

		err := h.AdminMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "Secret")
		})(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "Secret", rec.Body.String())
		}
	})
}

func TestHandleLoginView(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewAdminHandler(nil)

	if assert.NoError(t, h.HandleLoginView(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandleLoginAction(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	mockRepo := &mock.MockListingRepository{
		SaveUserFn: func(ctx context.Context, u domain.User) error {
			assert.Equal(t, domain.UserRoleAdmin, u.Role)
			return nil
		},
	}
	h := handler.NewAdminHandler(mockRepo)

	t.Run("Success", func(t *testing.T) {
		form := make(url.Values)
		form.Set("code", "agbalumo2024")
		req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("User", domain.User{ID: "u1", Role: domain.UserRoleUser})

		if assert.NoError(t, h.HandleLoginAction(c)) {
			assert.Equal(t, http.StatusFound, rec.Code)
			assert.Equal(t, "/admin", rec.Header().Get("Location"))
		}
	})

	t.Run("Invalid Code", func(t *testing.T) {
		form := make(url.Values)
		form.Set("code", "wrong")
		req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.HandleLoginAction(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// Should render template with error
		}
	})
}

func TestHandleDashboard(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	mockRepo := &mock.MockListingRepository{
		GetPendingListingsFn: func(ctx context.Context) ([]domain.Listing, error) {
			return []domain.Listing{{ID: "1", Title: "Pending"}}, nil
		},
		GetUserCountFn: func(ctx context.Context) (int, error) {
			return 100, nil
		},
		GetFeedbackCountsFn: func(ctx context.Context) (map[domain.FeedbackType]int, error) {
			return map[domain.FeedbackType]int{domain.FeedbackTypeIssue: 5}, nil
		},
	}
	h := handler.NewAdminHandler(mockRepo)

	if assert.NoError(t, h.HandleDashboard(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandleApprove(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/1/approve", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/listings/:id/approve")
	c.SetParamNames("id")
	c.SetParamValues("1")

	saved := false
	mockRepo := &mock.MockListingRepository{
		FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
			return domain.Listing{ID: "1", Status: domain.ListingStatusPending, IsActive: true}, nil
		},
		SaveFn: func(ctx context.Context, l domain.Listing) error {
			saved = true
			assert.Equal(t, domain.ListingStatusApproved, l.Status)
			assert.True(t, l.IsActive)
			return nil
		},
	}
	h := handler.NewAdminHandler(mockRepo)

	if assert.NoError(t, h.HandleApprove(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, saved)
	}
}

func TestHandleReject(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/1/reject", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/listings/:id/reject")
	c.SetParamNames("id")
	c.SetParamValues("1")

	saved := false
	mockRepo := &mock.MockListingRepository{
		FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
			return domain.Listing{ID: "1", Status: domain.ListingStatusPending, IsActive: true}, nil
		},
		SaveFn: func(ctx context.Context, l domain.Listing) error {
			saved = true
			assert.Equal(t, domain.ListingStatusRejected, l.Status)
			assert.False(t, l.IsActive)
			return nil
		},
	}
	h := handler.NewAdminHandler(mockRepo)

	if assert.NoError(t, h.HandleReject(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, saved)
	}
}

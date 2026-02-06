package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
)

func TestHandleLoginView(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewAdminHandler(nil)

	if err := h.HandleLoginView(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rec.Code)
	}
}

func TestHandleLoginAction_Success(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	form := strings.NewReader("code=agbalumo2024")
	req := httptest.NewRequest(http.MethodPost, "/admin/login", form)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewAdminHandler(nil)

	if err := h.HandleLoginAction(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusFound {
		t.Errorf("Expected 302 Found, got %d", rec.Code)
	}

	// Check Cookie
	cookie := rec.Result().Cookies()[0]
	if cookie.Name != "admin_session" || cookie.Value != "authenticated" {
		t.Errorf("Invalid cookie: %v", cookie)
	}
}

func TestHandleLoginAction_Failure(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	form := strings.NewReader("code=wrongpass")
	req := httptest.NewRequest(http.MethodPost, "/admin/login", form)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewAdminHandler(nil)

	if err := h.HandleLoginAction(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK (re-render), got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Invalid Access Code") {
		t.Errorf("Expected error message, got %s", rec.Body.String())
	}
}

func TestHandleDashboard(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{
		FindAllFn: func(ctx context.Context, filterType, query string, includeInactive bool) ([]domain.Listing, error) {
			if !includeInactive {
				t.Error("Admin should see inactive listings")
			}
			return []domain.Listing{{Title: "Item 1"}}, nil
		},
	}
	h := handler.NewAdminHandler(mockRepo)

	if err := h.HandleDashboard(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rec.Code)
	}
}

func TestHandleDelete_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/admin/listings/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	saved := false
	mockRepo := &mock.MockListingRepository{
		FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
			return domain.Listing{ID: "1", IsActive: true}, nil
		},
		SaveFn: func(ctx context.Context, l domain.Listing) error {
			saved = true
			if l.IsActive {
				t.Error("Expected IsActive to be false")
			}
			return nil
		},
	}

	h := handler.NewAdminHandler(mockRepo)

	// Execute
	if err := h.HandleDelete(c); err != nil {
		t.Fatalf("HandleDelete failed: %v", err)
	}

	// Verify
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if !saved {
		t.Error("Expected listing to be marked inactive (soft deleted)")
	}
}

func TestAuthMiddleware_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo)

	// Dummy handler to wrap
	next := func(c echo.Context) error {
		return c.String(http.StatusOK, "Secret")
	}

	// Execute Middleware
	handlerFunc := h.AuthMiddleware(next)
	handlerFunc(c)

	// Verify Redirect
	if rec.Code != http.StatusFound {
		t.Errorf("Expected status 302 Found (Redirect), got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/login" {
		t.Errorf("Expected redirect to /admin/login, got %s", loc)
	}
}

func TestAuthMiddleware_Authorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "admin_session", Value: "authenticated"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo)

	// Dummy handler to wrap
	next := func(c echo.Context) error {
		return c.String(http.StatusOK, "Secret")
	}

	// Execute Middleware
	handlerFunc := h.AuthMiddleware(next)
	handlerFunc(c)

	// Verify Access
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "Secret" {
		t.Errorf("Expected body 'Secret', got %s", rec.Body.String())
	}
}

func TestHandleDashboard_RepoError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{
		FindAllFn: func(ctx context.Context, filterType, query string, includeInactive bool) ([]domain.Listing, error) {
			return nil, errors.New("db disconnect")
		},
	}
	h := handler.NewAdminHandler(mockRepo)

	if err := h.HandleDashboard(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rec.Code)
	}
}

func TestHandleDelete_RepoNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/admin/listings/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockRepo := &mock.MockListingRepository{
		FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
			return domain.Listing{}, errors.New("not found")
		},
	}
	h := handler.NewAdminHandler(mockRepo)

	if err := h.HandleDelete(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rec.Code)
	}
}

func TestHandleDelete_RepoSaveError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodDelete, "/admin/listings/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockRepo := &mock.MockListingRepository{
		FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
			return domain.Listing{ID: "1"}, nil
		},
		SaveFn: func(ctx context.Context, l domain.Listing) error {
			return errors.New("save failed")
		},
	}
	h := handler.NewAdminHandler(mockRepo)

	if err := h.HandleDelete(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rec.Code)
	}
}

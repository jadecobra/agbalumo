package handler_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAdminHandler_HandleAdminDeleteAction_Success(t *testing.T) {
	cfg := config.LoadConfig()
	cfg.AdminCode = "secret"

	formData := url.Values{}
	formData.Set("admin_code", "secret")
	formData.Add("id", "l1")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("Delete", testifyMock.Anything, "l1").Return(nil)

	h := handler.NewAdminHandler(mockRepo, nil, cfg)
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestHandleAdminDeleteView(t *testing.T) {
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})
	mockRepo.On("FindByID", testifyMock.Anything, "listing1").Return(domain.Listing{ID: "listing1", Title: "To Delete"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/listings/delete?id=listing1", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplateForPage(t, "admin_delete_confirm.html")}
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleAdminDeleteView(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandleAdminDeleteView_NoIDs_Redirects(t *testing.T) {
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/delete", nil)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteView_FindByIDError_Returns404(t *testing.T) {
	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "bad-id").Return(domain.Listing{}, assert.AnError)
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

	req := httptest.NewRequest(http.MethodGet, "/admin/listings/delete?id=bad-id", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	e.Renderer = &AdminMockRenderer{}
	c := e.NewContext(req, rec)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleAdminDeleteAction_NoIDs_Redirects(t *testing.T) {
	formData := url.Values{}
	formData.Set("admin_code", "secret")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))

	cfg := config.LoadConfig()
	cfg.AdminCode = "secret"

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, cfg)

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteAction_WrongCode_RendersConfirmWithError(t *testing.T) {
	formData := url.Values{}
	formData.Set("admin_code", "wrong")
	formData.Add("id", "l1")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))

	cfg := config.LoadConfig()
	cfg.AdminCode = "correct"

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, cfg)

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleAdminDeleteAction_DeleteError_Logs(t *testing.T) {
	cfg := config.LoadConfig()
	cfg.AdminCode = "secret"

	formData := url.Values{}
	formData.Set("admin_code", "secret")
	formData.Add("id", "l1")
	formData.Add("id", "l2")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("Delete", testifyMock.Anything, "l1").Return(nil)
	mockRepo.On("Delete", testifyMock.Anything, "l2").Return(assert.AnError)

	h := handler.NewAdminHandler(mockRepo, nil, cfg)
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)
}

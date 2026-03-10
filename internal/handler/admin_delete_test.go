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

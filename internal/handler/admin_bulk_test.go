package handler_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleBulkUpload(t *testing.T) {
	e := echo.New()
	// CSV headers: title,type,description,origin,email,phone,whatsapp
	csvContent := "title,type,description,origin,email,phone,whatsapp\nTest Biz,Business,Description,Nigeria,test@test.com,,"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("csv_file", "upload.csv")
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, service.NewCSVService(), config.LoadConfig())
	err := h.HandleBulkUpload(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify listing was saved
	listings, _ := repo.FindByTitle(context.Background(), "Test Biz")
	assert.Len(t, listings, 1)
	assert.Equal(t, "Test Biz", listings[0].Title)
}

func TestAdminHandler_HandleBulkUpload_NoFile(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/upload", nil)
	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	h := handler.NewAdminHandler(nil, service.NewCSVService(), config.LoadConfig())
	_ = h.HandleBulkUpload(c)
	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestAdminHandler_HandleBulkUpload_ParseError(t *testing.T) {
	e := echo.New()
	// Missing required "description"
	csvContent := "title,type,origin,email\nTest Biz,Business,Nigeria,test@test.com"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("csv_file", "upload.csv")
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, service.NewCSVService(), config.LoadConfig())
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	// Verify no listing was saved
	listings, _ := repo.FindAll(context.Background(), "", "", "", "", true, 10, 0)
	assert.Empty(t, listings)
}

func TestHandleBulkAction_NoSelection(t *testing.T) {
	e := echo.New()
	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, nil, &config.Config{})

	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
	}
}

func TestHandleBulkAction_Approve(t *testing.T) {
	e := echo.New()
	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, nil, &config.Config{})

	_ = repo.Save(context.Background(), domain.Listing{ID: "l1", Title: "L1", Status: domain.ListingStatusPending})

	form := url.Values{}
	form.Add("action", "approve")
	form.Add("selectedListings", "l1")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		l, _ := repo.FindByID(context.Background(), "l1")
		assert.Equal(t, domain.ListingStatusApproved, l.Status)
	}
}

func TestHandleBulkAction_Reject(t *testing.T) {
	e := echo.New()
	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, nil, &config.Config{})

	_ = repo.Save(context.Background(), domain.Listing{ID: "l1", Title: "L1", Status: domain.ListingStatusPending})

	form := url.Values{}
	form.Add("action", "reject")
	form.Add("selectedListings", "l1")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		l, _ := repo.FindByID(context.Background(), "l1")
		assert.Equal(t, domain.ListingStatusRejected, l.Status)
	}
}

func TestHandleBulkAction_Delete(t *testing.T) {
	e := echo.New()
	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, nil, &config.Config{})

	form := url.Values{}
	form.Add("action", "delete")
	form.Add("selectedListings", "l1")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Contains(t, rec.Header().Get("Location"), "/admin/listings/delete-confirm?id=l1")
	}
}


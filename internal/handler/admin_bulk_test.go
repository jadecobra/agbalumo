package handler_test

import (
	"bytes"
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
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAdminHandler_HandleBulkUpload(t *testing.T) {
	e := echo.New()
	csvContent := "title,type,description,origin,email\nTest Biz,Business,Description,Nigeria,test@test.com"

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

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("Save", testifyMock.Anything, testifyMock.Anything).Return(nil)

	h := handler.NewAdminHandler(mockRepo, service.NewCSVService(), config.LoadConfig())
	err := h.HandleBulkUpload(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))
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
	csvContent := "invalid,csv,content\nmissing,columns"

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

	mockRepo := &mock.MockListingRepository{}

	// Create a mock for CSV service that returns an error
	// To keep it simple, let's use the real service but pass a repo that fails to save
	// This exercises the parse/import error logic
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, assert.AnError)
	mockRepo.On("Save", testifyMock.Anything, testifyMock.Anything).Return(assert.AnError)

	h := handler.NewAdminHandler(mockRepo, service.NewCSVService(), config.LoadConfig())
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestHandleBulkAction_NoSelection(t *testing.T) {
	e := echo.New()
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

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
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

	mockRepo.On("FindByID", testifyMock.Anything, "l1").Return(domain.Listing{ID: "l1", Status: domain.ListingStatusPending}, nil)
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.ID == "l1" && l.Status == domain.ListingStatusApproved
	})).Return(nil)

	form := url.Values{}
	form.Add("action", "approve")
	form.Add("selectedListings", "l1")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
	}
}

func TestHandleBulkAction_Reject(t *testing.T) {
	e := echo.New()
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

	mockRepo.On("FindByID", testifyMock.Anything, "l1").Return(domain.Listing{ID: "l1", Status: domain.ListingStatusPending}, nil)
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.ID == "l1" && l.Status == domain.ListingStatusRejected
	})).Return(nil)

	form := url.Values{}
	form.Add("action", "reject")
	form.Add("selectedListings", "l1")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
	}
}

func TestHandleBulkAction_Delete(t *testing.T) {
	e := echo.New()
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

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

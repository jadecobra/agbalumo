package handler

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAdminHandler_HandleDashboard(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock User
	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{{ID: "cr1", ListingTitle: "Test Listing"}}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(10, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{domain.FeedbackTypeIssue: 2}, nil)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{
		{ID: "f1", UserID: "u1", Type: domain.FeedbackTypeIssue, Content: "Test Feedback"},
	}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil)
	mockRepo.On("GetAllUsers", testifyMock.Anything, 10, 0).Return([]domain.User{}, nil)

	h := NewAdminHandler(mockRepo, service.NewCSVService(), config.LoadConfig())

	// Set Renderer (mock)
	e.Renderer = &mock.MockRenderer{}

	err := h.HandleDashboard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleUsers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("GetAllUsers", testifyMock.Anything, 50, 0).Return([]domain.User{{ID: "u1", Name: "User 1"}}, nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
	e.Renderer = &mock.MockRenderer{}

	err := h.HandleUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleApproveClaim(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/claims/cr1/approve", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/claims/:id/approve")
	c.SetParamNames("id")
	c.SetParamValues("cr1")

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("UpdateClaimRequestStatus", testifyMock.Anything, "cr1", domain.ClaimStatusApproved).Return(nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleApproveClaim(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleRejectClaim(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/claims/cr1/reject", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/claims/:id/reject")
	c.SetParamNames("id")
	c.SetParamValues("cr1")

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("UpdateClaimRequestStatus", testifyMock.Anything, "cr1", domain.ClaimStatusRejected).Return(nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleRejectClaim(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleBulkAction_Approve(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("action", "approve")
	formData.Add("selectedListings", "1")
	formData.Add("selectedListings", "2")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", Status: domain.ListingStatusPending}, nil)
	mockRepo.On("FindByID", testifyMock.Anything, "2").Return(domain.Listing{ID: "2", Status: domain.ListingStatusPending}, nil)

	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.Status == domain.ListingStatusApproved && l.IsActive && (l.ID == "1" || l.ID == "2")
	})).Return(nil).Twice()

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleBulkAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleBulkAction_Reject(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("action", "reject")
	formData.Add("selectedListings", "1")
	formData.Add("selectedListings", "2")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", Status: domain.ListingStatusPending}, nil)
	mockRepo.On("FindByID", testifyMock.Anything, "2").Return(domain.Listing{ID: "2", Status: domain.ListingStatusPending}, nil)

	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.Status == domain.ListingStatusRejected && !l.IsActive && (l.ID == "1" || l.ID == "2")
	})).Return(nil).Twice()

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleBulkAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleBulkAction_DeleteRedirect(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("action", "delete")
	formData.Add("selectedListings", "1")
	formData.Add("selectedListings", "2")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	h := NewAdminHandler(nil, nil, config.LoadConfig())

	err := h.HandleBulkAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings/delete-confirm?id=1&id=2", rec.Header().Get("Location"))
}

func TestAdminHandler_HandleBulkAction_NoSelection(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("action", "approve")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	h := NewAdminHandler(nil, nil, config.LoadConfig())

	err := h.HandleBulkAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestAdminHandler_HandleBulkAction_FindByIDError(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("action", "approve")
	formData.Add("selectedListings", "1")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{}, errors.New("not found"))

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleBulkAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleBulkAction_UnknownAction(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("action", "unknown")
	formData.Add("selectedListings", "1")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1"}, nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleBulkAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Save", testifyMock.Anything, testifyMock.Anything)
}

func TestAdminHandler_HandleBulkAction_SaveError(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("action", "approve")
	formData.Add("selectedListings", "1")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", Status: domain.ListingStatusPending}, nil)
	mockRepo.On("Save", testifyMock.Anything, testifyMock.Anything).Return(errors.New("save failed"))

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleBulkAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleLoginAction_Success(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("code", "agbalumo2024")
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := domain.User{ID: "user1", Role: domain.UserRoleUser}
	c.Set("User", user)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.MatchedBy(func(u domain.User) bool {
		return u.Role == domain.UserRoleAdmin
	})).Return(nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleLoginAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleLoginAction_InvalidCode(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("code", "wrongcode")
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewAdminHandler(nil, nil, config.LoadConfig())
	e.Renderer = &mock.MockRenderer{}

	err := h.HandleLoginAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_HandleAdminDeleteView(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/listings/delete-confirm?id=1&id=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1"}, nil)
	mockRepo.On("FindByID", testifyMock.Anything, "2").Return(domain.Listing{ID: "2"}, nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
	e.Renderer = &mock.MockRenderer{}

	err := h.HandleAdminDeleteView(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_HandleAdminDeleteAction_Success(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("admin_code", "agbalumo2024")
	formData.Add("id", "1")
	formData.Add("id", "2")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("Delete", testifyMock.Anything, "1").Return(nil)
	mockRepo.On("Delete", testifyMock.Anything, "2").Return(nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleAdminDeleteAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleAdminDeleteAction_WrongCode(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("admin_code", "wrong")
	formData.Add("id", "1")
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	h := NewAdminHandler(nil, nil, config.LoadConfig())
	e.Renderer = &mock.MockRenderer{}

	err := h.HandleAdminDeleteAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code) // Re-renders form
}

func TestAdminHandler_HandleBulkUpload_Success(t *testing.T) {
	e := echo.New()

	// Create CSV content
	csvContent := "title,type,description,origin,email\nTest Biz,Business,Description,Nigeria,test@test.com"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("csv_file", "upload.csv")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(csvContent))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	// Expect Save to be called for the valid row
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.Title == "Test Biz" && l.Type == domain.Business && l.OwnerOrigin == "Nigeria"
	})).Return(nil)

	h := NewAdminHandler(mockRepo, service.NewCSVService(), config.LoadConfig())

	err = h.HandleBulkUpload(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleBulkUpload_NoFile(t *testing.T) {
	e := echo.New()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	h := NewAdminHandler(nil, service.NewCSVService(), config.LoadConfig())
	e.Renderer = &mock.MockRenderer{}

	err := h.HandleBulkUpload(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))
}

func TestAdminHandler_HandleBulkUpload_ParseError(t *testing.T) {
	e := echo.New()

	csvContent := "title,description\nTest,Desc"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("csv_file", "upload.csv")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(csvContent))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	h := NewAdminHandler(nil, service.NewCSVService(), config.LoadConfig())

	err = h.HandleBulkUpload(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))
}

func TestAdminHandler_HandleBulkUpload_WithErrors(t *testing.T) {
	e := echo.New()

	csvContent := "title,type,description,origin,email\nGood,Business,Desc,NG,a@b.com\nBad,Business,,NG,b@b.com\nBad2,Business,,NG,c@b.com\nBad3,Business,,NG,d@b.com"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("csv_file", "upload.csv")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(csvContent))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	store := customMiddleware.NewTestSessionStore()
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("Save", testifyMock.Anything, testifyMock.Anything).Return(nil).Maybe()

	h := NewAdminHandler(mockRepo, service.NewCSVService(), config.LoadConfig())

	err = h.HandleBulkUpload(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
}

// --- AdminMiddleware Tests ---

func TestAdminMiddleware_NoUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewAdminHandler(nil, nil, config.LoadConfig())

	called := false
	handler := h.AdminMiddleware(func(c echo.Context) error {
		called = true
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.False(t, called, "Next handler should not be called when no user")
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/auth/google/login", rec.Header().Get("Location"))
}

func TestAdminMiddleware_NonAdminUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	regularUser := domain.User{ID: "u1", Role: domain.UserRoleUser}
	c.Set("User", regularUser)

	h := NewAdminHandler(nil, nil, config.LoadConfig())

	called := false
	handler := h.AdminMiddleware(func(c echo.Context) error {
		called = true
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.False(t, called, "Next handler should not be called for non-admin")
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/admin/login", rec.Header().Get("Location"))
}

func TestAdminMiddleware_AdminUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	h := NewAdminHandler(nil, nil, config.LoadConfig())

	called := false
	handler := h.AdminMiddleware(func(c echo.Context) error {
		called = true
		return c.String(http.StatusOK, "dashboard")
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called, "Next handler should be called for admin")
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- HandleLoginView Tests ---

func TestAdminHandler_HandleLoginView_AlreadyAdmin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	h := NewAdminHandler(nil, nil, config.LoadConfig())

	err := h.HandleLoginView(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))
}

func TestAdminHandler_HandleLoginView_NotAdmin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Renderer = &mock.MockRenderer{}

	// No user set — should render the login form
	h := NewAdminHandler(nil, nil, config.LoadConfig())

	err := h.HandleLoginView(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- HandleDashboard Error Path Tests ---

func TestAdminHandler_HandleDashboard_PendingClaimRequestsError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Renderer = &mock.MockRenderer{}

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, assert.AnError)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, nil).Maybe()
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil).Maybe()

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleDashboard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAdminHandler_HandleDashboard_UserCountError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Renderer = &mock.MockRenderer{}

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(0, assert.AnError)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, nil).Maybe()
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil).Maybe()

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleDashboard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAdminHandler_HandleDashboard_FeedbackCountsError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Renderer = &mock.MockRenderer{}

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, assert.AnError)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, nil).Maybe()
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil).Maybe()

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleDashboard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAdminHandler_HandleDashboard_ListingGrowthError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Renderer = &mock.MockRenderer{}

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, nil).Maybe()
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, assert.AnError)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil).Maybe()

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleDashboard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAdminHandler_HandleDashboard_UserGrowthError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Renderer = &mock.MockRenderer{}

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, nil).Maybe()
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, assert.AnError)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil).Maybe()

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleDashboard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAdminHandler_HandleDashboard_GetAllFeedbackError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Renderer = &mock.MockRenderer{}

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, assert.AnError)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleDashboard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- HandleLoginAction Error Path Tests ---

func TestAdminHandler_HandleLoginAction_NoUser(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("code", "agbalumo2024")
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No user set

	h := NewAdminHandler(nil, nil, config.LoadConfig())

	err := h.HandleLoginAction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/auth/google/login", rec.Header().Get("Location"))
}

type AdminMockRenderer struct{}

func (m *AdminMockRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return nil
}

func TestAdminHandler_HandleAllListings(t *testing.T) {
	e := echo.New()
	e.Renderer = &AdminMockRenderer{}

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/listings?page=1&category=Business", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo := &mock.MockListingRepository{}
		mockRepo.On("FindAll", testifyMock.Anything, "Business", "", "", "", true, 50, 0).Return([]domain.Listing{{ID: "1"}}, nil)
		mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{domain.Business: 5, domain.Food: 3}, nil)
		mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

		h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
		c.Set("User", domain.User{ID: "admin-1"})

		err := h.HandleAllListings(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DBError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/listings", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo := &mock.MockListingRepository{}
		mockRepo.On("FindAll", testifyMock.Anything, "", "", "", "", true, 50, 0).Return([]domain.Listing{}, assert.AnError)
		mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil).Maybe()
		mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

		h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
		c.Set("User", domain.User{ID: "admin-1"})

		err := h.HandleAllListings(c)
		assert.NoError(t, err) // RespondError handles it and returns nil
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestAdminHandler_HandleToggleFeatured(t *testing.T) {
	e := echo.New()
	e.Renderer = &mock.MockRenderer{}

	t.Run("Feature", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/listings/123/featured", strings.NewReader("featured=true"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("123")

		mockRepo := &mock.MockListingRepository{}
		mockRepo.On("SetFeatured", testifyMock.Anything, "123", true).Return(nil)

		h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
		c.Set("User", domain.User{ID: "admin-1", Role: domain.UserRoleAdmin})

		err := h.HandleToggleFeatured(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unfeature", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/listings/456/featured", strings.NewReader("featured=false"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("456")

		mockRepo := &mock.MockListingRepository{}
		mockRepo.On("SetFeatured", testifyMock.Anything, "456", false).Return(nil)

		h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
		c.Set("User", domain.User{ID: "admin-1", Role: domain.UserRoleAdmin})

		err := h.HandleToggleFeatured(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DBError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/listings/123/featured", strings.NewReader("featured=true"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("123")

		mockRepo := &mock.MockListingRepository{}
		mockRepo.On("SetFeatured", testifyMock.Anything, "123", true).Return(errors.New("db error"))

		h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
		c.Set("User", domain.User{ID: "admin-1", Role: domain.UserRoleAdmin})

		err := h.HandleToggleFeatured(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidBoolean", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/listings/123/featured", strings.NewReader("featured=notabool"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("123")

		mockRepo := &mock.MockListingRepository{}
		// Current implementation parses "notabool" == "true" as false and proceeds to DB
		mockRepo.On("SetFeatured", testifyMock.Anything, "123", false).Return(nil)

		h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
		c.Set("User", domain.User{ID: "admin-1", Role: domain.UserRoleAdmin})
		e.Renderer = &mock.MockRenderer{}

		err := h.HandleToggleFeatured(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MissingID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/listings//featured", strings.NewReader("featured=true"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		mockRepo := &mock.MockListingRepository{}
		h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
		c.Set("User", domain.User{ID: "admin-1", Role: domain.UserRoleAdmin})
		e.Renderer = &mock.MockRenderer{}

		err := h.HandleToggleFeatured(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAdminHandler_HandleApproveClaim_ErrorPaths(t *testing.T) {
	e := echo.New()
	e.Renderer = &mock.MockRenderer{}

	t.Run("NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/claims/cr999/approve", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/admin/claims/:id/approve")
		c.SetParamNames("id")
		c.SetParamValues("cr999")
		c.Set("User", domain.User{Role: domain.UserRoleAdmin})

		mockRepo := &mock.MockListingRepository{}
		mockRepo.On("UpdateClaimRequestStatus", testifyMock.Anything, "cr999", domain.ClaimStatusApproved).Return(errors.New("not found"))

		h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
		err := h.HandleApproveClaim(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAdminHandler_HandleRejectClaim_RepoError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/claims/cr1/reject", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/claims/:id/reject")
	c.SetParamNames("id")
	c.SetParamValues("cr1")

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("UpdateClaimRequestStatus", testifyMock.Anything, "cr1", domain.ClaimStatusRejected).Return(errors.New("repo error"))

	h := NewAdminHandler(mockRepo, nil, &config.Config{})
	_ = h.HandleRejectClaim(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAdminHandler_HandleAdminDeleteView_RepoError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/listings/delete-confirm?id=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{}, errors.New("repo error"))

	h := NewAdminHandler(mockRepo, nil, &config.Config{})
	_ = h.HandleAdminDeleteView(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAdminHandler_HandleUsers_RepoError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("GetAllUsers", testifyMock.Anything, 50, 0).Return([]domain.User{}, errors.New("repo error"))

	h := NewAdminHandler(mockRepo, nil, &config.Config{})
	e.Renderer = &mock.MockRenderer{}
	_ = h.HandleUsers(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

package handler

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
	mockRepo.On("GetPendingListings", testifyMock.Anything).Return([]domain.Listing{{ID: "1", Title: "Pending Listing"}}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(10, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{domain.FeedbackTypeIssue: 2}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)

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
	mockRepo.On("GetAllUsers", testifyMock.Anything).Return([]domain.User{{ID: "u1", Name: "User 1"}}, nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())
	e.Renderer = &mock.MockRenderer{}

	err := h.HandleUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleApprove(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/1/approve", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/listings/:id/approve")
	c.SetParamNames("id")
	c.SetParamValues("1")

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", Status: domain.ListingStatusPending}, nil)
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.Status == domain.ListingStatusApproved && l.IsActive
	})).Return(nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleApprove(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAdminHandler_HandleReject(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/listings/1/reject", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/listings/:id/reject")
	c.SetParamNames("id")
	c.SetParamValues("1")

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", Status: domain.ListingStatusPending}, nil)
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.Status == domain.ListingStatusRejected && !l.IsActive
	})).Return(nil)

	h := NewAdminHandler(mockRepo, nil, config.LoadConfig())

	err := h.HandleReject(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
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

	h := NewAdminHandler(nil, service.NewCSVService(), config.LoadConfig())
	e.Renderer = &mock.MockRenderer{}

	err := h.HandleBulkUpload(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
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

func TestAdminHandler_HandleDashboard_PendingListingsError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	e.Renderer = &mock.MockRenderer{}

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("GetPendingListings", testifyMock.Anything).Return([]domain.Listing{}, assert.AnError)

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
	mockRepo.On("GetPendingListings", testifyMock.Anything).Return([]domain.Listing{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(0, assert.AnError)

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
	mockRepo.On("GetPendingListings", testifyMock.Anything).Return([]domain.Listing{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, assert.AnError)

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
	mockRepo.On("GetPendingListings", testifyMock.Anything).Return([]domain.Listing{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, assert.AnError)

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
	mockRepo.On("GetPendingListings", testifyMock.Anything).Return([]domain.Listing{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, assert.AnError)

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

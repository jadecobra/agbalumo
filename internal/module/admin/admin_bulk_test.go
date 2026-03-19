package admin_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"

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
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: service.NewCSVService(), Cfg: config.LoadConfig()})

	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	err := h.HandleBulkUpload(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify listing was saved
	listings, _ := repo.FindByTitle(context.Background(), "Test Biz")
	assert.Len(t, listings, 1)
	assert.Equal(t, "Test Biz", listings[0].Title)
}

func TestAdminHandler_HandleBulkUpload_InvalidCSV(t *testing.T) {
	e := echo.New()
	// Junk content
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("csv_file", "junk.csv")
	_, _ = part.Write([]byte("invalid,csv,data\n1,2,3"))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	repo := handler.SetupTestRepository(t)
	// We need to inject a mock session store because HandleBulkUpload uses it for flash messages
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: service.NewCSVService(), Cfg: config.LoadConfig()})
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestAdminHandler_HandleBulkUpload_NoFile(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/upload", nil)
	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)

	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: nil, FeedbackStore: nil, AnalyticsStore: nil, CategoryStore: nil, UserStore: nil, ListingStore: nil, ClaimRequestStore: nil, CSVService: service.NewCSVService(), Cfg: config.LoadConfig()})
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
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: service.NewCSVService(), Cfg: config.LoadConfig()})
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	// Verify no listing was saved
	listings, _, _ := repo.FindAll(context.Background(), "", "", "", "", true, 10, 0)
	assert.Empty(t, listings)
}

func TestHandleBulkAction_NoSelection(t *testing.T) {
	e := echo.New()
	repo := handler.SetupTestRepository(t)
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: &config.Config{}})

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
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: &config.Config{}})

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
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: &config.Config{}})

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
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: &config.Config{}})

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

func TestHandleBulkAction_ChangeCategory(t *testing.T) {
	e := echo.New()
	repo := handler.SetupTestRepository(t)
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: &config.Config{}})

	// Create listings with "Business" category
	_ = repo.Save(context.Background(), domain.Listing{ID: "l1", Title: "L1", Type: domain.Business, City: "Lagos", OwnerOrigin: "Nigeria"})
	_ = repo.Save(context.Background(), domain.Listing{ID: "l2", Title: "L2", Type: domain.Business, City: "Accra", OwnerOrigin: "Ghana"})

	form := url.Values{}
	form.Add("action", "change_category")
	form.Add("selectedListings", "l1")
	form.Add("selectedListings", "l2")
	form.Add("new_category", string(domain.Job))

	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		l1, _ := repo.FindByID(context.Background(), "l1")
		l2, _ := repo.FindByID(context.Background(), "l2")
		assert.Equal(t, domain.Job, l1.Type)
		assert.Equal(t, domain.Job, l2.Type)
	}
}

func TestHandleBulkAction_NoAction(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/bulk", strings.NewReader("ids[]=1"))
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	repo := handler.SetupTestRepository(t)
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
	_ = h.HandleBulkAction(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestAdminHandler_HandleBulkUpload_ManyErrors(t *testing.T) {
	e := echo.New()
	// CSV with 4 invalid rows (missing title/desc)
	csvContent := "title,type,description\n,,\n,,\n,,\n,,\n"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("csv_file", "junk.csv")
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	repo := handler.SetupTestRepository(t)
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: service.NewCSVService(), Cfg: config.LoadConfig()})
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestHandleBulkAction_ChangeToCustomCategory(t *testing.T) {
	e := echo.New()
	repo := handler.SetupTestRepository(t)
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: &config.Config{}})

	// Add a custom category
	customCategory := domain.CategoryData{
		ID:     "tech-startups",
		Name:   "Tech Startups",
		Active: true,
	}
	_ = repo.SaveCategory(context.Background(), customCategory)

	// Create listings with "Business" category
	_ = repo.Save(context.Background(), domain.Listing{ID: "l1", Title: "L1", Type: domain.Business, City: "Lagos", OwnerOrigin: "Nigeria"})
	_ = repo.Save(context.Background(), domain.Listing{ID: "l2", Title: "L2", Type: domain.Business, City: "Accra", OwnerOrigin: "Ghana"})

	form := url.Values{}
	form.Add("action", "change_category")
	form.Add("selectedListings", "l1")
	form.Add("selectedListings", "l2")
	// The frontend now sends the Category Name as the value
	form.Add("new_category", customCategory.Name)

	req := httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		l1, _ := repo.FindByID(context.Background(), "l1")
		l2, _ := repo.FindByID(context.Background(), "l2")
		// Assert that the listing type is now the Category Name
		assert.Equal(t, domain.Category("Tech Startups"), l1.Type)
		assert.Equal(t, domain.Category("Tech Startups"), l2.Type)
	}
}

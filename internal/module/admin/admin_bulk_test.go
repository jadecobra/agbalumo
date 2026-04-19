package admin_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleBulkUpload(t *testing.T) {
	t.Parallel()
	// CSV headers: title,type,description,origin,email,phone,whatsapp
	csvContent := "title,type,description,origin,email,phone,whatsapp\nTest Biz,Business,Description,Nigeria,test@test.com,,"

	body, contentType := testutil.SetupCSVUploadBody(t, "csv_file", "upload.csv", csvContent)

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	env.App.CSVService = service.NewCSVService()
	h := admin.NewAdminHandler(env.App)
	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/upload", body)
	c.Request().Header.Set(echo.HeaderContentType, contentType)

	err := h.HandleBulkUpload(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify listing was saved
	testutil.AssertListingExists(t, env.App.DB, "Test Biz")
}

func TestAdminHandler_HandleBulkUpload_InvalidCSV(t *testing.T) {
	t.Parallel()
	// Junk content
	body, contentType := testutil.SetupCSVUploadBody(t, "csv_file", "junk.csv", "invalid,csv,data\n1,2,3")

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	env.App.CSVService = service.NewCSVService()
	h := admin.NewAdminHandler(env.App)
	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/upload", body)
	c.Request().Header.Set(echo.HeaderContentType, contentType)
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestAdminHandler_HandleBulkUpload_NoFile(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)
	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/upload", nil)
	_ = h.HandleBulkUpload(c)
	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestAdminHandler_HandleBulkUpload_ParseError(t *testing.T) {
	t.Parallel()
	// Missing required "description"
	csvContent := "title,type,origin,email\nTest Biz,Business,Nigeria,test@test.com"

	body, contentType := testutil.SetupCSVUploadBody(t, "csv_file", "upload.csv", csvContent)

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	env.App.CSVService = service.NewCSVService()
	h := admin.NewAdminHandler(env.App)
	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/upload", body)
	c.Request().Header.Set(echo.HeaderContentType, contentType)
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	// Verify no listing was saved
	listings, _, _ := env.App.DB.FindAll(context.Background(), "", "", "", "", "", true, 10, 0)
	assert.Empty(t, listings)
}

func TestHandleBulkAction_NoSelection(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/listings/bulk", nil)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
	}
}

func TestHandleBulkAction_StatusChanges(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		action       string
		expectStatus domain.ListingStatus
	}{
		{
			name:         "Approve",
			action:       "approve",
			expectStatus: domain.ListingStatusApproved,
		},
		{
			name:         "Reject",
			action:       "reject",
			expectStatus: domain.ListingStatusRejected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()
			h := admin.NewAdminHandler(env.App)

			_ = env.App.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "L1", Status: domain.ListingStatusPending})

			form := url.Values{}
			form.Add("action", tt.action)
			form.Add("selectedListings", "l1")
			c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			if assert.NoError(t, h.HandleBulkAction(c)) {
				assert.Equal(t, http.StatusFound, rec.Code)
				l, _ := env.App.DB.FindByID(context.Background(), "l1")
				assert.Equal(t, tt.expectStatus, l.Status)
			}
		})
	}
}

func TestHandleBulkAction_Delete(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	form := url.Values{}
	form.Add("action", "delete")
	form.Add("selectedListings", "l1")
	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Contains(t, rec.Header().Get("Location"), "/admin/listings/delete-confirm?id=l1")
	}
}

func TestHandleBulkAction_Categories(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		setup          func(ctx context.Context, db domain.ListingRepository) domain.Category
		expectedTarget domain.Category
	}{
		{
			name:           "StandardCategory",
			setup:          func(ctx context.Context, db domain.ListingRepository) domain.Category { return domain.Job },
			expectedTarget: domain.Job,
		},
		{
			name: "CustomCategory",
			setup: func(ctx context.Context, db domain.ListingRepository) domain.Category {
				cat := domain.CategoryData{ID: "tech", Name: "Tech", Active: true}
				_ = db.SaveCategory(ctx, cat)
				return domain.Category(cat.Name)
			},
			expectedTarget: domain.Category("Tech"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()
			h := admin.NewAdminHandler(env.App)
			ctx := context.Background()

			targetCat := tt.setup(ctx, env.App.DB)

			_ = env.App.DB.Save(ctx, domain.Listing{ID: "l1", Title: "L1", Type: domain.Business, City: "Lagos", OwnerOrigin: "Nigeria"})
			_ = env.App.DB.Save(ctx, domain.Listing{ID: "l2", Title: "L2", Type: domain.Business, City: "Accra", OwnerOrigin: "Ghana"})

			form := url.Values{}
			form.Add("action", "change_category")
			form.Add("selectedListings", "l1")
			form.Add("selectedListings", "l2")
			form.Add("new_category", string(targetCat))

			c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			if assert.NoError(t, h.HandleBulkAction(c)) {
				assert.Equal(t, http.StatusFound, rec.Code)
				l1, _ := env.App.DB.FindByID(ctx, "l1")
				_ = "admin_bulk_diff"
				l2, _ := env.App.DB.FindByID(ctx, "l2")
				assert.Equal(t, tt.expectedTarget, l1.Type)
				assert.Equal(t, tt.expectedTarget, l2.Type)
			}
		})
	}
}

func TestHandleBulkAction_NoAction(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/bulk", strings.NewReader("ids[]=1"))

	_ = h.HandleBulkAction(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestAdminHandler_HandleBulkUpload_ManyErrors(t *testing.T) {
	t.Parallel()
	csvContent := "title,type,description\n,,\n,,\n,,\n,,\n"
	body, contentType := testutil.SetupCSVUploadBody(t, "csv_file", "junk.csv", csvContent)

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	env.App.CSVService = service.NewCSVService()
	h := admin.NewAdminHandler(env.App)
	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/upload", body)
	c.Request().Header.Set(echo.HeaderContentType, contentType)
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

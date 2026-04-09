package admin_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleBulkUpload(t *testing.T) {
	t.Parallel()
	// CSV headers: title,type,description,origin,email,phone,whatsapp
	csvContent := "title,type,description,origin,email,phone,whatsapp\nTest Biz,Business,Description,Nigeria,test@test.com,,"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("csv_file", "upload.csv")
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	app, h, c, rec, cleanup := setupAdminBulkTest(t, http.MethodPost, "/admin/upload", body)
	defer cleanup()
	c.Request().Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	err := h.HandleBulkUpload(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify listing was saved
	listings, _ := app.DB.FindByTitle(context.Background(), "Test Biz")
	assert.Len(t, listings, 1)
	assert.Equal(t, "Test Biz", listings[0].Title)
}

func TestAdminHandler_HandleBulkUpload_InvalidCSV(t *testing.T) {
	t.Parallel()
	// Junk content
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("csv_file", "junk.csv")
	_, _ = part.Write([]byte("invalid,csv,data\n1,2,3"))
	_ = writer.Close()

	_, h, c, rec, cleanup := setupAdminBulkTest(t, http.MethodPost, "/admin/upload", body)
	defer cleanup()
	c.Request().Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestAdminHandler_HandleBulkUpload_NoFile(t *testing.T) {
	t.Parallel()
	_, h, c, rec, cleanup := setupAdminBulkTest(t, http.MethodPost, "/admin/upload", nil)
	defer cleanup()
	_ = h.HandleBulkUpload(c)
	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestAdminHandler_HandleBulkUpload_ParseError(t *testing.T) {
	t.Parallel()
	// Missing required "description"
	csvContent := "title,type,origin,email\nTest Biz,Business,Nigeria,test@test.com"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("csv_file", "upload.csv")
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	app, h, c, rec, cleanup := setupAdminBulkTest(t, http.MethodPost, "/admin/upload", body)
	defer cleanup()
	c.Request().Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	// Verify no listing was saved
	listings, _, _ := app.DB.FindAll(context.Background(), "", "", "", "", true, 10, 0)
	assert.Empty(t, listings)
}

func TestHandleBulkAction_NoSelection(t *testing.T) {
	t.Parallel()
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/bulk", nil)
	setupAdminAuth(t, c)
	_, h, cleanup := setupAdminTest(t)
	defer cleanup()

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
	}
}

func TestHandleBulkAction_Approve(t *testing.T) {
	t.Parallel()
	app, h, cleanup := setupAdminTest(t)
	defer cleanup()

	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "L1", Status: domain.ListingStatusPending})

	form := url.Values{}
	form.Add("action", "approve")
	form.Add("selectedListings", "l1")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	setupAdminAuth(t, c)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		l, _ := app.DB.FindByID(context.Background(), "l1")
		assert.Equal(t, domain.ListingStatusApproved, l.Status)
	}
}

func TestHandleBulkAction_Reject(t *testing.T) {
	t.Parallel()
	app, h, cleanup := setupAdminTest(t)
	defer cleanup()

	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "L1", Status: domain.ListingStatusPending})

	form := url.Values{}
	form.Add("action", "reject")
	form.Add("selectedListings", "l1")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	setupAdminAuth(t, c)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		l, _ := app.DB.FindByID(context.Background(), "l1")
		assert.Equal(t, domain.ListingStatusRejected, l.Status)
	}
}

func TestHandleBulkAction_Delete(t *testing.T) {
	t.Parallel()
	_, h, cleanup := setupAdminTest(t)
	defer cleanup()

	form := url.Values{}
	form.Add("action", "delete")
	form.Add("selectedListings", "l1")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	setupAdminAuth(t, c)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Contains(t, rec.Header().Get("Location"), "/admin/listings/delete-confirm?id=l1")
	}
}

func TestHandleBulkAction_ChangeCategory(t *testing.T) {
	t.Parallel()
	app, h, cleanup := setupAdminTest(t)
	defer cleanup()

	// Create listings with "Business" category
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "L1", Type: domain.Business, City: "Lagos", OwnerOrigin: "Nigeria"})
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l2", Title: "L2", Type: domain.Business, City: "Accra", OwnerOrigin: "Ghana"})

	form := url.Values{}
	form.Add("action", "change_category")
	form.Add("selectedListings", "l1")
	form.Add("selectedListings", "l2")
	form.Add("new_category", string(domain.Job))

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	setupAdminAuth(t, c)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		l1, _ := app.DB.FindByID(context.Background(), "l1")
		l2, _ := app.DB.FindByID(context.Background(), "l2")
		assert.Equal(t, domain.Job, l1.Type)
		assert.Equal(t, domain.Job, l2.Type)
	}
}

func TestHandleBulkAction_NoAction(t *testing.T) {
	t.Parallel()
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/bulk", strings.NewReader("ids[]=1"))
	setupAdminAuth(t, c)

	_, h, cleanup := setupAdminTest(t)
	defer cleanup()
	_ = h.HandleBulkAction(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestAdminHandler_HandleBulkUpload_ManyErrors(t *testing.T) {
	t.Parallel()
	// CSV with 4 invalid rows (missing title/desc)
	csvContent := "title,type,description\n,,\n,,\n,,\n,,\n"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("csv_file", "junk.csv")
	_, _ = part.Write([]byte(csvContent))
	_ = writer.Close()

	_, h, c, rec, cleanup := setupAdminBulkTest(t, http.MethodPost, "/admin/upload", body)
	defer cleanup()
	c.Request().Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	_ = h.HandleBulkUpload(c)

	assert.Equal(t, http.StatusFound, rec.Code)
}

func TestHandleBulkAction_ChangeToCustomCategory(t *testing.T) {
	t.Parallel()
	app, h, cleanup := setupAdminTest(t)
	defer cleanup()

	// Add a custom category
	customCategory := domain.CategoryData{
		ID:     "tech-startups",
		Name:   "Tech Startups",
		Active: true,
	}
	_ = app.DB.SaveCategory(context.Background(), customCategory)

	// Create listings with "Business" category
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "L1", Type: domain.Business, City: "Lagos", OwnerOrigin: "Nigeria"})
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l2", Title: "L2", Type: domain.Business, City: "Accra", OwnerOrigin: "Ghana"})

	form := url.Values{}
	form.Add("action", "change_category")
	form.Add("selectedListings", "l1")
	form.Add("selectedListings", "l2")
	// The frontend now sends the Category Name as the value
	form.Add("new_category", customCategory.Name)

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/bulk", strings.NewReader(form.Encode()))
	setupAdminAuth(t, c)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	if assert.NoError(t, h.HandleBulkAction(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		l1, _ := app.DB.FindByID(context.Background(), "l1")
		l2, _ := app.DB.FindByID(context.Background(), "l2")
		// Assert that the listing type is now the Category Name
		assert.Equal(t, domain.Category("Tech Startups"), l1.Type)
		assert.Equal(t, domain.Category("Tech Startups"), l2.Type)
	}
}

package admin_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleExportListings(t *testing.T) {
	t.Parallel()

	// t.Parallel() requirement: subtests that call t.Parallel() must have their own setup 
	// OR the parent must not return until they are done. 
	// Moving setup inside subtests for isolation.

	t.Run("Unauthorized access", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/admin/listings/export", nil)
		rec := httptest.NewRecorder()
		_ = echo.New().NewContext(req, rec)
		// No user set, will fail in middleware (not actually tested here as we call handler directly)
	})

	t.Run("Successful Export", func(t *testing.T) {
		t.Parallel()
		app, cleanup := testutil.SetupTestAppEnv(t)
		defer cleanup()

		h := admin.NewAdminHandler(app)
		app.CSVService = service.NewCSVService()

		ctx := context.Background()
		// Seed some data
		_ = app.DB.Save(ctx, domain.Listing{
			ID:           "test-1",
			Title:        "Test Listing",
			Type:         domain.Business,
			Description:  "Desc",
			OwnerOrigin:  "Nigeria",
			ContactEmail: "test@example.com",
			City:         "Lagos",
			IsActive:     true,
			Status:       domain.ListingStatusApproved,
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/listings/export", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.HandleExportListings(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "text/csv", rec.Header().Get(echo.HeaderContentType))
		assert.Contains(t, rec.Header().Get(echo.HeaderContentDisposition), "attachment")
		assert.Contains(t, rec.Header().Get(echo.HeaderContentDisposition), "listings.csv")

		body, _ := io.ReadAll(rec.Body)
		content := string(body)

		// Check headers
		assert.Contains(t, content, "ID,Title,Type,Description,City,Address")
		// Check data
		assert.Contains(t, content, "test-1,Test Listing,Business,Desc,Lagos")
		assert.Contains(t, content, "test@example.com")
	})
}

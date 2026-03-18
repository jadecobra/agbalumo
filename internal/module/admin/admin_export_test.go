package admin_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleExportListings(t *testing.T) {
	e := echo.New()
	repo := handler.SetupTestRepository(t)
	csvSvc := service.NewCSVService()
	h := admin.NewAdminHandler(repo, repo, repo, repo, repo, repo, repo, csvSvc, nil)

	ctx := context.Background()
	// Seed some data
	_ = repo.Save(ctx, domain.Listing{
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

	t.Run("Unauthorized access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/listings/export", nil)
		rec := httptest.NewRecorder()
		_ = e.NewContext(req, rec)
	})

	t.Run("Successful Export", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/listings/export", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

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

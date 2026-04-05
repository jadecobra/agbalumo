package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdate_JobSuccess(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	// Existing Job Listing
	existingListing := domain.Listing{
		ID:           "job-1",
		OwnerID:      "owner-1",
		Type:         domain.Job,
		Title:        "Senior Go Dev",
		Description:  "Write Go code",
		Company:      "Tech Corp",
		Skills:       "Go, SQL",
		PayRange:     "100k-150k",
		JobStartDate: time.Now().Add(24 * time.Hour),
		JobApplyURL:  "https://example.com",
		City:         "Lagos",
		OwnerOrigin:  "Nigeria",
		CreatedAt:    time.Now(),
		IsActive:     true,
		ContactEmail: "old@example.com",
	}
	_ = app.DB.Save(context.Background(), existingListing)

	jobStart := time.Now().Add(48 * time.Hour).Format("2006-01-02T15:04")

	h := listmod.NewListingHandler(app)

	formData := "title=Senior+Go+Dev+Updated&type=Job&owner_origin=Nigeria&description=Updated+Desc&contact_email=job@example.com&city=Lagos" +
		"&company=Updated+Corp&skills=Go,+Rust&pay_range=200k&job_apply_url=https://updated.com&job_start_date=" + jobStart

	c, rec := setupTestContext(http.MethodPost, "/listings/job-1", strings.NewReader(formData))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("job-1")
	c.Set("User", domain.User{ID: "owner-1"})

	// Execute
	err := h.HandleUpdate(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify DB state
	updated, _ := app.DB.FindByID(context.Background(), "job-1")
	assert.Equal(t, "Senior Go Dev Updated", updated.Title)
	assert.Equal(t, "Updated Corp", updated.Company)
	assert.Equal(t, "Go, Rust", updated.Skills)
	assert.Equal(t, "200k", updated.PayRange)
}

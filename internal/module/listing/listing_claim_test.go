package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleClaim(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "1", Title: "Biz", Type: domain.Business, Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", City: "Lagos", Address: "123 St"})
	_ = app.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})

	c, rec := setupTestContext(http.MethodPost, "/listings/1/claim", nil)
	c.SetPath("/listings/:id/claim")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "claimer", Name: "Claimer", Email: "c@e.com"})

	h := listmod.NewListingHandler(app)
	_ = h.HandleClaim(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

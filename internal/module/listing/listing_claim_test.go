package listing_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestHandleClaim(t *testing.T) {
	h, app, cleanup := setupListingHandler(t)
	defer cleanup()
	saveTestListing(t, app.DB, "1", "Biz")
	_ = app.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})

	c, rec := setupTestContext(http.MethodPost, "/listings/1/claim", nil)
	c.SetPath("/listings/:id/claim")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", newTestUser("claimer", domain.UserRoleUser))

	_ = h.HandleClaim(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

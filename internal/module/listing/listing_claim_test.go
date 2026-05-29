package listing_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleClaim(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	testutil.SaveTestListing(t, env.App.DB, "1", "Biz")
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})

	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/1/claim", nil)
	c.SetPath("/listings/:id/claim")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "claimer", Role: domain.UserRoleUser})

	_ = h.HandleClaim(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleClaim_Anonymous(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	testutil.SaveTestListing(t, env.App.DB, "1", "Biz")
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})

	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/1/claim", nil)
	c.SetPath("/listings/:id/claim")
	c.SetParamNames("id")
	c.SetParamValues("1")
	// No User set in context

	_ = h.HandleClaim(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Contains(t, rec.Header().Get("Location"), "/auth/google/login")
}

func TestHandleClaim_AlreadyClaimed(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)

	// Create a listing that already has an OwnerID
	_ = env.App.DB.Save(context.Background(), domain.Listing{
		ID:       "1",
		Title:    "Biz",
		Type:     domain.Business,
		OwnerID:  "original_owner",
		IsActive: true,
	})
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})

	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/1/claim", nil)
	c.SetPath("/listings/:id/claim")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "new_claimer", Role: domain.UserRoleUser})

	_ = h.HandleClaim(c)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandleClaim_PendingClaimExists(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)

	testutil.SaveTestListing(t, env.App.DB, "1", "Biz")
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})

	// Seed an existing pending claim request
	_ = env.App.DB.SaveClaimRequest(context.Background(), domain.ClaimRequest{
		ID:        "cr1",
		UserID:    "claimer",
		ListingID: "1",
		Status:    domain.ClaimStatusPending,
	})

	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/1/claim", nil)
	c.SetPath("/listings/:id/claim")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "claimer", Role: domain.UserRoleUser})

	_ = h.HandleClaim(c)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

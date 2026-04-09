package admin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleApproveClaim(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := admin.NewAdminHandler(app)

	// Seed data
	_ = app.DB.SaveClaimRequest(context.Background(), domain.ClaimRequest{ID: "claim1", UserID: "u1", ListingID: "l1", Status: domain.ClaimStatusPending})

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/claim1/approve", nil)
	c.SetParamNames("id")
	c.SetParamValues("claim1")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	if assert.NoError(t, h.HandleApproveClaim(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify database state
		claim, err := app.DB.GetClaimRequestByUserAndListing(context.Background(), "u1", "l1")
		assert.NoError(t, err)
		assert.Equal(t, domain.ClaimStatusApproved, claim.Status)
	}
}

func TestHandleApproveClaim_Error(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := admin.NewAdminHandler(app)

	// No data seeded for "bad" ID should result in error when trying to update
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/bad/approve", nil)
	c.SetParamNames("id")
	c.SetParamValues("bad")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	_ = h.HandleApproveClaim(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleRejectClaim(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := admin.NewAdminHandler(app)

	// Seed data
	_ = app.DB.SaveClaimRequest(context.Background(), domain.ClaimRequest{ID: "claim1", UserID: "u1", ListingID: "l1", Status: domain.ClaimStatusPending})

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/claim1/reject", nil)
	c.SetParamNames("id")
	c.SetParamValues("claim1")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	if assert.NoError(t, h.HandleRejectClaim(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify database state
		claim, err := app.DB.GetClaimRequestByUserAndListing(context.Background(), "u1", "l1")
		assert.NoError(t, err)
		assert.Equal(t, domain.ClaimStatusRejected, claim.Status)
	}
}

func TestHandleRejectClaim_Error(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := admin.NewAdminHandler(app)

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/bad/reject", nil)
	c.SetParamNames("id")
	c.SetParamValues("bad")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	_ = h.HandleRejectClaim(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

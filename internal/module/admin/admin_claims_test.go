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
	testClaimAction(t, "approve", domain.ClaimStatusPending, domain.ClaimStatusApproved, http.StatusOK)
}

func TestHandleApproveClaim_Error(t *testing.T) {
	testClaimAction(t, "approve", "", "", http.StatusNotFound)
}

func TestHandleRejectClaim(t *testing.T) {
	testClaimAction(t, "reject", domain.ClaimStatusPending, domain.ClaimStatusRejected, http.StatusOK)
}

func TestHandleRejectClaim_Error(t *testing.T) {
	testClaimAction(t, "reject", "", "", http.StatusNotFound)
}

func testClaimAction(t *testing.T, action string, initialStatus, expectedStatus domain.ClaimStatus, expectedCode int) {
	t.Helper()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	claimID := "claim1"
	if expectedCode == http.StatusNotFound {
		claimID = "bad"
	} else {
		// Seed data
		_ = env.App.DB.SaveClaimRequest(context.Background(), domain.ClaimRequest{
			ID:        claimID,
			UserID:    "u1",
			ListingID: "l1",
			Status:    initialStatus,
		})
	}

	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/claims/"+claimID+"/"+action, nil)
	c.SetParamNames("id")
	c.SetParamValues(claimID)

	var err error
	if action == "approve" {
		err = h.HandleApproveClaim(c)
	} else {
		err = h.HandleRejectClaim(c)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedCode, rec.Code)

	if expectedCode == http.StatusOK {
		// Verify database state
		claim, err := env.App.DB.GetClaimRequestByUserAndListing(context.Background(), "u1", "l1")
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, claim.Status)
	}
}

package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestHandleApproveClaim(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, nil, config.LoadConfig())

	// Seed data
	_ = repo.SaveClaimRequest(context.Background(), domain.ClaimRequest{ID: "claim1", UserID: "u1", ListingID: "l1", Status: domain.ClaimStatusPending})

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/claim1/approve", nil)
	c.SetParamNames("id")
	c.SetParamValues("claim1")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	if assert.NoError(t, h.HandleApproveClaim(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify database state
		claim, err := repo.GetClaimRequestByUserAndListing(context.Background(), "u1", "l1")
		assert.NoError(t, err)
		assert.Equal(t, domain.ClaimStatusApproved, claim.Status)
	}
}

func TestHandleApproveClaim_Error(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, nil, config.LoadConfig())

	// No data seeded for "bad" ID should result in error when trying to update
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/bad/approve", nil)
	c.SetParamNames("id")
	c.SetParamValues("bad")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	_ = h.HandleApproveClaim(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleRejectClaim(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, nil, config.LoadConfig())

	// Seed data
	_ = repo.SaveClaimRequest(context.Background(), domain.ClaimRequest{ID: "claim1", UserID: "u1", ListingID: "l1", Status: domain.ClaimStatusPending})

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/claim1/reject", nil)
	c.SetParamNames("id")
	c.SetParamValues("claim1")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	if assert.NoError(t, h.HandleRejectClaim(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify database state
		claim, err := repo.GetClaimRequestByUserAndListing(context.Background(), "u1", "l1")
		assert.NoError(t, err)
		assert.Equal(t, domain.ClaimStatusRejected, claim.Status)
	}
}

func TestHandleRejectClaim_Error(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	h := handler.NewAdminHandler(repo, nil, config.LoadConfig())

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/bad/reject", nil)
	c.SetParamNames("id")
	c.SetParamValues("bad")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	_ = h.HandleRejectClaim(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

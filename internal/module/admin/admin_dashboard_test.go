package admin_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleDashboard(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	setupAdminAuth(c)
	_, h, cleanup := setupAdminTest(t)
	defer cleanup()

	err := h.HandleDashboard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_HandleDashboard_HappyPath(t *testing.T) {
	app, h, cleanup := setupAdminTest(t)
	defer cleanup()
	ctx := context.Background()

	// 1. Users
	_ = app.DB.SaveUser(ctx, domain.User{ID: "u1", GoogleID: "g1", Role: domain.UserRoleAdmin})
	_ = app.DB.SaveUser(ctx, domain.User{ID: "u2", GoogleID: "g2", Role: domain.UserRoleUser})

	// 2. Listings
	_ = app.DB.Save(ctx, domain.Listing{ID: "l1", Title: "Business A", Type: domain.Business, IsActive: true, OwnerOrigin: "Nigeria", Address: "123 Lagos St"})
	_ = app.DB.Save(ctx, domain.Listing{ID: "l2", Title: "Job B", Type: domain.Job, IsActive: true, OwnerOrigin: "Ghana"})

	// 3. Claim Requests
	_ = app.DB.SaveClaimRequest(ctx, domain.ClaimRequest{ID: "c1", UserID: "u2", ListingID: "l1", Status: domain.ClaimStatusPending})

	// 4. Feedback
	_ = app.DB.SaveFeedback(ctx, domain.Feedback{ID: "f1", Type: domain.FeedbackTypeIssue, Content: "Help!"})

	// 5. Categories
	_ = app.DB.SaveCategory(ctx, domain.CategoryData{ID: "music", Name: "Music"})

	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	setupAdminAuth(c)

	err := h.HandleDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_HandleDashboard_GrowthMetrics(t *testing.T) {
	app, h, cleanup := setupAdminTest(t)
	defer cleanup()
	ctx := context.Background()

	// Seed multiple users/listings on different days if possible, or just enough to show it doesn't crash
	now := time.Now()
	_ = app.DB.SaveUser(ctx, domain.User{ID: "u1", CreatedAt: now})
	_ = app.DB.Save(ctx, domain.Listing{ID: "l1", CreatedAt: now})

	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	setupAdminAuth(c)

	err := h.HandleDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

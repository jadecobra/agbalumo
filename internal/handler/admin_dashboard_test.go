package handler_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleDashboard_HappyPath(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	// Seed data for various components of the dashboard
	ctx := context.Background()

	// 1. Users
	_ = repo.SaveUser(ctx, domain.User{ID: "u1", GoogleID: "g1", Role: domain.UserRoleAdmin})
	_ = repo.SaveUser(ctx, domain.User{ID: "u2", GoogleID: "g2", Role: domain.UserRoleUser})

	// 2. Listings
	_ = repo.Save(ctx, domain.Listing{ID: "l1", Title: "Business A", Type: domain.Business, IsActive: true, OwnerOrigin: "Nigeria", Address: "123 Lagos St"})
	_ = repo.Save(ctx, domain.Listing{ID: "l2", Title: "Job B", Type: domain.Job, IsActive: true, OwnerOrigin: "Ghana"})

	// 3. Claim Requests
	_ = repo.SaveClaimRequest(ctx, domain.ClaimRequest{ID: "c1", UserID: "u2", ListingID: "l1", Status: domain.ClaimStatusPending})

	// 4. Feedback
	_ = repo.SaveFeedback(ctx, domain.Feedback{ID: "f1", Type: domain.FeedbackTypeIssue, Content: "Help!"})

	// 5. Categories
	_ = repo.SaveCategory(ctx, domain.CategoryData{ID: "music", Name: "Music"})

	h := handler.NewAdminHandler(repo, repo, repo, repo, repo, repo, repo, nil, config.LoadConfig())
	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	err := h.HandleDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_HandleDashboard_GrowthMetrics(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	ctx := context.Background()

	// Seed multiple users/listings on different days if possible, or just enough to show it doesn't crash
	now := time.Now()
	_ = repo.SaveUser(ctx, domain.User{ID: "u1", CreatedAt: now})
	_ = repo.Save(ctx, domain.Listing{ID: "l1", CreatedAt: now})

	h := handler.NewAdminHandler(repo, repo, repo, repo, repo, repo, repo, nil, config.LoadConfig())
	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	err := h.HandleDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

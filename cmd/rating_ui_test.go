package cmd_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestRatingLegibility(t *testing.T) {
	// Setup: Manually seed a listing with a rating to ensure it renders
	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "file:test_ui.db?mode=memory&cache=shared"
	}

	// We need to import sqlite package
	repo, err := sqlite.NewSQLiteRepository(dbURL)
	assert.NoError(t, err)
	defer func() { _ = repo.Close() }()

	_ = repo.Save(ctx, domain.Listing{
		ID:          "test-rating-id",
		Title:       "Rating Test Venue",
		ReviewCount: 10,
		Rating:      4.5,
		Type:        domain.Food,
		IsActive:    true,
		Status:      domain.ListingStatusApproved,
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()

	// RED TEST: Currently listings use a single ⭐.
	// We want to see 5 🍊 icons representing the scale.
	// We check for at least 5 instances of the Agbalumo icon (🍊).
	// (Excluding the featured badge and footer which also use 🍊)
	count := strings.Count(body, "🍊")
	assert.GreaterOrEqual(t, count, 5, "Expected at least 5 Agbalumo icons (🍊) to represent the 5-star rating scale")

	// Check for legibility classes (no text-[10px] for ratings)
	assert.NotContains(t, body, `text-[10px] md:text-xs text-yellow-500 font-bold flex items-center gap-0.5" data-testid="listing-rating"`)
}

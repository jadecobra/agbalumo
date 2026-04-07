package listing_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestContext(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	return testutil.SetupTestContext(method, target, body)
}

func setupRequest(method, target string, body io.Reader) *http.Request {
	return httptest.NewRequest(method, target, body)
}

func setupResponseRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func saveTestListing(t *testing.T, db domain.ListingRepository, id, title string, extra ...func(*domain.Listing)) {
	t.Helper()
	l := domain.Listing{
		ID:           id,
		Title:        title,
		Type:         domain.Business,
		Status:       domain.ListingStatusApproved,
		IsActive:     true,
		Address:      "Lagos",
		ContactEmail: "test@example.com",
		OwnerOrigin:  "Nigeria",
	}
	for _, f := range extra {
		f(&l)
	}
	err := db.Save(context.Background(), l)
	if err != nil {
		t.Fatalf("Failed to save test listing %s: %v", id, err)
	}
}

func assertFeaturedStatus(t *testing.T, db domain.ListingRepository, id string, expected bool) {
	t.Helper()
	l, err := db.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, expected, l.Featured, "Listing %s featured status mismatch", id)
}

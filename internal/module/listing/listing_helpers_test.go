package listing_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
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
	testutil.SaveTestListing(t, db, id, title, extra...)
}

func newTestUser(id string, role domain.UserRole) domain.User {
	return domain.User{
		ID:   id,
		Role: role,
	}
}

func newTestListing(id, title string, extra ...func(*domain.Listing)) domain.Listing {
	l := domain.Listing{
		ID:           id,
		Title:        title,
		Type:         domain.Business,
		Status:       domain.ListingStatusApproved,
		IsActive:     true,
		Address:      "Lagos",
		ContactEmail: "test@example.com",
		OwnerOrigin:  "Nigeria",
		CreatedAt:    time.Now(),
	}
	for _, f := range extra {
		f(&l)
	}
	return l
}

func setupListingHandler(t *testing.T) (*listmod.ListingHandler, *env.AppEnv, func()) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	h := listmod.NewListingHandler(app)
	return h, app, cleanup
}

func assertContainsPagination(t testing.TB, body string) {
	testutil.AssertContainsPagination(t, body)
}

func assertErrorPage(t testing.TB, body string) {
	testutil.AssertErrorPage(t, body)
}

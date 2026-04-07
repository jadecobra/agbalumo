package listing_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
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


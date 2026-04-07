package listing_test

import (
	"io"
	"net/http"
	"net/http/httptest"

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

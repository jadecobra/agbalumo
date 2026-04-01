package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHelloAgent(t *testing.T) {
	e := echo.New()

	// In a real scenario, we would call the function that registers routes.
	// But in the RED phase, that function doesn't exist or doesn't have the route yet.

	req := httptest.NewRequest(http.MethodGet, "/hello-agent", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.GET("/hello-agent", handler.HandleHelloAgent)
	e.ServeHTTP(rec, req)

	// Test that the handler returns the correct response when registered.
	assert.Equal(t, http.StatusOK, rec.Code, "Route /hello-agent should be registered and return 200 OK")

	expectedBody := `{"message":"Hello, Agent!"}`
	assert.JSONEq(t, expectedBody, rec.Body.String(), "Response body should match")
}

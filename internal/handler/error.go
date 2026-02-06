package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RespondError logs the actual error internally and renders a friendly error page to the user.
// This prevents sensitive details (like DB errors) from leaking to the client.
func RespondError(c echo.Context, err error) error {
	// Log the actual error for debugging
	log.Printf("[ERROR] Request ID: %s | Error: %v", c.Response().Header().Get(echo.HeaderXRequestID), err)

	// Render the friendly error page
	// We return 500 Internal Server Error, but the user sees the friendly message.
	return c.Render(http.StatusInternalServerError, "error.html", nil)
}

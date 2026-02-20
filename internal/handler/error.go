package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RespondError logs the actual error internally and renders a friendly error page to the user.
// This prevents sensitive details (like DB errors) from leaking to the client.
func RespondError(c echo.Context, err error) error {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	// Log the actual error for debugging
	slog.Error("Request failed",
		slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
		slog.Any("error", err),
	)

	// Render the friendly error page
	// We return the specific code (e.g. 400 or 500)
	return c.Render(code, "error.html", nil)
}

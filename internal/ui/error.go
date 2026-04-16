package ui

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
	var message string
	if he, ok := err.(*echo.HTTPError); ok {
		message = he.Message.(string)
	}

	return c.Render(code, "error.html", map[string]interface{}{
		"Message": message,
	})
}

// RespondErrorMsg is a convenient shorthand for returning an HTTPError with a specific message.
// This reduces code duplication across handlers.
func RespondErrorMsg(c echo.Context, code int, message string) error {
	return RespondError(c, echo.NewHTTPError(code, message))
}

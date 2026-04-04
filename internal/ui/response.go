package ui

import "github.com/labstack/echo/v4"

// ErrorResponse represents a standardized JSON error structure.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// RespondJSONError sends a standardized JSON error response.
func RespondJSONError(c echo.Context, code int, errMessage string) error {
	return c.JSON(code, &ErrorResponse{
		Error: errMessage,
		Code:  code,
	})
}

package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HandleHelloAgent returns a JSON hello message.
func HandleHelloAgent(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello, Agent!"})
}

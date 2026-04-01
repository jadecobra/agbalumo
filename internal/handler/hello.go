package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HelloResponse defines the JSON structure for the hello-agent endpoint.
type HelloResponse struct {
	Message string `json:"message"`
}

// HandleHelloAgent returns a JSON hello message.
func HandleHelloAgent(c echo.Context) error {
	return c.JSON(http.StatusOK, HelloResponse{
		Message: "Hello, Agent!",
	})
}

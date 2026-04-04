package cmd

import (
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/infra/server"
	"github.com/labstack/echo/v4"
)

// SetupServer initializes the Echo server and its dependencies by calling the infra layer.
func SetupServer() (*echo.Echo, error) {
	cfg := config.LoadConfig()
	return server.Setup(cfg)
}

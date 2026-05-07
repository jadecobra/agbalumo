package maintenance

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/server"
)

// GenerateRouteMap boots the server in dry-run mode and returns the Echo route table.
func GenerateRouteMap() (string, error) {
	// Backup original env
	origEnv := os.Getenv(domain.EnvKeyAppEnv)
	origDryRun := os.Getenv("AGBALUMO_DRY_RUN")

	_ = os.Setenv(domain.EnvKeyAppEnv, "test")
	_ = os.Setenv("AGBALUMO_DRY_RUN", "true")
	defer func() {
		_ = os.Setenv(domain.EnvKeyAppEnv, origEnv)
		_ = os.Setenv("AGBALUMO_DRY_RUN", origDryRun)
	}()

	cfg := config.LoadConfig()
	// Use in-memory DB to avoid side effects
	cfg.DatabaseURL = domain.SQLiteMemory

	e, cleanup, err := server.Setup(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to setup server: %w", err)
	}
	defer cleanup()

	var lines []string
	for _, r := range e.Routes() {
		if !isInternalRoute(r.Method, r.Path, r.Name) {
			lines = append(lines, fmt.Sprintf("%-6s %-30s →  %s", r.Method, r.Path, r.Name))
		}
	}

	sort.Strings(lines)
	return strings.Join(lines, "\n"), nil
}

func isInternalRoute(method, path, name string) bool {
	if method == "echo_route_not_found" {
		return true
	}
	if strings.Contains(name, "github.com/labstack/echo") && !strings.Contains(name, "agbalumo") {
		if path == "/*" || path == "/static/*" {
			return true
		}
	}
	return false
}

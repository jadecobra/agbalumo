package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetupServer(t *testing.T) {
	// go test ./cmd runs from project root, so template paths are already correct

	// Setup Env
	os.Setenv("AGBALUMO_ENV", "test")
	defer os.Unsetenv("AGBALUMO_ENV")

	// We need a dummy DB or allow it to fail?
	// SetupServer connects to DB.
	// If agbalumo.db exists, it works.
	// If not, it creates it.
	// Best to use a temp db file.
	os.Setenv("DATABASE_URL", "test_server.db")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Remove("test_server.db")
	}()

	e, err := SetupServer()
	require.NoError(t, err)
	require.NotNil(t, e)

	// Verify Routes
	// We can't easily inspect routes map in Echo without private access or iterating
	// But we can check if e.Renderer is set
	require.NotNil(t, e.Renderer)

	// Verify Middlewares (count)
	// SecureHeaders, RateLimit, CSRF, Session => 4 global middlewares?
	// Actually Echo has pre-middlewares and use-middlewares.
	// We can't strictly count without being fragile.
	require.NotEmpty(t, e.Routes())
}

func TestServeCmd_Run(t *testing.T) {
	// go test ./cmd runs from project root, so template paths are already correct

	os.Setenv("AGBALUMO_ENV", "test")
	os.Setenv("AGBALUMO_DRY_RUN", "true")
	os.Setenv("DATABASE_URL", "test_serve_cmd.db") // Diff DB
	defer func() {
		os.Unsetenv("AGBALUMO_ENV")
		os.Unsetenv("AGBALUMO_DRY_RUN")
		os.Unsetenv("DATABASE_URL")
		os.Remove("test_serve_cmd.db")
	}()

	// Execute Run
	// serveCmd is global in cmd package
	serveCmd.Run(serveCmd, []string{})
}

func TestSetupServerProduction(t *testing.T) {
	// Test SetupServer with production environment (JSON logging)
	os.Setenv("AGBALUMO_ENV", "production")
	os.Setenv("SESSION_SECRET", "production-secret-key")
	os.Setenv("DATABASE_URL", "test_server_prod.db")
	os.Setenv("ADMIN_CODE", "test-admin-code")
	defer func() {
		os.Unsetenv("AGBALUMO_ENV")
		os.Unsetenv("SESSION_SECRET")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("ADMIN_CODE")
		os.Remove("test_server_prod.db")
	}()

	e, err := SetupServer()
	require.NoError(t, err)
	require.NotNil(t, e)
}

func TestSetupBackgroundServicesProduction(t *testing.T) {
	// Test setupBackgroundServices in production (should skip seeding)
	os.Setenv("AGBALUMO_ENV", "production")
	os.Setenv("SESSION_SECRET", "test-secret-key")
	os.Setenv("DATABASE_URL", "test_bg_prod.db")
	os.Setenv("ADMIN_CODE", "test-admin-code")
	defer func() {
		os.Unsetenv("AGBALUMO_ENV")
		os.Unsetenv("SESSION_SECRET")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("ADMIN_CODE")
		os.Remove("test_bg_prod.db")
	}()

	e, err := SetupServer()
	require.NoError(t, err)
	require.NotNil(t, e)
	// Background services are started in goroutine, so we just verify server setup
}

package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthzEndpoint(t *testing.T) {
	_ = os.Setenv("AGBALUMO_ENV", "test")
	_ = os.Setenv("DATABASE_URL", "@tester/test_healthz.db")
	defer func() {
		_ = os.Unsetenv("AGBALUMO_ENV")
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Remove("@tester/test_healthz.db")
	}()

	e, err := SetupServer()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, "ok", body["status"])
}

func TestSetupServer(t *testing.T) {
	// go test ./cmd runs from project root, so template paths are already correct

	// Setup Env
	_ = os.Setenv("AGBALUMO_ENV", "test")
	defer func() { _ = os.Unsetenv("AGBALUMO_ENV") }()

	// We need a dummy DB or allow it to fail?
	// SetupServer connects to DB.
	// If agbalumo.db exists, it works.
	// If not, it creates it.
	// Best to use a temp db file.
	_ = os.Setenv("DATABASE_URL", "@tester/test_server.db")
	defer func() {
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Remove("@tester/test_server.db")
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

	_ = os.Setenv("AGBALUMO_ENV", "test")
	_ = os.Setenv("AGBALUMO_DRY_RUN", "true")
	_ = os.Setenv("DATABASE_URL", "@tester/test_serve_cmd.db") // Diff DB
	defer func() {
		_ = os.Unsetenv("AGBALUMO_ENV")
		_ = os.Unsetenv("AGBALUMO_DRY_RUN")
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Remove("@tester/test_serve_cmd.db")
	}()

	// Execute Run
	// serveCmd is global in cmd package
	serveCmd.Run(serveCmd, []string{})
}

func TestSetupServerProduction(t *testing.T) {
	// Test SetupServer with production environment (JSON logging)
	_ = os.Setenv("AGBALUMO_ENV", "production")
	_ = os.Setenv("SESSION_SECRET", "production-secret-key")
	_ = os.Setenv("DATABASE_URL", "@tester/test_server_prod.db")
	_ = os.Setenv("ADMIN_CODE", "test-admin-code")
	defer func() {
		_ = os.Unsetenv("AGBALUMO_ENV")
		_ = os.Unsetenv("SESSION_SECRET")
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Unsetenv("ADMIN_CODE")
		_ = os.Remove("@tester/test_server_prod.db")
	}()

	e, err := SetupServer()
	require.NoError(t, err)
	require.NotNil(t, e)
}

func TestSetupBackgroundServicesProduction(t *testing.T) {
	// Test setupBackgroundServices in production (should skip seeding)
	_ = os.Setenv("AGBALUMO_ENV", "production")
	_ = os.Setenv("SESSION_SECRET", "test-secret-key")
	_ = os.Setenv("DATABASE_URL", "@tester/test_bg_prod.db")
	_ = os.Setenv("ADMIN_CODE", "test-admin-code")
	defer func() {
		_ = os.Unsetenv("AGBALUMO_ENV")
		_ = os.Unsetenv("SESSION_SECRET")
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Unsetenv("ADMIN_CODE")
		_ = os.Remove("@tester/test_bg_prod.db")
	}()

	e, err := SetupServer()
	require.NoError(t, err)
	require.NotNil(t, e)
	// Background services are started in goroutine, so we just verify server setup
}

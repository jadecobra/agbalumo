package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetupServer(t *testing.T) {
	// Change CWD to project root for templates
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	// Revert CWD after test
	defer func() {
		os.Chdir("cmd")
	}()

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
	// Change CWD to project root for templates
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Chdir("cmd")
	}()

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

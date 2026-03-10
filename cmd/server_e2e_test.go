package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	// Setup temporary environment
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "e2e.db")

	// Set required env vars for SetupServer
	_ = os.Setenv("DATABASE_URL", dbPath)
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("SESSION_SECRET", "e2e-secret-key-123")
	_ = os.Setenv("UPLOAD_DIR", filepath.Join(tmpDir, "uploads"))
	defer func() {
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Unsetenv("ENV")
		_ = os.Unsetenv("SESSION_SECRET")
		_ = os.Unsetenv("UPLOAD_DIR")
	}()

	e, err := SetupServer()
	if err != nil {
		t.Fatalf("Failed to setup server: %v", err)
	}

	ts := httptest.NewServer(e)
	defer ts.Close()

	client := ts.Client()

	// 1. Health Check
	t.Run("Healthz", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/healthz")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 2. Home Page
	t.Run("Home Page", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "find what you want")
	})

	// 3. Search Fragment
	t.Run("Search Fragment", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/listings/fragment?q=test")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

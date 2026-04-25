package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/cmd"
	"github.com/stretchr/testify/assert"
)

/**
 * Search Latency Constraint
 * The search/filter endpoint MUST return a response in under 200ms
 * to ensure a smooth discovery experience for users like Ada.
 */
func TestSearchLatency_Constraint(t *testing.T) {
	// Change CWD to project root so relative paths for templates/DB work
	originalCwd, _ := os.Getwd()
	root := originalCwd
	for !strings.HasSuffix(root, "/agbalumo") && root != "/" {
		root = filepath.Dir(root)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("failed to change directory to root %s: %v", root, err)
	}
	defer func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Logf("failed to restore original cwd: %v", err)
		}
	}()

	t.Setenv("AGBALUMO_ENV", "test")

	// Initialize the server and dependencies
	e, cleanup, err := cmd.SetupServer()
	if err != nil {
		t.Fatalf("failed to setup server: %v", err)
	}
	defer cleanup()

	// Use a realistic search query that Ada might use
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment?q=Nigerian", nil)
	rec := httptest.NewRecorder()

	// Measure the end-to-end latency of the handler
	start := time.Now()
	e.ServeHTTP(rec, req)
	duration := time.Since(start)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code, "Search should return 200 OK")

	// Strict performance budget: 200ms
	// *Insight:* This forces efficient database indexing and parallel execution
	// in the handler (ListingHandler.HandleFragment).
	assert.Less(t, duration, 200*time.Millisecond, "Search API latency MUST be under 200ms budget")

	t.Logf("Search latency: %v", duration)
}

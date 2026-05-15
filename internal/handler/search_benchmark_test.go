package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/cmd"
)

func BenchmarkSearchLatency(b *testing.B) {
	// Change CWD to project root so relative paths for templates/DB work
	originalCwd, _ := os.Getwd()
	root := originalCwd
	for !strings.HasSuffix(root, "/agbalumo") && root != "/" {
		root = filepath.Dir(root)
	}
	if err := os.Chdir(root); err != nil {
		b.Fatalf("failed to change directory to root %s: %v", root, err)
	}
	defer func() {
		_ = os.Chdir(originalCwd)
	}()

	if err := os.Setenv("AGBALUMO_ENV", "test"); err != nil {
		b.Fatalf("failed to set env: %v", err)
	}

	// Initialize the server and dependencies
	e, cleanup, err := cmd.SetupServer()
	if err != nil {
		b.Fatalf("failed to setup server: %v", err)
	}
	defer cleanup()

	// Use a realistic search query
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment?q=Nigerian", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
	}
}

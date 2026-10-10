package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCICmdHasWithDockerFlag(t *testing.T) {
	flag := ciCmd.Flags().Lookup("with-docker")
	if flag == nil {
		t.Fatal("ciCmd should have a --with-docker flag")
	}
	if flag.DefValue != "false" {
		t.Errorf("expected default false, got %s", flag.DefValue)
	}
}

func TestUptimeCmdHasPathFlag(t *testing.T) {
	flag := uptimeCmd.Flags().Lookup("path")
	if flag == nil {
		t.Fatal("uptimeCmd should have a --path flag")
	}
	if flag.DefValue != "" {
		t.Errorf("expected default empty string, got %s", flag.DefValue)
	}
}

func TestUptimeCmdPathFlagExecution(t *testing.T) {
	var receivedPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	t.Setenv("APP_URL", ts.URL)
	defer func() {
		_ = uptimeCmd.Flags().Set("path", "")
	}()

	if err := uptimeCmd.Flags().Set("path", "/listings/1/og.png"); err != nil {
		t.Fatal(err)
	}

	if err := uptimeCmd.RunE(uptimeCmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedPath != "/listings/1/og.png" {
		t.Fatalf("expected path /listings/1/og.png, got %s", receivedPath)
	}
}

func TestRunTrivyScanFunctionExists(t *testing.T) {
	// Verify localCIImageTag constant exists (this will fail compilation initially)
	tag := localCIImageTag
	if tag == "" {
		t.Fatal("localCIImageTag constant must not be empty")
	}
}

func TestCICmdWithDockerFlagDescription(t *testing.T) {
	flag := ciCmd.Flags().Lookup("with-docker")
	if flag == nil {
		t.Fatal("--with-docker flag missing from ciCmd")
	}
	if !strings.Contains(flag.Usage, "trivy") {
		t.Errorf("--with-docker flag description should mention trivy; got: %s", flag.Usage)
	}
}

func TestBrowserCmdRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "browser" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("browser subcommand is not registered")
	}
}

func TestPipeDrainCmdRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "pipe-drain" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("pipe-drain subcommand is not registered")
	}
}

func TestComplexityCmdRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "complexity" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("complexity subcommand is not registered")
	}
}

func TestHandlerConcurrencyCmdRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "handler-concurrency" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("handler-concurrency subcommand is not registered")
	}
}

func TestErrorSwallowCmdRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "error-swallow" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("error-swallow subcommand is not registered")
	}
}

func TestGetVerificationOpts(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			t.Error(err)
		}
	}()

	tests := []struct {
		name         string
		expectedPath string
		setupFiles   []string
	}{
		{
			name:         "canonical path exists",
			setupFiles:   []string{".agents/coverage.json"},
			expectedPath: ".agents/coverage.json",
		},
		{
			name:         "canonical and fallback both exist",
			setupFiles:   []string{".agents/coverage.json", ".metrics/coverage"},
			expectedPath: ".agents/coverage.json",
		},
		{
			name:         "only fallback exists",
			setupFiles:   []string{".metrics/coverage"},
			expectedPath: ".metrics/coverage",
		},
		{
			name:         "none exist, default to fallback",
			setupFiles:   []string{},
			expectedPath: ".metrics/coverage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			performVerificationTest(t, tt.setupFiles, tt.expectedPath)
		})
	}
}

func performVerificationTest(t *testing.T, files []string, expected string) {
	subDir := t.TempDir()
	if err := os.Chdir(subDir); err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if err := os.MkdirAll(filepath.Dir(f), 0750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f, []byte("{}"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	cmd := &cobra.Command{}
	setupVerifyFlags(cmd)
	_, path := getVerificationOpts(cmd)

	if path != expected {
		t.Errorf("expected path %s, got %s", expected, path)
	}
}

func TestCICmdHasFocusFlag(t *testing.T) {
	t.Parallel()
	flag := ciCmd.Flags().Lookup("focus")
	if flag == nil {
		t.Fatal("ciCmd should have a --focus flag")
	}
	if flag.DefValue != "" {
		t.Errorf("expected default empty string, got %s", flag.DefValue)
	}
}

func TestDBPathCmdRegistered(t *testing.T) {
	t.Parallel()
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "db-path" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("db-path command should be registered in rootCmd")
	}
}

func TestAppendOrOverrideEnv(t *testing.T) {
	base := []string{"FOO=bar", "BAZ=qux", "KEY=val1"}
	overrides := []string{"KEY=val2", "NEW=entry"}

	result := appendOrOverrideEnv(base, overrides...)

	expected := []string{"FOO=bar", "BAZ=qux", "KEY=val2", "NEW=entry"}
	if len(result) != len(expected) {
		t.Fatalf("expected len %d, got %d: %v", len(expected), len(result), result)
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("at index %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestBuildCITasksIncludesTZMatrix(t *testing.T) {
	tasks := buildCITasks(ciCmd, nil)

	var hasUTC, hasChicago, hasOld bool
	for _, task := range tasks {
		switch task.Name {
		case "Running Heavy Tests (with -race, TZ=UTC)":
			hasUTC = true
		case "Running Heavy Tests (with -race, TZ=America/Chicago)":
			hasChicago = true
		case "Running Heavy Tests (with -race)":
			hasOld = true
		}
	}

	if !hasUTC {
		t.Error("missing task: 'Running Heavy Tests (with -race, TZ=UTC)'")
	}
	if !hasChicago {
		t.Error("missing task: 'Running Heavy Tests (with -race, TZ=America/Chicago)'")
	}
	if hasOld {
		t.Error("legacy task 'Running Heavy Tests (with -race)' should not be present")
	}
}

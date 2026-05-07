package main

import (
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

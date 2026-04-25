package main

import (
	"strings"
	"testing"
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


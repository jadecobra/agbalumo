package agent

import (
	"os"
	"testing"
)

func TestVerifyApiSpecGate(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommand
	defer func() { ExecCommand = orig }()

	VerifyApiSpec("feature")
}

func TestVerifyApiSpec_DriftAggregation(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommand
	defer func() { ExecCommand = orig }()

	// Mock file system in a temp dir
	tmpDir := t.TempDir()

	err := os.MkdirAll(tmpDir+"/cmd", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(tmpDir+"/docs/cli", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(tmpDir+"/internal/handler", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(tmpDir+"/internal/module", 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Code commands
	_ = os.WriteFile(tmpDir+"/cmd/main.go", []byte(`package main
import "github.com/spf13/cobra"
var cmd = &cobra.Command{ Use: "serve" }`), 0644)

	// OpenApi and Docs
	_ = os.WriteFile(tmpDir+"/docs/openapi.yaml", []byte(`
paths:
  /auth:
    get: {}`), 0644)

	_ = os.WriteFile(tmpDir+"/docs/api.md", []byte(`
| Method | Path |
|--------|------|
| GET | /missing |`), 0644)

	// MD Commands missing 'serve' and adding 'unknown'
	_ = os.WriteFile(tmpDir+"/docs/cli.md", []byte(`### unknown`), 0644)

	// Change dir to TempDir to trick hardcoded paths
	origWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(origWd) }()

	// The system shouldn't pass
	if VerifyApiSpec("feature") {
		t.Error("VerifyApiSpec should have failed due to drift violations")
	}
}

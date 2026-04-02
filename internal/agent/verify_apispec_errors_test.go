package agent

import (
	"os"
	"strings"
	"testing"
)

func TestVerifyApiSpec_ExtractRoutesError(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	// Use a temp dir with no code files at all to trigger ExtractRoutes error
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(origWd) }()

	output := captureStdout(t, func() {
		result := VerifyApiSpec("feature")
		if result {
			t.Error("expected VerifyApiSpec to fail when ExtractRoutes errors")
		}
	})

	if !strings.Contains(output, "Error extracting routes from code") {
		t.Errorf("expected output to contain 'Error extracting routes from code', got: %s", output)
	}
}

func TestVerifyApiSpec_OpenAPIBundleError(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpecFail
	defer func() { ExecCommand = orig }()

	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/labstack/echo/v4"
func routes(e *echo.Echo) {
	e.GET("/api/listings", nil)
}
`,
		openapiContent: `paths:
  /api/listings:
    get:
      summary: List
`,
		apiMdContent: `| Method | Path |`,
		cliMdContent: `### serve`,
	})
	defer cleanup()

	output := captureStdout(t, func() {
		result := VerifyApiSpec("feature")
		if result {
			t.Error("expected VerifyApiSpec to fail on OpenAPI bundle error")
		}
	})

	if !strings.Contains(output, "Error bundling docs/openapi.yaml") {
		t.Errorf("expected output to contain 'Error bundling docs/openapi.yaml', got: %s", output)
	}
}

func TestVerifyApiSpec_ApiMdMissing(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	// Create temp with everything except docs/api.md
	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/labstack/echo/v4"
func routes(e *echo.Echo) {
	e.GET("/api/listings", nil)
}
`,
		openapiContent: `paths:
  /api/listings:
    get:
      summary: List
`,
		cliMdContent: `### serve`,
	})
	defer cleanup()

	// Remove docs/api.md to simulate missing file
	_ = os.Remove("docs/api.md")

	output := captureStdout(t, func() {
		result := VerifyApiSpec("feature")
		if result {
			t.Error("expected VerifyApiSpec to fail when docs/api.md is missing")
		}
	})

	if !strings.Contains(output, "Error reading docs/api.md") {
		t.Errorf("expected output to contain 'Error reading docs/api.md', got: %s", output)
	}
}

func TestVerifyApiSpec_CLICodeExtractError(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		openapiContent: `paths:
  /api/listings:
    get:
      summary: List
`,
		apiMdContent: `| Method | Path |
|--------|------|
| GET | /api/listings |
`,
		cliMdContent: `### serve`,
	})
	defer cleanup()

	// Put a valid Go file in internal/handler so ExtractRoutes finds at least one file
	_ = os.WriteFile("internal/handler/routes.go", []byte(`package handler
import "github.com/labstack/echo/v4"
func RegisterRoutes(e *echo.Echo) {
	e.GET("/api/listings", nil)
}
`), 0644)

	// Make cmd directory unreadable to force filepath.Walk to fail
	_ = os.Chmod("cmd", 0000)
	defer func() { _ = os.Chmod("cmd", 0755) }()

	output := captureStdout(t, func() {
		result := VerifyApiSpec("feature")
		if result {
			t.Error("expected VerifyApiSpec to fail when CLI code extraction errors")
		}
	})

	if !strings.Contains(output, "Error extracting CLI code cmds") {
		t.Errorf("expected output to contain 'Error extracting CLI code cmds', got: %s", output)
	}
}

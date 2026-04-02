package agent

import (
	"os"
	"strings"
	"testing"
)

func TestVerifyApiSpec_APIDrift(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	// Code has routes that OpenAPI doesn't, and vice versa
	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/labstack/echo/v4"
func routes(e *echo.Echo) {
	e.GET("/code-only", nil)
}
`,
		openapiContent: `paths:
  /openapi-only:
    get:
      summary: OpenAPI only
`,
		apiMdContent: `| Method | Path | Description |
|--------|------|-------------|
| GET | /md-only | MD only |
`,
		cliMdContent: `### serve
Description
`,
	})
	defer cleanup()

	// Write the cmd Go file with CLI Use for "serve"
	_ = os.WriteFile("cmd/server.go", []byte(`package main
import "github.com/spf13/cobra"
var cmd = &cobra.Command{ Use: "serve" }
`), 0644)

	output := captureStdout(t, func() {
		result := VerifyApiSpec("feature")
		if result {
			t.Error("expected VerifyApiSpec to return false (drift detected)")
		}
	})

	if !strings.Contains(output, "Gate FAIL") {
		t.Errorf("expected output to contain 'Gate FAIL', got: %s", output)
	}
	if !strings.Contains(output, "contract drift detected") {
		t.Errorf("expected output to contain 'contract drift detected', got: %s", output)
	}
}

func TestVerifyApiSpec_CLIDrift(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/spf13/cobra"
var cmd = &cobra.Command{ Use: "serve" }
`,
		openapiContent: `paths: {}`,
		apiMdContent:   `| Method | Path |`,
		cliMdContent:   `### unknown-cmd`,
	})
	defer cleanup()

	output := captureStdout(t, func() {
		result := VerifyApiSpec("feature")
		if result {
			t.Error("expected VerifyApiSpec to return false (CLI drift)")
		}
	})

	if !strings.Contains(output, "Gate FAIL") {
		t.Errorf("expected output to contain 'Gate FAIL', got: %s", output)
	}
}

func TestVerifyApiSpec_WorkflowRefactor(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/spf13/cobra"
var cmd = &cobra.Command{ Use: "serve" }
`,
		openapiContent: `paths: {}`,
		apiMdContent:   `| Method | Path |`,
		cliMdContent:   `### unknown-cmd`,
	})
	defer cleanup()

	output := captureStdout(t, func() {
		result := VerifyApiSpec("refactor")
		if result {
			t.Error("expected VerifyApiSpec to fail for refactor workflow with drift")
		}
	})

	if !strings.Contains(output, "mandatory passive validations") {
		t.Errorf("expected output to contain 'mandatory passive validations', got: %s", output)
	}
}

func TestVerifyApiSpec_WorkflowBugfix(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/spf13/cobra"
var cmd = &cobra.Command{ Use: "serve" }
`,
		openapiContent: `paths: {}`,
		apiMdContent:   `| Method | Path |`,
		cliMdContent:   `### unknown-cmd`,
	})
	defer cleanup()

	output := captureStdout(t, func() {
		result := VerifyApiSpec("bugfix")
		if result {
			t.Error("expected VerifyApiSpec to fail for bugfix workflow with drift")
		}
	})

	if !strings.Contains(output, "mandatory passive validations") {
		t.Errorf("expected output to contain 'mandatory passive validations', got: %s", output)
	}
	if !strings.Contains(output, "bugfix") {
		t.Errorf("expected output to contain 'bugfix', got: %s", output)
	}
}

func TestVerifyApiSpec_FullDriftReport(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	// Create complex drift across all sources
	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/labstack/echo/v4"
func routes(e *echo.Echo) {
	e.GET("/code-only", nil)
	e.POST("/shared", nil)
}
`,
		openapiContent: `paths:
  /openapi-only:
    get: {}
  /shared:
    post: {}
`,
		apiMdContent: `| Method | Path |
|--------|------|
| GET | /md-only |
| POST | /shared |
`,
		cliMdContent: `### md-cmd`,
	})
	defer cleanup()

	_ = os.WriteFile("cmd/server.go", []byte(`package main
import "github.com/spf13/cobra"
var cmd = &cobra.Command{ Use: "code-cmd" }
`), 0644)

	output := captureStdout(t, func() {
		result := VerifyApiSpec("feature")
		if result {
			t.Error("expected VerifyApiSpec to fail with full drift")
		}
	})

	// Assertions for truncated drift reporting
	if !strings.Contains(output, "❌ Drift detected (showing first):") {
		t.Errorf("expected output to contain '❌ Drift detected (showing first):', got: %s", output)
	}
	// It should show at least the first one, and then a count of the rest
	if !strings.Contains(output, "... and 7 more drifts.") {
		t.Errorf("expected output to contain '... and 7 more drifts.', got: %s", output)
	}
}

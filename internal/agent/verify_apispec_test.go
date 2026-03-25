package agent

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// setupVerifyApiSpecEnv creates a temp directory with the required
// filesystem layout for VerifyApiSpec and changes the working directory
// to it. The returned cleanup function restores the original cwd.
func setupVerifyApiSpecEnv(t *testing.T, opts apiSpecEnvOpts) func() {
	t.Helper()
	tmpDir := t.TempDir()

	for _, dir := range []string{
		"cmd", "docs/cli", "internal/handler", "internal/module",
	} {
		if err := os.MkdirAll(tmpDir+"/"+dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Write Go source with cobra Use field for CLI extraction
	if opts.cmdGoContent != "" {
		_ = os.WriteFile(tmpDir+"/cmd/main.go", []byte(opts.cmdGoContent), 0644)
	}

	// Write OpenAPI spec
	if opts.openapiContent != "" {
		_ = os.WriteFile(tmpDir+"/docs/openapi.yaml", []byte(opts.openapiContent), 0644)
	}

	// Write API markdown doc
	if opts.apiMdContent != "" {
		_ = os.WriteFile(tmpDir+"/docs/api.md", []byte(opts.apiMdContent), 0644)
	}

	// Write CLI markdown doc
	if opts.cliMdContent != "" {
		_ = os.WriteFile(tmpDir+"/docs/cli.md", []byte(opts.cliMdContent), 0644)
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	return func() { _ = os.Chdir(origWd) }
}

type apiSpecEnvOpts struct {
	cmdGoContent   string
	openapiContent string
	apiMdContent   string
	cliMdContent   string
}

// captureStdout captures stdout output during the execution of fn.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// mockExecCommandApiSpec returns exit 0 with empty stdout.
func mockExecCommandApiSpec(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcessApiSpec", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcessApiSpec is the helper process for api spec mocks.
func TestHelperProcessApiSpec(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for _, a := range args {
		if a == "swagger-cli" {
			data, err := os.ReadFile("docs/openapi.yaml")
			if err == nil {
				fmt.Print(string(data))
			} else {
				fmt.Println("paths: {}")
			}
			os.Exit(0)
		}
	}
	os.Exit(0)
}

// mockExecCommandApiSpecFail returns exit 1 to simulate command failure.
func mockExecCommandApiSpecFail(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcessApiSpecFail", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcessApiSpecFail is the helper process that always fails.
func TestHelperProcessApiSpecFail(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintln(os.Stderr, "simulated command failure")
	os.Exit(1)
}

// mockExecCommandApiSpecWithRoutes returns a realistic OpenAPI YAML with
// path entries matching what is passed in the env var. By default returns
// a minimal YAML with paths that match the test fixtures.
func mockExecCommandApiSpecWithRoutes(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcessApiSpecWithRoutes", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcessApiSpecWithRoutes outputs valid OpenAPI YAML.
func TestHelperProcessApiSpecWithRoutes(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for _, a := range args {
		if a == "swagger-cli" {
			fmt.Print(`paths:
  /api/listings:
    get:
      summary: List
    post:
      summary: Create
`)
			os.Exit(0)
		}
	}
	os.Exit(0)
}

func TestVerifyApiSpec_Success(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpecWithRoutes
	defer func() { ExecCommand = orig }()

	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/labstack/echo/v4"
func routes(e *echo.Echo) {
	e.GET("/api/listings", nil)
	e.POST("/api/listings", nil)
}
`,
		openapiContent: `paths:
  /api/listings:
    get:
      summary: List
    post:
      summary: Create
`,
		apiMdContent: `| Method | Path | Description |
|--------|------|-------------|
| GET | /api/listings | List |
| POST | /api/listings | Create |
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
		if !result {
			t.Error("expected VerifyApiSpec to return true (no drift)")
		}
	})

	if !strings.Contains(output, "Gate PASS") {
		t.Errorf("expected output to contain 'Gate PASS', got: %s", output)
	}
}

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
		// apiMdContent intentionally empty — don't create the file
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

	// Assertions for specific drift messages
	expectedDrifts := []string{
		"Missing in OpenAPI (docs/openapi.yaml): GET /code-only",
		"❌ Missing in Code (cmd/server.go): GET /openapi-only",
		"Missing in API Docs (docs/api.md): GET /code-only",
		"❌ Missing in Code (cmd/server.go): GET /md-only",
		"❌ Missing in CLI Docs: code-cmd (found in Code)",
		"❌ Missing in Code: md-cmd (found in CLI Docs)",
	}

	for _, exp := range expectedDrifts {
		if !strings.Contains(output, exp) {
			t.Errorf("missing expected drift message: %s\nFull output:\n%s", exp, output)
		}
	}
}

func TestVerifyApiSpec_EmptyFiles(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/labstack/echo/v4"
func routes(e *echo.Echo) {}
`,
		openapiContent: `paths: {}`,
		apiMdContent:   "# fake", // non-empty but no routes
		cliMdContent:   "",       // empty CLI MD
	})
	defer cleanup()

	_ = os.WriteFile("cmd/server.go", []byte(`package main`), 0644)

	output := captureStdout(t, func() {
		result := VerifyApiSpec("feature")
		// Should pass because empty files = no routes = no drift
		if !result {
			t.Error("expected VerifyApiSpec to pass with empty documentation files")
		}
	})

	if !strings.Contains(output, "Gate PASS") {
		t.Errorf("expected Gate PASS, got: %s", output)
	}
}

func TestVerifyApiSpec_MarkdownPathEdgeCases(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommandApiSpec
	defer func() { ExecCommand = orig }()

	// Test backticks in markdown and trailing slashes
	// We use a unique route /api/v1/edge to avoid any overlap
	cleanup := setupVerifyApiSpecEnv(t, apiSpecEnvOpts{
		cmdGoContent: `package main
import "github.com/labstack/echo/v4"
func routes(e *echo.Echo) {
	e.GET("/api/v1/edge", nil)
}
`,
		openapiContent: `paths:
  /api/v1/edge:
    get:
      summary: edge
`,
		apiMdContent: `| Method | Path |
|--------|------|
| GET | /api/v1/edge |
`,
		cliMdContent: `### serve`,
	})
	defer cleanup()

	_ = os.WriteFile("cmd/server.go", []byte(`package main
import "github.com/spf13/cobra"
var cmd = &cobra.Command{ Use: "serve" }
`), 0644)

	output := captureStdout(t, func() {
		result := VerifyApiSpec("feature")
		if !result {
			t.Error("expected VerifyApiSpec to pass with normalized markdown paths")
		}
	})

	if !strings.Contains(output, "Gate PASS") {
		t.Errorf("expected Gate PASS, got: %s", output)
	}
}

func TestVerifyApiSpec_ExtractMarkdownRoutes_Direct(t *testing.T) {
	// Directly test ExtractMarkdownRoutes with various formats
	content := []byte(`
| GET | /users | List users |
| POST | ` + "`/users`" + ` | Create user |
| PUT | /users/ | Update user |
`)
	routes, err := ExtractMarkdownRoutes(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(routes) != 3 { // GET /users, POST /users, PUT /users
		t.Errorf("expected 3 unique routes, got %d", len(routes))
	}
	
	foundGet := false
	foundPost := false
	foundPut := false
	for _, r := range routes {
		if r.Method == "GET" && r.Path == "/users" { foundGet = true }
		if r.Method == "POST" && r.Path == "/users" { foundPost = true }
		if r.Method == "PUT" && r.Path == "/users" { foundPut = true }
	}
	if !foundGet || !foundPost || !foundPut {
		t.Errorf("missing expected routes: GET (/users): %v, POST (/users): %v, PUT (/users): %v", foundGet, foundPost, foundPut)
	}
}

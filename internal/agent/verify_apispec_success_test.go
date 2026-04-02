package agent

import (
	"os"
	"strings"
	"testing"
)

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
			t.Error("expected VerifyApiSpec to return true (no drift)")
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
		if r.Method == "GET" && r.Path == "/users" {
			foundGet = true
		}
		if r.Method == "POST" && r.Path == "/users" {
			foundPost = true
		}
		if r.Method == "PUT" && r.Path == "/users" {
			foundPut = true
		}
	}
	if !foundGet || !foundPost || !foundPut {
		t.Errorf("missing expected routes: GET (/users): %v, POST (/users): %v, PUT (/users): %v", foundGet, foundPost, foundPut)
	}
}

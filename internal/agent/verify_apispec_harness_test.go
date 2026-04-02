package agent

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
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

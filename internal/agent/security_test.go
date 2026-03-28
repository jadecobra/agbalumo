package agent

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSQLi(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Safe query with placeholder",
			code: `package main
func main() {
	db.Query("SELECT * FROM users WHERE id = ?", id)
}`,
			expected: 0,
		},
		{
			name: "Unsafe query with concatenation",
			code: `package main
func main() {
	db.Query("SELECT * FROM users WHERE id = " + id)
}`,
			expected: 1,
		},
		{
			name: "Unsafe query with fmt.Sprintf",
			code: `package main
import "fmt"
func main() {
	db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %s", id))
}`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}
			violations := checkSQLi(f, fset)
			if len(violations) != tt.expected {
				t.Errorf("Expected %d violations, got %d", tt.expected, len(violations))
			}
		})
	}
}

func TestCheckXSS(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Safe template.HTML with literal",
			code: `package main
import "html/template"
func main() {
	_ = template.HTML("<div>Static</div>")
}`,
			expected: 0,
		},
		{
			name: "Unsafe template.HTML with variable",
			code: `package main
import "html/template"
func main() {
	_ = template.HTML(userInput)
}`,
			expected: 1,
		},
		{
			name: "Unsafe Echo c.HTML with raw string",
			code: `package main
func Handler(c echo.Context) error {
	return c.HTML(200, "<div>" + userInput + "</div>")
}`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}
			violations := checkXSS(f, fset)
			if len(violations) != tt.expected {
				t.Errorf("Expected %d violations, got %d", tt.expected, len(violations))
			}
		})
	}
}

func TestCheckSecrets(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "AWS Access Key",
			content:  `aws_key := "AKIA1234567890ABCDEF"`, // gitleaks:allow
			expected: 1,
		},
		{
			name:     "Slack Webhook",
			content:  `url := "https://hooks.slack.com/services/T123/B456/789"`, // gitleaks:ignore
			expected: 1,
		},
		{
			name:     "High Entropy Secret",
			content:  `secret := "r5n7v9xc2mq4p6z8k1j3h5g7f9d1s3a2z4x6c8v0A1B2C3D4E5F6G7H8I9J0K1L2M3"`, // gitleaks:allow // 64 random alphanumeric chars
			expected: 1,
		},
		{
			name:     "Ignored secret",
			content:  `key := "AKIA1234567890ABCDEF" // #nosec - test key`, // gitleaks:allow
			expected: 0,
		},
		{
			name:     "Safe string",
			content:  `msg := "This is a normal string with no secrets"`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.go")
			code := "package main\nfunc main() {\n" + tt.content + "\n}"
			if err := os.WriteFile(tmpFile, []byte(code), 0600); err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			violations, err := VerifySecurityStatic(tmpDir)
			if err != nil {
				t.Fatalf("VerifySecurityStatic failed: %v", err)
			}

			count := 0
			for _, v := range violations {
				if v.Type == "Secret" || v.Type == "Entropy" {
					count++
				}
			}

			if count != tt.expected {
				t.Errorf("Expected %d secret violations, got %d", tt.expected, count)
			}
		})
	}
}

func TestVerifySecurityStatic_NonGoFiles(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create a .env file with a secret
	envFile := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envFile, []byte("API_KEY=AKIA1234567890ABCDEF"), 0600); err != nil { // gitleaks:ignore
		t.Fatalf("Failed to write .env: %v", err)
	}

	// Create a YAML file with a secret
	yamlFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(yamlFile, []byte("password: \"r5n7v9xc2mq4p6z8k1j3h5g7f9d1s3a2z4x6c8v0\""), 0600); err != nil { // gitleaks:ignore
		t.Fatalf("Failed to write yaml: %v", err)
	}

	violations, err := VerifySecurityStatic(tmpDir)
	if err != nil {
		t.Fatalf("VerifySecurityStatic failed: %v", err)
	}

	// We expect at least one violation for the .env and one for config.yaml
	foundEnv := false
	foundYaml := false
	for _, v := range violations {
		if strings.HasSuffix(v.File, ".env") {
			foundEnv = true
		}
		if strings.HasSuffix(v.File, "config.yaml") {
			foundYaml = true
		}
	}

	if !foundEnv {
		t.Error("Expected to find secret in .env file")
	}
	if !foundYaml {
		t.Error("Expected to find secret in config.yaml file")
	}
}

func TestCheckStructuralRaw(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filename string
		wantType string
	}{
		{
			name:     "Insecure onclick",
			content:  `<button onclick="alert('hi')">Click me</button>`,
			filename: "test.html",
			wantType: "Structural",
		},
		{
			name:     "Forbidden CDN",
			content:  `<script src="https://unpkg.com/htmx.org@1.9.10"></script>`,
			filename: "test.html",
			wantType: "Structural",
		},
		{
			name:     "Dangerous JS eval",
			content:  `eval("window.x = 1");`,
			filename: "test.js",
			wantType: "Structural",
		},
		{
			name:     "Inline Script",
			content:  `<script>console.log("inline");</script>`,
			filename: "test.html",
			wantType: "Structural",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(path, []byte(tt.content), 0600)
			assert.NoError(t, err)

			violations, err := checkStructuralRaw(path)
			assert.NoError(t, err)
			assert.NotEmpty(t, violations, "Should find violations in %s", tt.name)
			assert.Equal(t, tt.wantType, violations[0].Type)
		})
	}
}
func TestCheckFile_UsesSafeOpen(t *testing.T) {
	called := false
	orig := internalOpen
	internalOpen = func(path string) (*os.File, error) {
		called = true
		return os.Open(path)
	}
	defer func() { internalOpen = orig }()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.go")
	_ = os.WriteFile(path, []byte("package main"), 0600)

	_, _ = checkFile(path)
	assert.True(t, called, "Expected checkFile to use internalOpen hook for Go files")
}

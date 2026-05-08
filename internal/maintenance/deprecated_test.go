package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckDeprecatedPatterns(t *testing.T) {
	tmpDir := t.TempDir()

	setupDeprecatedFixtures(t, tmpDir)

	violations, err := CheckDeprecatedPatterns(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expect 3 violations: 
	// 1. map[string]interface{} in internal/module/bad.go
	// 2. map[string]interface{} in internal/handler/bad.go
	// 3. RenderWithBaseContext in internal/service/bad.go
	expected := 3
	if len(violations) != expected {
		t.Errorf("expected %d violations, got %d", expected, len(violations))
	}

	validateDeprecatedViolations(t, violations)
}

func setupDeprecatedFixtures(t *testing.T, tmpDir string) {
	// internal/module/bad.go
	moduleBad := filepath.Join(tmpDir, "internal", "module", "bad.go")
	_ = os.MkdirAll(filepath.Dir(moduleBad), 0750)
	_ = os.WriteFile(moduleBad, []byte(`package module
func Handle() {
	data := map[string]interface{}{}
	_ = data
}`), 0600)

	// internal/handler/bad.go
	handlerBad := filepath.Join(tmpDir, "internal", "handler", "bad.go")
	_ = os.MkdirAll(filepath.Dir(handlerBad), 0750)
	_ = os.WriteFile(handlerBad, []byte(`package handler
func Handle() {
	data := map[string]interface{}{}
	_ = data
}`), 0600)

	// internal/service/bad.go (RenderWithBaseContext)
	serviceBad := filepath.Join(tmpDir, "internal", "service", "bad.go")
	_ = os.MkdirAll(filepath.Dir(serviceBad), 0750)
	_ = os.WriteFile(serviceBad, []byte(`package service
func Render() {
	RenderWithBaseContext(w, r, "tmpl", data)
}`), 0600)

	// internal/module/good.go
	moduleGood := filepath.Join(tmpDir, "internal", "module", "good.go")
	_ = os.WriteFile(moduleGood, []byte(`package module
type MyViewModel struct{}
func Handle() {
	data := MyViewModel{}
	_ = data
}`), 0600)

	// internal/module/bad_test.go (should be ignored)
	moduleTest := filepath.Join(tmpDir, "internal", "module", "bad_test.go")
	_ = os.WriteFile(moduleTest, []byte(`package module
import "testing"
func TestHandle(t *testing.T) {
	data := map[string]interface{}{}
	_ = data
}`), 0600)
}

func validateDeprecatedViolations(t *testing.T, violations []DeprecatedViolation) {
	checkMapViolations(t, violations)
	checkRenderViolations(t, violations)
}

func checkMapViolations(t *testing.T, violations []DeprecatedViolation) {
	hasMapModule := false
	hasMapHandler := false

	for _, v := range violations {
		if v.Pattern == "map[string]interface{}" {
			updateFoundFlags(v.File, &hasMapModule, &hasMapHandler)
		}
	}

	if !hasMapModule {
		t.Error("expected violation for map[string]interface{} in internal/module")
	}
	if !hasMapHandler {
		t.Error("expected violation for map[string]interface{} in internal/handler")
	}
}

func updateFoundFlags(file string, module, handler *bool) {
	if strings.Contains(file, "internal/module/bad.go") {
		*module = true
	}
	if strings.Contains(file, "internal/handler/bad.go") {
		*handler = true
	}
}

func checkRenderViolations(t *testing.T, violations []DeprecatedViolation) {
	hasRender := false
	for _, v := range violations {
		if v.Pattern == "RenderWithBaseContext" {
			hasRender = true
			break
		}
	}
	if !hasRender {
		t.Error("expected violation for RenderWithBaseContext")
	}
}

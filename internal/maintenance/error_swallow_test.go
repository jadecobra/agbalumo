package maintenance_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/maintenance"
)

func setupTestModule(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	goMod := "module testpkg\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0600); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	for relPath, content := range files {
		fullPath := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
			t.Fatalf("failed to create dir for %s: %v", relPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0600); err != nil {
			t.Fatalf("failed to write file %s: %v", relPath, err)
		}
	}
	return dir
}

func TestCheckErrorSwallow(t *testing.T) {
	t.Run("MultiReturnErrorSwallowed", testMultiReturnErrorSwallowed)
	t.Run("SingleReturnErrorSwallowed", testSingleReturnErrorSwallowed)
	t.Run("CommaOkAndTypeAssertionNoViolation", testCommaOkAndTypeAssertionNoViolation)
	t.Run("DeferredCallsNoViolation", testDeferredCallsNoViolation)
	t.Run("ErrorProperlyAssignedNoViolation", testErrorProperlyAssignedNoViolation)
	t.Run("DiscardedErrorWithNolintCommentNoViolation", testDiscardedErrorWithNolintCommentNoViolation)
	t.Run("DefaultTargetDirsGracefulNonExistent", testDefaultTargetDirsGracefulNonExistent)
	t.Run("ExemptMethodCallsNoViolation", testExemptMethodCallsNoViolation)
}

func testMultiReturnErrorSwallowed(t *testing.T) {
	src := `package testpkg

import "errors"

func multi() (string, error) {
	return "ok", errors.New("fail")
}

func Caller() {
	val, _ := multi()
	_ = val
}
`
	dir := setupTestModule(t, map[string]string{"pkg/sample.go": src})
	violations, err := maintenance.CheckErrorSwallow(dir, "pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %+v", len(violations), violations)
	}
	if violations[0].Line != 10 {
		t.Errorf("expected violation at line 10, got %d", violations[0].Line)
	}
}

func testSingleReturnErrorSwallowed(t *testing.T) {
	src := `package testpkg

import "errors"

func single() error {
	return errors.New("fail")
}

func Caller() {
	_ = single()
}
`
	dir := setupTestModule(t, map[string]string{"pkg/sample.go": src})
	violations, err := maintenance.CheckErrorSwallow(dir, "pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %+v", len(violations), violations)
	}
	if violations[0].Line != 10 {
		t.Errorf("expected violation at line 10, got %d", violations[0].Line)
	}
}

func testCommaOkAndTypeAssertionNoViolation(t *testing.T) {
	src := `package testpkg

func CommaOk() {
	m := map[string]int{"k": 1}
	val, _ := m["k"]
	_ = val

	var x any = 100
	num, _ := x.(int)
	_ = num
}
`
	dir := setupTestModule(t, map[string]string{"pkg/sample.go": src})
	violations, err := maintenance.CheckErrorSwallow(dir, "pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func testDeferredCallsNoViolation(t *testing.T) {
	src := `package testpkg

type Closer interface {
	Close() error
}

type Tx interface {
	Rollback() error
}

func DeferFunc(r Closer, tx Tx) {
	defer func() {
		_ = r.Close()
	}()
	defer func() {
		_ = tx.Rollback()
	}()
}
`
	dir := setupTestModule(t, map[string]string{"pkg/sample.go": src})
	violations, err := maintenance.CheckErrorSwallow(dir, "pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func testErrorProperlyAssignedNoViolation(t *testing.T) {
	src := `package testpkg

import "errors"

func multi() (string, error) {
	return "ok", errors.New("fail")
}

func Caller() error {
	val, err := multi()
	if err != nil {
		return err
	}
	_ = val
	return nil
}
`
	dir := setupTestModule(t, map[string]string{"pkg/sample.go": src})
	violations, err := maintenance.CheckErrorSwallow(dir, "pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func testDiscardedErrorWithNolintCommentNoViolation(t *testing.T) {
	src := `package testpkg

import "errors"

func single() error {
	return errors.New("fail")
}

func multi() (string, error) {
	return "ok", errors.New("fail")
}

func Caller() {
	_ = single() // nolint:error-swallow intentional
	val, _ := multi() // explicit-ignore: known reason
	_ = val
}
`
	dir := setupTestModule(t, map[string]string{"pkg/sample.go": src})
	violations, err := maintenance.CheckErrorSwallow(dir, "pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func testDefaultTargetDirsGracefulNonExistent(t *testing.T) {
	srcRepo := `package repository

import "errors"

func fail() error {
	return errors.New("fail")
}

func Work() {
	_ = fail()
}
`
	dir := setupTestModule(t, map[string]string{
		"internal/repository/repo.go": srcRepo,
	})
	violations, err := maintenance.CheckErrorSwallow(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %+v", len(violations), violations)
	}
}

func testExemptMethodCallsNoViolation(t *testing.T) {
	src := `package testpkg

type Closer interface {
	Close() error
	Rollback() error
	CloseWithError(error) error
}

func Cleanup(c Closer) {
	_ = c.Close()
	_ = c.Rollback()
	_ = c.CloseWithError(nil)
}
`
	dir := setupTestModule(t, map[string]string{"pkg/sample.go": src})
	violations, err := maintenance.CheckErrorSwallow(dir, "pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

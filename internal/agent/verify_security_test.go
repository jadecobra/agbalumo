package agent

import (
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/util"
)

func TestVerifySecurityStaticGate(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("ViolationDetected_Query", func(t *testing.T) {
		content := `package main
func main() {
	db.Query("SELECT * FROM users WHERE id = " + id)
}`
		err := util.SafeWriteFile(filepath.Join(tmpDir, "unsafe_query.go"), []byte(content))
		if err != nil {
			t.Fatalf("failed to write unsafe file: %v", err)
		}

		if VerifySecurityStaticGate(tmpDir) {
			t.Error("VerifySecurityStaticGate should have failed due to SQLi violation in Query")
		}
	})

	t.Run("ViolationDetected_Prepare", func(t *testing.T) {
		content := `package main
func main() {
	db.Prepare("SELECT * FROM users WHERE id = " + id)
}`
		prepareDir := t.TempDir()
		err := util.SafeWriteFile(filepath.Join(prepareDir, "unsafe_prepare.go"), []byte(content))
		if err != nil {
			t.Fatalf("failed to write unsafe file: %v", err)
		}

		// This should currently PASS (fail to detect) because Prepare is not in the list
		if VerifySecurityStaticGate(prepareDir) {
			t.Error("VerifySecurityStaticGate should have failed due to SQLi violation in Prepare")
		}
	})

	t.Run("NoViolation", func(t *testing.T) {
		content := `package main
func main() {
	db.Query("SELECT * FROM users WHERE id = ?", id)
}`
		safeDir := t.TempDir()
		err := util.SafeWriteFile(filepath.Join(safeDir, "safe.go"), []byte(content))
		if err != nil {
			t.Fatalf("failed to write safe file: %v", err)
		}

		if !VerifySecurityStaticGate(safeDir) {
			t.Error("VerifySecurityStaticGate should have passed for safe code")
		}
	})
}

package agent

import (
	"testing"
)

func TestSecurityStaticIncremental(t *testing.T) {
	// Re-applying implementation: VerifySecurityStatic now supports multiple targets.
	_, err := VerifySecurityStatic("internal/agent/security.go", "internal/util/fs.go")
	if err != nil {
		t.Fatalf("Performance optimization failure: %v", err)
	}
}

package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfileComplexity_ThresholdAndTop(t *testing.T) {
	tmpDir := t.TempDir()

	simpleCode := `package sample
func Simple() int {
	return 42
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "simple.go"), []byte(simpleCode), 0o600); err != nil {
		t.Fatalf("failed to write simple.go: %v", err)
	}

	complexCode := `package sample
func HighComplexity(a, b, c bool, n int) int {
	res := 0
	if a {
		if b {
			res += 1
		} else if c {
			res += 2
		}
	}
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			res += i
		}
	}
	return res
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "complex.go"), []byte(complexCode), 0o600); err != nil {
		t.Fatalf("failed to write complex.go: %v", err)
	}

	opts := ComplexityProfileOptions{
		Threshold:   5,
		Top:         10,
		Diagnostics: true,
	}

	stats, err := ProfileComplexity(tmpDir, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("expected 1 function with complexity >= 5, got %d", len(stats))
	}

	if stats[0].FuncName != "HighComplexity" {
		t.Errorf("expected HighComplexity, got %s", stats[0].FuncName)
	}

	if len(stats[0].Diagnostics) == 0 {
		t.Errorf("expected diagnostics when enabled, got 0")
	}
}

func TestProfileComplexity_TopSorting(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package sample
func FnOne(a, b bool) {
	if a {
		if b {}
	}
}
func FnTwo(a, b, c bool) {
	if a {
		if b {
			if c {}
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(code), 0o600); err != nil {
		t.Fatalf("failed to write test.go: %v", err)
	}

	opts := ComplexityProfileOptions{
		Threshold: 1,
		Top:       1,
	}

	stats, err := ProfileComplexity(tmpDir, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("expected top 1 result, got %d", len(stats))
	}

	if stats[0].FuncName != "FnTwo" {
		t.Errorf("expected highest complexity function FnTwo, got %s", stats[0].FuncName)
	}
}

func TestProfileComplexity_SkipSpecialDirs(t *testing.T) {
	tmpDir := t.TempDir()
	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0o750); err != nil {
		t.Fatalf("failed to create vendor dir: %v", err)
	}

	code := `package vendor
func Ignored(a, b, c bool) {
	if a {
		if b {
			if c {}
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(vendorDir, "vendor.go"), []byte(code), 0o600); err != nil {
		t.Fatalf("failed to write vendor file: %v", err)
	}

	opts := ComplexityProfileOptions{
		Threshold: 1,
	}

	stats, err := ProfileComplexity(tmpDir, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("expected vendor files to be skipped, got %d stats", len(stats))
	}
}

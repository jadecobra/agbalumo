package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCalculateContextCost(t *testing.T) {
	// Create a temporary directory structure mimicking a small codebase
	tmpDir := t.TempDir()

	// 1. Create agent-readable files
	//   file1.go: 5 non-empty lines (actually just 5 lines)
	file1Content := "package main\n\nimport \"fmt\"\n\nfunc main() {}\n"
	err := os.WriteFile(filepath.Join(tmpDir, "file1.go"), []byte(file1Content), 0644)
	if err != nil {
		t.Fatalf("Failed to write file1.go: %v", err)
	}

	//   file2.md: 10 lines
	file2Content := ""
	for i := 0; i < 9; i++ {
		file2Content += "line\n"
	}
	file2Content += "line"
	err = os.WriteFile(filepath.Join(tmpDir, "file2.md"), []byte(file2Content), 0644)
	if err != nil {
		t.Fatalf("Failed to write file2.md: %v", err)
	}

	// 2. Create an ignored directory and file
	vendorDir := filepath.Join(tmpDir, "vendor")
	err = os.MkdirAll(vendorDir, 0755)
	if err != nil {
		t.Fatalf("Failed to mkdir vendor: %v", err)
	}
	err = os.WriteFile(filepath.Join(vendorDir, "file3.go"), []byte("package vendor\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write vendor/file3.go: %v", err)
	}

	// 3. Create an ignored file extension
	err = os.WriteFile(filepath.Join(tmpDir, "image.png"), []byte("fake png binary"), 0644)
	if err != nil {
		t.Fatalf("Failed to write image.png: %v", err)
	}

	// 4. Create ignored programmatically generated lockfiles
	err = os.WriteFile(filepath.Join(tmpDir, "package-lock.json"), []byte("{\"lockfileVersion\": 3}\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write package-lock.json: %v", err)
	}
	err = os.WriteFile(filepath.Join(tmpDir, "go.sum"), []byte("github.com/foo/bar v1.0.0 h1:hash\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.sum: %v", err)
	}

	// Calculate Cost
	report, err := CalculateContextCost(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if report == nil {
		t.Fatalf("Expected a report, got nil")
	}

	// Validate Total Files Processed: should only be file1.go and file2.md
	if report.TotalFiles != 2 {
		t.Errorf("Expected 2 files processed, got %d", report.TotalFiles)
	}

	// Validate RMS
	// file1.go LOC = 5 (\n count, maybe 4 or 5 depending on logic. Let's assume standard newline counting)
	// file2.md LOC = 10
	// For file1.go: 5 lines. For file2.md: 10 lines.
	// sum = 15. avg = 7.5.
	// sq_sum = 25 + 100 = 125.
	// mean(sq_sum) = 62.5
	// rms = sqrt(62.5) ≈ 7.90569

	// Just make sure it processes the lines somewhat correctly. We'll be precise in implementation.
	if report.TotalLines != 15 {
		t.Errorf("Expected TotalLines 15, got %d", report.TotalLines)
	}

	if report.RMS < 7.9 || report.RMS > 8.0 {
		t.Errorf("Expected RMS ~7.905, got %f", report.RMS)
	}

	// Validate top expensive files
	if len(report.TopFiles) != 2 {
		t.Errorf("Expected 2 TopFiles, got %d", len(report.TopFiles))
	} else {
		// TopFiles are now sorted by Tokens descending
		if report.TopFiles[0].Tokens <= 0 {
			t.Errorf("Expected top file to have positive token count, got %d", report.TopFiles[0].Tokens)
		}
	}
	// Validate token counts are populated for each file
	if report.TotalTokens <= 0 {
		t.Errorf("Expected TotalTokens > 0, got %d", report.TotalTokens)
	}
	for _, fc := range report.TopFiles {
		if fc.Tokens <= 0 {
			t.Errorf("Expected Tokens > 0 for file %s, got %d", fc.FilePath, fc.Tokens)
		}
	}

	// Validate TokenRMS is positive
	if report.TokenRMS <= 0 {
		t.Errorf("Expected TokenRMS > 0, got %f", report.TokenRMS)
	}

	// Validate ContextWindowPct is between 0 and 100 (tiny test dir is well under 200k tokens)
	if report.ContextWindowPct <= 0 {
		t.Errorf("Expected ContextWindowPct > 0, got %f", report.ContextWindowPct)
	}
	if report.ContextWindowPct >= 100 {
		t.Errorf("Expected ContextWindowPct < 100 for tiny test dir, got %f", report.ContextWindowPct)
	}
}

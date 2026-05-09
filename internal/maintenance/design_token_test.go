package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckTokenStrictness(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		shouldFail bool
	}{
		{name: "arbitrary text color hex", content: `<div class="text-[#123]">`, shouldFail: true},
		{name: "arbitrary padding px", content: `<div class="p-[17px]">`, shouldFail: true},
		{name: "arbitrary width percent", content: `<div class="w-[32.5%]">`, shouldFail: true},
		{name: "arbitrary margin rem", content: `<div class="m-[2.5rem]">`, shouldFail: true},
		{name: "arbitrary height vh", content: `<div class="h-[100vh]">`, shouldFail: true},
		{name: "arbitrary border color", content: `<div class="border-[#fff]">`, shouldFail: true},
		{name: "standard text color", content: `<div class="text-surface-dark">`, shouldFail: false},
		{name: "standard padding", content: `<div class="p-4">`, shouldFail: false},
		{name: "standard width", content: `<div class="w-full">`, shouldFail: false},
		{name: "non-guarded arbitrary value (top)", content: `<div class="top-[20px]">`, shouldFail: false},
		{name: "arbitrary background", content: `<div class="bg-[#000]">`, shouldFail: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTokenTest(t, tt.content, tt.shouldFail)
		})
	}
}

func runTokenTest(t *testing.T, content string, shouldFail bool) {
	t.Helper()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.html")
	err := os.WriteFile(filePath, []byte(content), 0600)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	violations, err := CheckTokenStrictness(tmpDir)
	if err != nil {
		t.Fatalf("CheckTokenStrictness failed: %v", err)
	}

	if shouldFail && len(violations) == 0 {
		t.Errorf("expected violation for content: %s", content)
	}
	if !shouldFail && len(violations) > 0 {
		t.Errorf("expected no violation for content: %s, got %d issues", content, len(violations))
	}
}

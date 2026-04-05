package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckGosecRationale(t *testing.T) {
	// Setup temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "gosec-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tests := []struct {
		name     string
		content  string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid with hyphen",
			content:  "// #nosec G101 - this is a rationale",
			filename: "valid.go",
			wantErr:  false,
		},
		{
			name:     "valid with double hyphen",
			content:  "// #nosec G101 -- this is a rationale",
			filename: "valid2.go",
			wantErr:  false,
		},
		{
			name:     "invalid missing rationale",
			content:  "// #nosec G101",
			filename: "invalid.go",
			wantErr:  true,
		},
		{
			name:     "invalid with rule IDs but no rationale",
			content:  "// #nosec G101,G601",
			filename: "invalid2.go",
			wantErr:  true,
		},
		{
			name:     "no nosec directive",
			content:  "func main() {}",
			filename: "none.go",
			wantErr:  false,
		},
		{
			name:     "nosec without space",
			content:  "//#nosec G101",
			filename: "invalid3.go",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh temp directory for each test case
			caseDir, err := os.MkdirTemp(tmpDir, "case-*")
			if err != nil {
				t.Fatalf("failed to create case dir: %v", err)
			}
			// Don't remove caseDir here, as we need to write into it!
			// tmpDir cleanup will handle this via defer.

			filePath := filepath.Join(caseDir, tt.filename)
			if werr := os.WriteFile( /*nolint:gosec*/ filePath, []byte(tt.content), 0600); werr != nil {
				t.Fatalf("failed to write test file: %v", werr)
			}

			// We pass the caseDir as the root to CheckGosecRationale
			err = CheckGosecRationale(caseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckGosecRationale() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckSnapshotParity(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		wantViol  bool
		wantCount int
	}{
		{
			name:      "darwin only - fail",
			files:     []string{"test-darwin.png"},
			wantViol:  true,
			wantCount: 1,
		},
		{
			name:      "darwin and linux - pass",
			files:     []string{"test-darwin.png", "test-linux.png"},
			wantViol:  false,
			wantCount: 0,
		},
		{
			name:      "linux only - pass",
			files:     []string{"test-linux.png"},
			wantViol:  false,
			wantCount: 0,
		},
		{
			name:      "multiple darwin, some missing linux",
			files:     []string{"a-darwin.png", "a-linux.png", "b-darwin.png"},
			wantViol:  true,
			wantCount: 1,
		},
		{
			name:      "non-snapshot files ignored",
			files:     []string{"readme.md", "test.go"},
			wantViol:  false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			setupSnapshots(t, tmpDir, tt.files)

			violations, err := CheckSnapshotParity(tmpDir)
			if err != nil {
				t.Fatalf("CheckSnapshotParity failed: %v", err)
			}

			if (len(violations) > 0) != tt.wantViol || len(violations) != tt.wantCount {
				t.Errorf("got %d violations, want %d", len(violations), tt.wantCount)
			}
		})
	}
}

func setupSnapshots(t *testing.T, root string, files []string) {
	t.Helper()
	dir := filepath.Join(root, "tests/e2e/visual.spec.ts-snapshots")
	if err := os.MkdirAll(dir, 0750); err != nil { // G301 fix
		t.Fatalf("failed to create snapshots dir: %v", err)
	}
	for _, f := range files {
		path := filepath.Join(dir, f)
		if err := os.WriteFile(path, []byte("test"), 0600); err != nil { // G306 fix
			t.Fatalf("failed to write file %s: %v", f, err)
		}
	}
}

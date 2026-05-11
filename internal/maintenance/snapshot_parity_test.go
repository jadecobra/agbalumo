package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckSnapshotParity(t *testing.T) {
	tests := []struct { //nolint:govet
		wantCount int
		wantViol  bool
		name      string
		files     map[string]string
	}{
		{
			name:      "darwin only - fail",
			files:     map[string]string{"test-darwin.png": "content1"},
			wantViol:  true,
			wantCount: 1,
		},
		{
			name:      "darwin and linux different - pass",
			files:     map[string]string{"test-darwin.png": "content1", "test-linux.png": "content2"},
			wantViol:  false,
			wantCount: 0,
		},
		{
			name:      "darwin and linux identical - fail",
			files:     map[string]string{"test-darwin.png": "identical", "test-linux.png": "identical"},
			wantViol:  true,
			wantCount: 1,
		},
		{
			name:      "linux only - pass",
			files:     map[string]string{"test-linux.png": "content1"},
			wantViol:  false,
			wantCount: 0,
		},
		{
			name:      "multiple darwin, some missing linux",
			files:     map[string]string{"a-darwin.png": "c1", "a-linux.png": "c2", "b-darwin.png": "c3"},
			wantViol:  true,
			wantCount: 1,
		},
		{
			name:      "non-snapshot files ignored",
			files:     map[string]string{"readme.md": "c1", "test.go": "c2"},
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

func setupSnapshots(t *testing.T, root string, files map[string]string) {
	t.Helper()
	dir := filepath.Join(root, "tests/e2e/visual.spec.ts-snapshots")
	if err := os.MkdirAll(dir, 0750); err != nil {
		t.Fatalf("failed to create snapshots dir: %v", err)
	}
	for f, content := range files {
		path := filepath.Join(dir, f)
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatalf("failed to write file %s: %v", f, err)
		}
	}
}

package maintenance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type snapshotParityTestCase struct {
	name      string
	files     string // JSON map
	wantCount int
	wantViol  bool
}

//nolint:gocognit
func TestCheckSnapshotParity(t *testing.T) {
	tests := []snapshotParityTestCase{
		{
			name:      "darwin only - fail",
			files:     `{"test-darwin.png": "content1"}`,
			wantViol:  true,
			wantCount: 1,
		},
		{
			name:      "darwin and linux different - pass",
			files:     `{"test-darwin.png": "content1", "test-linux.png": "content2"}`,
			wantViol:  false,
			wantCount: 0,
		},
		{
			name:      "darwin and linux identical - fail",
			files:     `{"test-darwin.png": "identical", "test-linux.png": "identical"}`,
			wantViol:  true,
			wantCount: 1,
		},
		{
			name:      "linux only - pass",
			files:     `{"test-linux.png": "content1"}`,
			wantViol:  false,
			wantCount: 0,
		},
		{
			name:      "multiple darwin, some missing linux",
			files:     `{"a-darwin.png": "c1", "a-linux.png": "c2", "b-darwin.png": "c3"}`,
			wantViol:  true,
			wantCount: 1,
		},
		{
			name:      "non-snapshot files ignored",
			files:     `{"readme.md": "c1", "test.go": "c2"}`,
			wantViol:  false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			var fileMap map[string]string
			if err := json.Unmarshal([]byte(tt.files), &fileMap); err != nil {
				t.Fatalf("failed to unmarshal files: %v", err)
			}
			setupSnapshots(t, tmpDir, fileMap)

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

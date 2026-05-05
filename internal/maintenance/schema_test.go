package maintenance

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDumpSQLiteSchema(t *testing.T) {
	tests := []struct {
		setup   func(string) string
		name    string
		wantErr bool
	}{
		{
			name: "missing file returns error",
			setup: func(dir string) string {
				return filepath.Join(dir, "missing.db")
			},
			wantErr: true,
		},
		{
			name: "valid sqlite file returns schema",
			setup: func(dir string) string {
				dbPath := filepath.Join(dir, "test.db")
				// Initialize an empty sqlite db with a table
				// #nosec G204
				cmd := exec.Command("sqlite3", dbPath, "CREATE TABLE users (id INTEGER PRIMARY KEY);")
				if err := cmd.Run(); err != nil {
					t.Fatalf("failed to setup mock db: %v", err)
				}
				return dbPath
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		runSchemaTest(t, tt.name, tt.setup, tt.wantErr)
	}
}

func runSchemaTest(t *testing.T, name string, setup func(string) string, wantErr bool) {
	t.Run(name, func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := setup(tmpDir)

		got, err := DumpSQLiteSchema(dbPath)
		if (err != nil) != wantErr {
			t.Errorf("DumpSQLiteSchema() error = %v, wantErr %v", err, wantErr)
			return
		}
		if !wantErr && got == "" {
			t.Errorf("DumpSQLiteSchema() returned empty string for valid db")
		}
	})
}

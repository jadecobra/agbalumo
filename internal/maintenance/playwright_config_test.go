package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyPlaywrightConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid config with open: never",
			content: "reporter: [['html', { open: 'never' }]],",
			wantErr: false,
		},
		{
			name:    "invalid config with default html reporter",
			content: "reporter: 'html',",
			wantErr: true,
		},
		{
			name:    "invalid config with missing setting",
			content: "reporter: [['list']],",
			wantErr: true,
		},
		{
			name:    "valid config with complex structure",
			content: `export default defineConfig({
  reporter: [['html', { open: 'never' }]],
  use: { ... }
});`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "playwright.config.ts")
			if err := os.WriteFile(configPath, []byte(tt.content), 0600); err != nil {
				t.Fatalf("failed to write mock config: %v", err)
			}

			err := VerifyPlaywrightConfig(tmpDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPlaywrightConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

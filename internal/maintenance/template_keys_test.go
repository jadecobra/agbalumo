package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckTemplateKeyGaps(t *testing.T) {
	tests := []struct {
		name       string
		files      map[string]string
		shouldFail bool
	}{
		{
			name: "missing key violation",
			files: map[string]string{
				"card.html": `{{ define "card" }}<div>{{ $.SavedIDs }}</div>{{ end }}`,
				"home.html": `<div>{{ template "card" dict "Listing" . }}</div>`,
			},
			shouldFail: true,
		},
		{
			name: "all keys present",
			files: map[string]string{
				"card.html": `{{ define "card" }}<div>{{ $.SavedIDs }}</div>{{ end }}`,
				"home.html": `<div>{{ template "card" dict "SavedIDs" $.SavedIDs "Listing" . }}</div>`,
			},
			shouldFail: false,
		},
		{
			name: "multiple keys and dot references",
			files: map[string]string{
				"card.html": `{{ define "card" }}<div>{{ .Listing.Title }} - {{ $.User.Name }}</div>{{ end }}`,
				"home.html": `<div>{{ template "card" dict "Listing" . "User" $.User }}</div>`,
			},
			shouldFail: false,
		},
		{
			name: "missing User reference",
			files: map[string]string{
				"card.html": `{{ define "card" }}<div>{{ $.User.Name }}</div>{{ end }}`,
				"home.html": `<div>{{ template "card" dict "Listing" . }}</div>`,
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "template_test_*")
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = os.RemoveAll(tmpDir) }()

			for name, content := range tt.files {
				if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			violations, err := CheckTemplateKeyGaps(tmpDir)
			if err != nil {
				t.Fatalf("CheckTemplateKeyGaps failed: %v", err)
			}

			if tt.shouldFail && len(violations) == 0 {
				t.Errorf("expected violations, got none")
			}
			if !tt.shouldFail && len(violations) > 0 {
				t.Errorf("expected no violations, got %d: %+v", len(violations), violations)
			}
		})
	}
}

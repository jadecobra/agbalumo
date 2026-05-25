package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

type templateTest struct {
	files      map[string]string
	name       string
	shouldFail bool
}

func TestCheckTemplateKeyGaps(t *testing.T) {
	tests := []templateTest{
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
		// YAML Schema checks
		{
			name: "YAML schema missing required prop",
			files: map[string]string{
				"button.html": `{{ define "btn" }}
<!-- props:
  Label: required
  Type: optional
-->
<button>{{ .Label }}</button>
{{ end }}`,
				"home.html": `<div>{{ template "btn" dict "Type" "submit" }}</div>`,
			},
			shouldFail: true,
		},
		{
			name: "YAML schema extraneous prop",
			files: map[string]string{
				"button.html": `{{ define "btn" }}
<!-- props:
  Label: required
  Type: optional
-->
<button>{{ .Label }}</button>
{{ end }}`,
				"home.html": `<div>{{ template "btn" dict "Label" "Submit" "UnknownProp" "value" }}</div>`,
			},
			shouldFail: true,
		},
		{
			name: "YAML schema correct props",
			files: map[string]string{
				"button.html": `{{ define "btn" }}
<!-- props:
  Label: required
  Type: optional
-->
<button>{{ .Label }}</button>
{{ end }}`,
				"home.html": `<div>{{ template "btn" dict "Label" "Submit" "Type" "submit" }}</div>`,
			},
			shouldFail: false,
		},
		{
			name: "YAML schema double quoted value is not treated as extraneous key",
			files: map[string]string{
				"badge.html": `{{ define "badge" }}
<!-- props:
  Label: required
  ColorClasses: optional
-->
<span>{{ .Label }}</span>
{{ end }}`,
				"home.html": `<div>{{ template "badge" dict "Label" "custom" "ColorClasses" "bg-green-500" }}</div>`,
			},
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTemplateKeyGapTest(t, tt)
		})
	}
}

func runTemplateKeyGapTest(t *testing.T, tt templateTest) {
	tmpDir, err := os.MkdirTemp("", "template_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	for name, content := range tt.files {
		path := filepath.Join(tmpDir, name)
		err = os.WriteFile(path, []byte(content), 0600)
		if err != nil {
			t.Fatal(err)
		}
	}

	var violations []DesignViolation
	violations, err = CheckTemplateKeyGaps(tmpDir)
	if err != nil {
		t.Fatalf("CheckTemplateKeyGaps failed: %v", err)
	}

	if tt.shouldFail && len(violations) == 0 {
		t.Errorf("expected violations, got none")
	}
	if !tt.shouldFail && len(violations) > 0 {
		t.Errorf("expected no violations, got %d: %+v", len(violations), violations)
	}
}

func TestParseDictKeys(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "simple keys",
			content:  `"Label" "Submit" "Classes" "bg-earth-dark"`,
			expected: []string{"Label", "Classes"},
		},
		{
			name:     "keys with parenthesized expressions",
			content:  `"ID" (print "edit-listing-modal-" .Listing.ID) "Title" "Edit"`,
			expected: []string{"ID", "Title"},
		},
		{
			name:     "keys with mixed types",
			content:  `"ID" 123 "IsForm" true "Data" .`,
			expected: []string{"ID", "IsForm", "Data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runSingleParseDictKeysTest(t, tt.content, tt.expected)
		})
	}
}

func runSingleParseDictKeysTest(t *testing.T, content string, expected []string) {
	keys := parseDictKeys(content)
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d: %v", len(expected), len(keys), keys)
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected key at index %d to be %q, got %q", i, expected[i], k)
		}
	}
}

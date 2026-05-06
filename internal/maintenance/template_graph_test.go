package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildTemplateGraph(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "template_graph_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	setupTestTemplates(t, tmpDir)

	graph, err := BuildTemplateGraph(tmpDir)
	if err != nil {
		t.Fatalf("BuildTemplateGraph failed: %v", err)
	}

	testCardNode(t, graph)
	testWrapperNode(t, graph)
}

func testCardNode(t *testing.T, graph map[string]*TemplateNode) {
	t.Run("CardNode", func(t *testing.T) {
		card, ok := graph["card"]
		if !ok {
			t.Fatal("expected 'card' node in graph")
		}
		checkCardReferences(t, card)
		checkCardCallers(t, card)
	})
}

func testWrapperNode(t *testing.T, graph map[string]*TemplateNode) {
	t.Run("WrapperNode", func(t *testing.T) {
		wrapper, ok := graph["wrapper"]
		if !ok {
			t.Fatal("expected 'wrapper' node in graph")
		}
		foundCardCall := false
		for _, call := range wrapper.Calls {
			if call == "card" {
				foundCardCall = true
			}
		}
		if !foundCardCall {
			t.Errorf("expected 'wrapper' to call 'card', got %v", wrapper.Calls)
		}
	})
}

func setupTestTemplates(t *testing.T, tmpDir string) {
	var err error
	file1 := filepath.Join(tmpDir, "card.html")
	content1 := `{{ define "card" }}
	<div class="user">{{ $.User }}</div>
	<div class="meta">{{ .Meta }}</div>
{{ end }}`
	if err = os.WriteFile(file1, []byte(content1), 0600); err != nil {
		t.Fatal(err)
	}

	file2 := filepath.Join(tmpDir, "page.html")
	content2 := `{{ template "card" dict "Listing" .Listing }}`
	if err = os.WriteFile(file2, []byte(content2), 0600); err != nil {
		t.Fatal(err)
	}

	file3 := filepath.Join(tmpDir, "nested.html")
	content3 := `{{ define "wrapper" }}
	{{ template "card" dict "User" .User }}
{{ end }}`
	if err = os.WriteFile(file3, []byte(content3), 0600); err != nil {
		t.Fatal(err)
	}
}

func checkCardReferences(t *testing.T, card *TemplateNode) {
	hasUser := false
	hasMeta := false
	for _, ref := range card.References {
		if ref == "User" {
			hasUser = true
		}
		if ref == "Meta" {
			hasMeta = true
		}
	}
	if !hasUser || !hasMeta {
		t.Errorf("expected references 'User' and 'Meta', got %v", card.References)
	}
}

func checkCardCallers(t *testing.T, card *TemplateNode) {
	if len(card.CalledBy) != 2 {
		t.Errorf("expected 2 callers for 'card', got %d", len(card.CalledBy))
	}

	var pageCall *TemplateCall
	for _, call := range card.CalledBy {
		if filepath.Base(call.File) == "page.html" {
			pageCall = &call
			break
		}
	}

	if pageCall == nil {
		t.Fatal("expected call from page.html")
	}

	foundListing := false
	for _, k := range pageCall.DictKeys {
		if k == "Listing" {
			foundListing = true
		}
	}
	if !foundListing {
		t.Errorf("expected 'Listing' in DictKeys for page.html call, got %v", pageCall.DictKeys)
	}
}

func TestExtractNodeReferences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "IgnoreJSAndCSS",
			input:    `<div>{{ .Listing.Title }}</div><script>let x = document.getElementById();</script><style>margin: .5em;</style>`,
			expected: []string{"Listing", "Title"},
		},
		{
			name:     "MultipleActions",
			input:    `{{ .A }} {{ .B }}`,
			expected: []string{"A", "B"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractNodeReferences(tt.input)
			if strings.Join(got, ",") != strings.Join(tt.expected, ",") {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

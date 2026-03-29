package agent

import (
	"reflect"
	"testing"
)

func TestParseMarkdownTracker(t *testing.T) {
	content := `
# Category 1
Description 1
- [x] Step 1 (Completed)
- [ ] Step 2

# Category 2
Description 2
- [x] Step 3 (Completed)
- [x] Step 4 (Completed)
`

	expected := ProgressTracker{
		Features: []Feature{
			{
				Category:    "Category 1",
				Description: "Description 1",
				Passes:      false,
				Steps:       []string{"Step 1 (Completed)", "Step 2"},
			},
			{
				Category:    "Category 2",
				Description: "Description 2",
				Passes:      true,
				Steps:       []string{"Step 3 (Completed)", "Step 4 (Completed)"},
			},
		},
	}

	got, err := ParseMarkdownTracker(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %+v, want %+v", got, expected)
	}
}

func TestToMarkdown(t *testing.T) {
	tracker := ProgressTracker{
		Features: []Feature{
			{
				Category:    "Category 1",
				Description: "Description 1",
				Passes:      false,
				Steps:       []string{"Step 1 (Completed)", "Step 2"},
			},
		},
	}

	expected := `# Category 1
Description 1
- [x] Step 1
- [ ] Step 2
`
	got := ToMarkdown(tracker)
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

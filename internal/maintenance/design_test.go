package maintenance

import (
	"os"
	"testing"
)

func TestCheckMinFontSize(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		shouldFail bool
	}{
		{name: "font size 8px", line: `text-[8px]`, shouldFail: true},
		{name: "font size 10px", line: `text-[10px]`, shouldFail: false},
		{name: "font size xs", line: `text-xs`, shouldFail: false},
		{name: "font size 12px", line: `text-[12px]`, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkMinFontSize("test.html", 1, tt.line)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation for %s, got none", tt.line)
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation for %s, got %d", tt.line, len(v))
			}
		})
	}
}

func TestCheckLowContrastOpacity(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		shouldFail bool
	}{
		{name: "opacity 60", line: `text-text-sub/60`, shouldFail: true},
		{name: "opacity 80", line: `text-text-sub/80`, shouldFail: false},
		{name: "plain text-text-sub", line: `text-text-sub`, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkLowContrastOpacity("test.html", 1, tt.line)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation for %s, got none", tt.line)
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation for %s, got %d", tt.line, len(v))
			}
		})
	}
}

func TestCheckHardcodedModalBg(t *testing.T) {
	tests := []struct {
		name       string
		file       string
		line       string
		shouldFail bool
	}{
		{name: "bg-earth-dark in modal_detail.html", file: "modal_detail.html", line: `bg-earth-dark`, shouldFail: true},
		{name: "bg-earth-dark in ui_components.html", file: "ui_components.html", line: `bg-earth-dark`, shouldFail: true},
		{name: "dark:bg-earth-dark in modal_detail.html", file: "modal_detail.html", line: `dark:bg-earth-dark`, shouldFail: false},
		{name: "bg-earth-dark in other.html", file: "other.html", line: `bg-earth-dark`, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkHardcodedModalBg(tt.file, 1, tt.line)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation for %s in %s, got none", tt.line, tt.file)
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation for %s in %s, got %d", tt.line, tt.file, len(v))
			}
		})
	}
}

func TestCheckInlineHandlers(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		shouldFail bool
	}{
		{name: "onclick handler", line: `<div onclick="foo()">`, shouldFail: true},
		{name: "onchange handler", line: `<select onchange="bar()">`, shouldFail: true},
		{name: "onsubmit handler", line: `<form onsubmit="baz()">`, shouldFail: true},
		{name: "onmouseover handler", line: `<span onmouseover="qux()">`, shouldFail: true},
		{name: "hx-get attribute", line: `<div hx-get="/path">`, shouldFail: false},
		{name: "data-action attribute", line: `<div data-action="click">`, shouldFail: false},
		{name: "no handler", line: `<div>`, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkInlineHandlers("test.html", 1, tt.line)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation for %s, got none", tt.line)
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation for %s, got %d", tt.line, len(v))
			}
		})
	}
}

func TestCheckInlineStyles(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		shouldFail bool
	}{
		{name: "inline style color", line: `<div style="color:red">`, shouldFail: true},
		{name: "inline style display", line: `<div style="display:none">`, shouldFail: true},
		{name: "inline style background-image", line: `<div style="background-image: url('/img.png')">`, shouldFail: false},
		{name: "tailwind class", line: `<div class="text-red-500">`, shouldFail: false},
		{name: "no style", line: `<div>`, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkInlineStyles("test.html", 1, tt.line)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation for %s, got none", tt.line)
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation for %s, got %d", tt.line, len(v))
			}
		})
	}
}

func TestCheckUppercaseDensity(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		shouldFail bool
	}{
		{
			name: "under limit (3)",
			content: `
				<div class="uppercase">1</div>
				<div class="uppercase">2</div>
				<div class="uppercase">3</div>
			`,
			shouldFail: false,
		},
		{
			name: "at limit (4)",
			content: `
				<div class="uppercase">1</div>
				<div class="uppercase">2</div>
				<div class="uppercase">3</div>
				<div class="uppercase">4</div>
			`,
			shouldFail: false,
		},
		{
			name: "over limit (5)",
			content: `
				<div class="uppercase">1</div>
				<div class="uppercase">2</div>
				<div class="uppercase">3</div>
				<div class="uppercase">4</div>
				<div class="uppercase">5</div>
			`,
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runUppercaseDensityTest(t, tt.content, tt.shouldFail)
		})
	}
}

func runUppercaseDensityTest(t *testing.T, content string, shouldFail bool) {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test_uppercase_*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

	if _, err = tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	v, err := checkUppercaseDensity(tmpfile.Name())
	if err != nil {
		t.Fatalf("checkUppercaseDensity failed: %v", err)
	}

	if shouldFail && len(v) == 0 {
		t.Errorf("expected violation, got none")
	}
	if !shouldFail && len(v) > 0 {
		t.Errorf("expected no violation, got %d", len(v))
	}
}

func TestCheckMissingTestIDs(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		shouldFail bool
	}{
		{name: "button with hx-post missing data-testid", line: `<button hx-post="/save">`, shouldFail: true},
		{name: "anchor with hx-get missing data-testid", line: `<a hx-get="/details">`, shouldFail: true},
		{name: "button with hx-post and data-testid", line: `<button hx-post="/save" data-testid="save-btn">`, shouldFail: false},
		{name: "anchor with hx-get and data-testid", line: `<a hx-get="/details" data-testid="detail-link">`, shouldFail: false},
		{name: "button without hx-attribute", line: `<button class="plain">`, shouldFail: false},
		{name: "anchor without hx-attribute", line: `<a href="/foo">`, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkMissingTestIDs("test.html", 1, tt.line)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation for %s, got none", tt.line)
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation for %s, got %d", tt.line, len(v))
			}
		})
	}
}

func TestCheckDeadSpacers(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		shouldFail bool
	}{
		{name: "empty section", line: `<section></section>`, shouldFail: true},
		{name: "section with attributes but empty", line: `<section class="py-4"></section>`, shouldFail: true},
		{name: "section with content", line: `<section>content</section>`, shouldFail: false},
		{name: "section with leading space only", line: `<section>   </section>`, shouldFail: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkDeadSpacers("test.html", 1, tt.line)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation for %s, got none", tt.line)
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation for %s, got %d", tt.line, len(v))
			}
		})
	}
}
func TestA11ySemantics(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		hasLabel   bool
		shouldFail bool
	}{
		{name: "icon button missing aria-label", line: `<button class="icon-close"><svg...`, shouldFail: true},
		{name: "icon button with aria-label", line: `<button aria-label="Close" class="icon-close"><svg...`, shouldFail: false},
		{name: "label without for", line: `<label>Name</label>`, shouldFail: true},
		{name: "label with for", line: `<label for="name">Name</label>`, shouldFail: false},
		{name: "input missing id (with label)", line: `<input type="text">`, hasLabel: true, shouldFail: true},
		{name: "input with id (with label)", line: `<input id="name">`, hasLabel: true, shouldFail: false},
		{name: "input missing id (no label)", line: `<input type="text">`, hasLabel: false, shouldFail: false},
		{name: "img missing alt", line: `<img src="x.jpg">`, shouldFail: true},
		{name: "img with decorative alt", line: `<img src="x.jpg" alt="">`, shouldFail: false},
		{name: "img with alt text", line: `<img src="x.jpg" alt="A photo">`, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkA11ySemantics("test.html", 1, tt.line, tt.hasLabel)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation for %s, got none", tt.line)
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation for %s, got %d", tt.line, len(v))
			}
		})
	}
}

func TestCheckFileStandardsIgnoreComments(t *testing.T) {
	content := `
		<!-- props:
		  Label: required
		  rounded-md inside a comment shouldn't fail
		  #D4A373 inside a comment shouldn't fail
		-->
		<div class="bg-earth-dark text-earth-cream text-sm">
			<script>
				console.log("rounded-md inside a script");
				var color = "#D4A373";
			</script>
			<span>Hello</span>
		</div>
	`

	tmpfile, err := os.CreateTemp("", "test_comment_ignore_*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

	if _, err = tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	violations, err := checkFileStandards(tmpfile.Name())
	if err != nil {
		t.Fatalf("checkFileStandards failed: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf("expected 0 violations (rounded-md/hex bypassed in comment/script), got %d: %+v", len(violations), violations)
	}
}

func TestCheckHtmxIndicator(t *testing.T) {
	tests := []struct {
		attrs      map[string]string
		name       string
		shouldFail bool
	}{
		{name: "hx-post without indicator", attrs: map[string]string{"hx-post": "/submit"}, shouldFail: true},
		{name: "hx-delete without indicator", attrs: map[string]string{"hx-delete": "/item"}, shouldFail: true},
		{name: "hx-post with indicator", attrs: map[string]string{"hx-post": "/submit", "hx-indicator": "#spinner"}, shouldFail: false},
		{name: "hx-get without indicator", attrs: map[string]string{"hx-get": "/items"}, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkHtmxIndicator("test.html", 1, "test-element", tt.attrs)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation for %v, got none", tt.attrs)
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation for %v, got %d", tt.attrs, len(v))
			}
		})
	}
}

func TestCheckTextSelectionUsability(t *testing.T) {
	tests := []struct {
		attrs      map[string]string
		name       string
		shouldFail bool
	}{
		{name: "select-none present", attrs: map[string]string{"class": "select-none flex p-4"}, shouldFail: true},
		{name: "select-none absent", attrs: map[string]string{"class": "flex p-4 text-sm"}, shouldFail: false},
		{name: "no class attribute", attrs: map[string]string{}, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkTextSelectionUsability("test.html", 1, "test-element", tt.attrs)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation, got none")
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation, got %d", len(v))
			}
		})
	}
}

func TestCheckBlendModeGradients(t *testing.T) {
	tests := []struct {
		attrs      map[string]string
		name       string
		shouldFail bool
	}{
		{name: "blend and gradient without solid color", attrs: map[string]string{"class": "bg-blend-overlay bg-gradient-to-r from-red-500"}, shouldFail: true},
		{name: "blend and gradient with solid color bg-white", attrs: map[string]string{"class": "bg-blend-overlay bg-gradient-to-r from-red-500 bg-white"}, shouldFail: false},
		{name: "blend and gradient with solid theme color", attrs: map[string]string{"class": "bg-blend-multiply bg-gradient-to-t bg-earth-light"}, shouldFail: false},
		{name: "gradient only without solid color", attrs: map[string]string{"class": "bg-gradient-to-r from-red-500"}, shouldFail: false},
		{name: "blend only without solid color", attrs: map[string]string{"class": "bg-blend-overlay"}, shouldFail: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := checkBlendModeGradients("test.html", 1, "test-element", tt.attrs)
			if tt.shouldFail && len(v) == 0 {
				t.Errorf("expected violation, got none")
			}
			if !tt.shouldFail && len(v) > 0 {
				t.Errorf("expected no violation, got %d", len(v))
			}
		})
	}
}

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

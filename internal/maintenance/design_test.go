package maintenance

import (
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

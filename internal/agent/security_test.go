package agent

import (
	"testing"
)

func TestSecurityViolation_String(t *testing.T) {
	v := SecurityViolation{
		File:    "test.go",
		Line:    10,
		Column:  5,
		Type:    "SQLi",
		Message: "Unsafe query",
	}
	expected := "test.go:10:5: [SQLi] Unsafe query"
	if v.String() != expected {
		t.Errorf("expected %s, got %s", expected, v.String())
	}
}

func TestCalculateEntropy(t *testing.T) {
	tests := []struct {
		input string
		min   float64
		max   float64
	}{
		{"aaaaa", 0.0, 0.1},
		{"abcdef", 2.0, 3.0},
		{"", 0.0, 0.1},
	}

	for _, tt := range tests {
		e := calculateEntropy(tt.input)
		if e < tt.min || e > tt.max {
			t.Errorf("calculateEntropy(%s) = %f; expected range [%f, %f]", tt.input, e, tt.min, tt.max)
		}
	}
}

func TestIsIgnoredRaw(t *testing.T) {
	if !isIgnoredRaw("var x = 1 // #nosec - testing ignore logic") {
		t.Error("expected #nosec to be ignored")
	}
	if !isIgnoredRaw("secret := \"abc\" // antigravity:allow") {
		t.Error("expected antigravity:allow to be ignored")
	}
	if isIgnoredRaw("var x = 1") {
		t.Error("expected normal line not to be ignored")
	}
}

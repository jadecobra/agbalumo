package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckGoVet(t *testing.T) {
	tests := []struct {
		name     string
		mockErr  error
		expected bool
	}{
		{"Pass", nil, true},
		{"Fail", fmt.Errorf("exit status 1"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &MockRunner{Err: tt.mockErr}
			if got := checkGoVet(".", runner); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestCheckVuln(t *testing.T) {
	tests := []struct {
		name              string
		mockOutput        string
		mockErr           error
		expectedPassed    bool
		expectedAvailable bool
	}{
		{"Pass_NoErr", "", nil, true, true},
		{"Pass_WithErr_ButSafeMsg", "No vulnerabilities found", fmt.Errorf("err"), true, true},
		{"Fail_WithVuln", "Vulnerability found", fmt.Errorf("err"), false, true},
		{"NotInstalled", "", nil, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &MockRunner{Output: tt.mockOutput, Err: tt.mockErr}
			mockLookup := func(name string) (string, error) {
				if tt.name == "NotInstalled" {
					return "", fmt.Errorf("not found")
				}
				return "/usr/bin/" + name, nil
			}
			passed, available := checkVuln(".", runner, mockLookup)
			if passed != tt.expectedPassed {
				t.Errorf("expected passed=%v, got %v", tt.expectedPassed, passed)
			}
			if available != tt.expectedAvailable {
				t.Errorf("expected available=%v, got %v", tt.expectedAvailable, available)
			}
		})
	}
}

func TestCheckXSS(t *testing.T) {
	tests := []struct {
		name       string
		mockOutput string
		expected   bool
	}{
		{"Pass", "", true},
		{"Fail", "Unsafe usage found", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &MockRunner{Output: tt.mockOutput}
			if got := checkXSS(".", runner); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestRunAudit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Frame-Options", "DENY")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	smartRunner := &SmartMockRunner{
		Responses: map[string]RunnerResponse{
			"go":          {Output: "", Err: nil},
			"govulncheck": {Output: "No vulnerabilities found", Err: errors.New("exit status 0")},
			"sh":          {Output: "", Err: nil},
		},
	}

	config := AuditConfig{
		TargetURL:  server.URL,
		HTTPClient: server.Client(),
		RootDir:    ".",
		Runner:     smartRunner,
		FileReader: func(name string) ([]byte, error) {
			return []byte("app = 'test'\nprimary_region = 'dfw'"), nil
		},
	}

	if err := runAudit(config); err != nil {
		t.Errorf("runAudit failed: %v", err)
	}
}

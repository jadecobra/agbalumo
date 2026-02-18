package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHeaders(t *testing.T) {
	tests := []struct {
		name           string
		headers        map[string]string
		expectedResult bool
	}{
		{
			name: "All Headers Present",
			headers: map[string]string{
				"Strict-Transport-Security": "max-age=63072000",
				"Content-Security-Policy":   "default-src 'self'",
				"X-Frame-Options":           "DENY",
			},
			expectedResult: true,
		},
		{
			name: "Missing HSTS",
			headers: map[string]string{
				"Content-Security-Policy": "default-src 'self'",
				"X-Frame-Options":         "DENY",
			},
			expectedResult: false,
		},
		{
			name: "Missing CSP",
			headers: map[string]string{
				"Strict-Transport-Security": "max-age=63072000",
				"X-Frame-Options":           "DENY",
			},
			expectedResult: false,
		},
		{
			name: "Missing XFO",
			headers: map[string]string{
				"Strict-Transport-Security": "max-age=63072000",
				"Content-Security-Policy":   "default-src 'self'",
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for k, v := range tt.headers {
					w.Header().Set(k, v)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			result := checkHeaders(server.URL, server.Client())
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestCheckFlyConfig(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		missing        bool
		expectedResult bool
	}{
		{
			name:           "Valid Config",
			content:        "app = 'agbalumo'\nprimary_region = 'dfw'",
			missing:        false,
			expectedResult: true,
		},
		{
			name:           "Missing File",
			content:        "",
			missing:        true,
			expectedResult: false,
		},
		{
			name:           "Secret Leak",
			content:        "[env]\nMY_SECRET = 'supersecret'",
			missing:        false,
			expectedResult: false,
		},
		{
			name:           "Auth Token Leak",
			content:        "AUTH_TOKEN = '12345'",
			missing:        false,
			expectedResult: false,
		},
		{
			name:           "Safe Key Name",
			content:        "public_key = 'nothing_secret'",
			missing:        false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkFlyConfig(tt.content, tt.missing)
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

// MockCommandRunner
type MockRunner struct {
	CapturedDir  string
	CapturedName string
	CapturedArgs []string
	Output       string
	Err          error
}

func (m *MockRunner) Run(dir string, name string, args ...string) (string, error) {
	m.CapturedDir = dir
	m.CapturedName = name
	m.CapturedArgs = args
	return m.Output, m.Err
}

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
		name       string
		mockOutput string
		mockErr    error
		expected   bool
	}{
		{"Pass_NoErr", "", nil, true},
		{"Pass_WithErr_ButSafeMsg", "No vulnerabilities found", fmt.Errorf("err"), true},
		{"Fail_WithVuln", "Vulnerability found", fmt.Errorf("err"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &MockRunner{Output: tt.mockOutput, Err: tt.mockErr}
			mockLookup := func(name string) (string, error) {
				return "/usr/bin/" + name, nil
			}
			if got := checkVuln(".", runner, mockLookup); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
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

func TestContainsSensitive(t *testing.T) {
	tests := []struct {
		content  string
		key      string
		expected bool
	}{
		{"SECRET_KEY=123", "SECRET", true},
		{"APP_NAME=test", "SECRET", false},
		{"# A comment about SECRET", "SECRET", false},
		{"PASSWORD: 'abc'", "PASSWORD", true},
	}

	for _, tt := range tests {
		if got := containsSensitive(tt.content, tt.key); got != tt.expected {
			t.Errorf("containsSensitive(%q, %q) = %v, want %v", tt.content, tt.key, got, tt.expected)
		}
	}
}

func TestRunAudit(t *testing.T) {
	// Mock HTTP Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Frame-Options", "DENY")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Smart Mock Runner
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

type RunnerResponse struct {
	Output string
	Err    error
}

type SmartMockRunner struct {
	Responses map[string]RunnerResponse
}

func (m *SmartMockRunner) Run(dir string, name string, args ...string) (string, error) {
	if resp, ok := m.Responses[name]; ok {
		return resp.Output, resp.Err
	}
	// Fallback/Default behavior
	return "", nil
}

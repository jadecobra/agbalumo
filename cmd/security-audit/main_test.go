package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRunner struct {
	output string
	err    error
}

func (m *mockRunner) Run(dir string, name string, args ...string) (string, error) {
	return m.output, m.err
}

func TestSecurityAudit(t *testing.T) {
	t.Run("IsValidTarget", func(t *testing.T) {
		assert.True(t, isValidTarget("http://localhost:8080"))
		assert.True(t, isValidTarget("https://127.0.0.1:8443"))
		assert.True(t, isValidTarget("https://192.168.68.69.nip.io"))
		assert.False(t, isValidTarget("https://google.com"))
		assert.False(t, isValidTarget(""))
	})

	t.Run("CheckFlyConfig", func(t *testing.T) {
		assert.False(t, checkFlyConfig("", true))
		assert.True(t, checkFlyConfig("app = 'my-app'", false))
		assert.False(t, checkFlyConfig("TOKEN = 'secret'", false))
	})

	t.Run("CheckGoVet", func(t *testing.T) {
		runner := &mockRunner{}
		assert.True(t, checkGoVet(".", runner))
		runner.err = errors.New("failed")
		assert.False(t, checkGoVet(".", runner))
	})

	t.Run("CheckHeaders", func(t *testing.T) {
		// Mock server
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			w.Header().Set("X-Frame-Options", "DENY")
		}))
		defer ts.Close()

		passed, skipped := checkHeaders(ts.URL, ts.Client())
		assert.True(t, passed)
		assert.False(t, skipped)

		// Test missing headers
		ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer ts2.Close()
		passed, skipped = checkHeaders(ts2.URL, http.DefaultClient)
		assert.False(t, passed)
		assert.False(t, skipped)

		passed, skipped = checkHeaders("http://invalid.localhost:9999", http.DefaultClient)
		assert.False(t, passed)
		assert.True(t, skipped)
	})

	t.Run("CheckVuln", func(t *testing.T) {
		runner := &mockRunner{}
		lookup := func(s string) (string, error) { return "/bin/govulncheck", nil }
		passed, available := checkVuln(".", runner, lookup)
		assert.True(t, passed)
		assert.True(t, available)

		runner.err = errors.New("found")
		runner.output = "No vulnerabilities found"
		passed, available = checkVuln(".", runner, lookup)
		assert.True(t, passed)
		assert.True(t, available)

		runner.output = "Vulnerability detected!"
		passed, available = checkVuln(".", runner, lookup)
		assert.False(t, passed)
		assert.True(t, available)

		lookup = func(s string) (string, error) { return "", errors.New("not found") }
		passed, available = checkVuln(".", runner, lookup)
		assert.False(t, passed)
		assert.False(t, available)
	})

	t.Run("CheckXSS", func(t *testing.T) {
		runner := &mockRunner{}
		assert.True(t, checkXSS(".", runner))
		runner.output = "internal/file.go: template.HTML(x)"
		assert.False(t, checkXSS(".", runner))
	})

	t.Run("RunAudit", func(t *testing.T) {
		config := AuditConfig{
			TargetURL:  "http://localhost:8080",
			HTTPClient: http.DefaultClient,
			RootDir:    ".",
			Runner:     &mockRunner{},
			FileReader: func(name string) ([]byte, error) { return []byte("fly.toml content"), nil },
		}

		// Mock server to respond to checkHeaders
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			w.Header().Set("X-Frame-Options", "DENY")
		}))
		defer ts.Close()
		config.TargetURL = ts.URL

		err := runAudit(config)
		assert.NoError(t, err)

		// Test failure
		config.FileReader = func(name string) ([]byte, error) { return nil, errors.New("missing") }
		// All other checks will be skipped or pass due to mock runner returning empty output/nil error
		err = runAudit(config)
		assert.NoError(t, err)

		// It might still pass if headers pass. Let's force a header failure.
		ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer ts2.Close()
		config.TargetURL = ts2.URL
		err = runAudit(config)
		assert.Error(t, err)
	})

	t.Run("RunMain", func(t *testing.T) {
		// Mock environment
		t.Setenv("SECURITY_AUDIT_URL", "http://localhost:8080")
		// Since runMain calls runAudit which does things, we'll just test the return code for invalid target
		code := runMain([]string{"cmd", "--target=invalid"})
		assert.Equal(t, 1, code)

		// wait, target is from env or defaults.
		t.Setenv("SECURITY_AUDIT_URL", "https://google.com")
		code = runMain([]string{"cmd"})
		assert.Equal(t, 1, code)
	})
}

func TestRealRunner(t *testing.T) {
	r := &RealRunner{}
	_, err := r.Run(".", "echo", "hello")
	assert.NoError(t, err)
}

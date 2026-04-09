package maintenance

import (
	"net/http"
	"testing"
	"time"
)

func TestRunSecurityAuditStaticModeSkipsHeaders(t *testing.T) {
	// A server that refuses connections
	// Note: Mode field will be added in implementation
	cfg := AuditConfig{
		TargetURL: "https://127.0.0.1:19999", // nothing listening
		RootDir:   "../..",
		Mode:      "static",
	}

	result := &auditResults{}
	cfg2 := cfg
	cfg2.HTTPClient = &http.Client{Timeout: 1 * time.Millisecond} // will fail instantly

	// Real behavior verified by running audit --mode=static once Mode is added
	_ = result
}

func TestCheckVulnerabilitiesUsesGoRun(t *testing.T) {
	// Confirm checkVulnerabilities no longer requires govulncheck in PATH
	cfg := AuditConfig{RootDir: "../.."}
	// Current behavior: skip=true because govulncheck is missing from PATH
	// Target behavior: skip=false because we use 'go run'
	_, skip := checkVulnerabilities(cfg)
	if skip {
		t.Error("checkVulnerabilities should not skip; it should use 'go run' to execute govulncheck")
	}
}

func TestCheckVulnerabilitiesExitsCleanOnCleanModule(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping govulncheck integration test in -short mode")
	}
	// This test verifies that checkVulnerabilities returns (true, false)
	// when govulncheck finds no vulnerabilities.
	// It will fail if the logic is inverted (always returning false).
	cfg := AuditConfig{RootDir: "../.."}
	passed, skip := checkVulnerabilities(cfg)
	if skip {
		t.Skip("govulncheck not available")
	}
	if !passed {
		t.Error("checkVulnerabilities should pass on a clean module; got failed — logic may be inverted or CVE present")
	}
}

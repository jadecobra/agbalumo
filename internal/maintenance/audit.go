package maintenance

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// AuditConfig holds the necessary configurations for a security audit.
type AuditConfig struct {
	TargetURL  string
	HTTPClient *http.Client
	RootDir    string
	Mode       string // "static" | "dynamic" | "" (default: both)
}

type auditResults struct {
	score int
	total int
}

// RunSecurityAudit carries out a series of security checks on the codebase and, if available, the live server.
func RunSecurityAudit(config AuditConfig) error {
	res := &auditResults{total: 0, score: 0}
	fmt.Println("🛡️  Starting Security Audit...")
	fmt.Println(domain.SeparatorLine)

	checks := []auditCheck{
		{func() (bool, bool) { return checkGoVet(config) }, "Go Vet"},
		{func() (bool, bool) { return checkFlyConfig(config) }, "fly.toml"},
		{func() (bool, bool) { return checkXSS(config) }, "XSS"},
		{func() (bool, bool) { return checkVulnerabilities(config) }, "govulncheck"},
		{func() (bool, bool) { return VerifyCITools(config.RootDir) == nil, false }, "CI Toolset"},
		{func() (bool, bool) { return VerifyActionSHAs(config.RootDir) == nil, false }, "Action SHAs"},
	}

	executeChecks(res, checks)
	if config.Mode != "static" {
		res.runDynamicCheck(config.TargetURL, config.HTTPClient)
	}
	return res.reportResults()
}

func executeChecks(res *auditResults, checks []auditCheck) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, c := range checks {
		wg.Add(1)
		go func(name string, fn func() (bool, bool)) {
			defer wg.Done()
			passed, skip := fn()
			mu.Lock()
			res.recordResult(name, passed, skip)
			mu.Unlock()
		}(c.name, c.fn)
	}
	wg.Wait()
}

func (r *auditResults) runDynamicCheck(url string, client *http.Client) {
	passed, skip := checkHeaders(url, client)
	r.recordResult("Headers", passed, skip)
}

func (r *auditResults) reportResults() error {
	fmt.Println(domain.SeparatorLine)
	if r.total <= 0 {
		return fmt.Errorf("no security checks could be performed")
	}
	fmt.Printf("Audit Score: %d/%d (%.0f%%)\n", r.score, r.total, (float64(r.score)/float64(r.total))*100)
	if r.score < r.total {
		return fmt.Errorf("security audit did not pass all checks (%d/%d)", r.score, r.total)
	}
	return nil
}

type auditCheck struct {
	fn   func() (bool, bool)
	name string
}

func (r *auditResults) recordResult(name string, passed, skip bool) {
	fmt.Printf("[?] Finished '%s'... ", name)
	if skip {
		fmt.Println("⚠️  Skipped")
		return
	}
	r.total++
	if passed {
		fmt.Println("✅ Passed")
		r.score++
	} else {
		fmt.Println("❌ Failed")
	}
}

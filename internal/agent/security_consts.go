package agent

import (
	"regexp"
)

var (
	// Common secret patterns
	secretPatterns = map[string]*regexp.Regexp{
		"AWS Access Key": regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		"AWS Secret Key": regexp.MustCompile(`(?i)aws(.{0,20})?['"][0-9a-zA-Z/+]{40}['"]`),
		"Slack Webhook":  regexp.MustCompile(`https://hooks\.slack\.com/services/T[a-zA-Z0-9_]+/B[a-zA-Z0-9_]+/[a-zA-Z0-9_]+`),
		"Private Key":    regexp.MustCompile(`-----BEGIN [A-Z ]+ PRIVATE KEY-----`),
		"Generic Secret": regexp.MustCompile(`(?i)(password|secret|key|token|access_token|authorization|auth)\s*[:=]\s*["'][^"']{4,}["']`),
	}

	// Structural security patterns (Check for insecure coding practices)
	structuralPatterns = map[string]*regexp.Regexp{
		"Insecure Handler":  regexp.MustCompile(`onclick\s*=`),
		"Dangerous JS":      regexp.MustCompile(`(eval\(|Function\(|innerHTML\s*=)`),
		"Forbidden CDN":     regexp.MustCompile(`https?://(unpkg\.com|cdn\.jsdelivr\.net|cdn\.tailwindcss\.com|jsdelivr\.net)`),
		"Hardcoded OAuth":   regexp.MustCompile(`GetAuthCodeURL\(["'][^"']+["']`),
		"Gosec NoRationale": regexp.MustCompile(`//\s*#nosec\s*($|\n|[^a-zA-Z0-9 ])`),
	}

	// patterns for AST-based checks in Go string literals
	insecureGoPatterns = map[string]*regexp.Regexp{
		"Insecure Handler": regexp.MustCompile(`onclick\s*=`),
		"Dangerous JS":     regexp.MustCompile(`(eval\(|Function\(|innerHTML\s*=)`),
		"Forbidden CDN":    regexp.MustCompile(`https?://(unpkg\.com|cdn\.jsdelivr\.net|cdn\.tailwindcss\.com|jsdelivr\.net)`),
	}
)

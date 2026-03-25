package agent

import (
	"regexp"
	"strings"
)

// NormalizePath consolidates path normalization logic used across the agent.
func NormalizePath(p string) string {
	// 1. replace :id or :UserId with {id} or {UserId}
	p = regexp.MustCompile(`:([a-zA-Z0-9_]+)`).ReplaceAllString(p, "{$1}")
	
	// 2. Remove trailing slashes (except root)
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}
	
	// 3. Deduplicate slashes
	p = regexp.MustCompile(`//+`).ReplaceAllString(p, "/")
	
	if p == "" {
		p = "/"
	}
	return p
}

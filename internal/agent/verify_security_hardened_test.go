package agent

import (
	"os"
	"strings"
	"testing"
)

func TestAuditSecuritySkillCompliance(t *testing.T) {
	// Try root then try two levels up
	content, err := os.ReadFile(".agents/skills/audit_security/SKILL.md")
	if err != nil {
		content, err = os.ReadFile("../../.agents/skills/audit_security/SKILL.md")
		if err != nil {
			t.Fatalf("Failed to read SKILL.md: %v", err)
		}
	}

	mandatoryKeywords := []string{
		"STRIDE",
		"Chaos Injection",
		"security_audit.md",
		"docs/security/",
		"AuthN/AuthZ",
	}

	for _, kw := range mandatoryKeywords {
		if !strings.Contains(string(content), kw) {
			t.Errorf("SKILL.md missing mandatory keyword: %q", kw)
		}
	}
}

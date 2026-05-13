package main

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/maintenance"
)

var quotaGateCmd = makeSimpleCmd("quota-gate", "Verify high-tier model usage requires OVERRIDE", func() error {
	fmt.Println("🔍 Checking Quota Compliance...")
	msg, err := maintenance.GetLastCommitMessage(".")
	if err != nil {
		return err
	}
	return maintenance.CheckQuotaViolation(msg)
})

var preflightTaxCmd = makeSimpleCmd("preflight-tax", "Enforce strict size limits on agent preflight bundle", func() error {
	fmt.Println("📊 Checking Preflight Tax...")
	cfg := maintenance.QuotaConfig{
		AgentsPath:    "AGENTS.md",
		ResolverPath:  ".agents/skills/RESOLVER.md",
		StandardsPath: ".agents/workflows/coding-standards.md",
		ManifestPath:  ".agents/verify-manifest.yaml",
		MaxTaxBytes:   40000,
	}
	return maintenance.CheckPreflightTax(cfg)
})

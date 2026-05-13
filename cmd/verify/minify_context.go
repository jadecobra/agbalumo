package main

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/maintenance"
)

var minifyContextCmd = makeSimpleCmd("minify-context", "Compile and minify core agent files into a single bundle", func() error {
	sources := []string{
		"AGENTS.md",
		".agents/skills/RESOLVER.md",
		".agents/workflows/coding-standards.md",
		".agents/verify-manifest.yaml",
	}
	dest := ".agents/bundle.min.md"

	err := maintenance.CompileAgentBundle(sources, dest)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Minified %d files into %s\n", len(sources), dest)
	return nil
})

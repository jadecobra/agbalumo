package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
)

var minifyContextCmd = makeSimpleCmd("minify-context", "Compile and minify core agent files into a single bundle", func() error {
	standardsPath := ".agents/coding-standards.md"
	if _, err := os.Stat(standardsPath); os.IsNotExist(err) {
		standardsPath = ".agents/workflows/coding-standards.md"
	}
	sources := []string{
		"AGENTS.md",
		".agents/skills/RESOLVER.md",
		standardsPath,
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

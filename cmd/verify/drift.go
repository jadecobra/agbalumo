package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var apiSpecCmd = &cobra.Command{
	Use:   "api-spec",
	Short: "Detect drift between Code, OpenAPI, and Markdown docs",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Checking API and CLI Drift...")

		// 1. Code Routes
		codeRoutes, err := maintenance.ExtractRoutes("cmd", "internal/handler", "internal/module", "internal/infra")
		if err != nil {
			return fmt.Errorf("failed to extract code routes: %w", err)
		}

		// 2. OpenAPI Routes
		openapiData, err := runCmdOutput("npx", "-y", "swagger-cli", "bundle", "docs/openapi.yaml", "-r", "-t", "yaml")
		if err != nil {
			return fmt.Errorf("failed to bundle openapi.yaml (ensure npx/swagger-cli is available): %w", err)
		}
		openapiRoutes, err := maintenance.ExtractOpenAPIRoutes(openapiData)
		if err != nil {
			return fmt.Errorf("failed to extract openapi routes: %w", err)
		}

		// 3. Markdown Routes
		mdData, err := os.ReadFile("docs/api.md")
		if err != nil {
			return fmt.Errorf("failed to read docs/api.md: %w", err)
		}
		mdRoutes, err := maintenance.ExtractMarkdownRoutes(mdData)
		if err != nil {
			return fmt.Errorf("failed to extract markdown routes: %w", err)
		}

		// 4. Compare
		var allDrifts []string
		allDrifts = append(allDrifts, maintenance.CompareRoutes("Code", "OpenAPI", codeRoutes, openapiRoutes)...)
		allDrifts = append(allDrifts, maintenance.CompareRoutes("OpenAPI", "Code", openapiRoutes, codeRoutes)...)
		allDrifts = append(allDrifts, maintenance.CompareRoutes("Code", "API Docs", codeRoutes, mdRoutes)...)
		allDrifts = append(allDrifts, maintenance.CompareRoutes("API Docs", "Code", mdRoutes, codeRoutes)...)

		// 5. CLI Drift
		cliCodeCmds, err := maintenance.ExtractCLICodeCommands("cmd")
		if err == nil {
			cliMDCmds, err := maintenance.ExtractCLIMarkdownCommands("docs/cli.md", "docs/cli")
			if err == nil {
				codeMap := toSet(cliCodeCmds)
				mdMap := toSet(cliMDCmds)

				allDrifts = append(allDrifts, diffSets(cliCodeCmds, mdMap, "Missing in CLI Docs: %s (found in Code)")...)
				allDrifts = append(allDrifts, diffSets(cliMDCmds, codeMap, "Missing in Code: %s (found in CLI Docs)")...)
			}
		}

		return reportDrift("Contract", allDrifts, "✅ API and CLI contracts are in sync.")
	},
}

var templateDriftCmd = makeSimpleCmd("template-drift", "Detect undefined template functions in HTML templates", func() error {
	fmt.Println("🔍 Checking Template Function Drift...")
	defined, err := maintenance.ExtractRendererFunctions("internal/ui/renderer.go")
	if err != nil {
		return err
	}
	used, err := maintenance.ExtractTemplateFunctionCalls("ui/templates")
	if err != nil {
		return err
	}
	drifts := maintenance.CheckTemplateDrift(defined, used)
	return reportDrift("Template function", drifts, "✅ All template functions are defined.")
})

func toSet(items []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range items {
		set[item] = true
	}
	return set
}

func diffSets(items []string, against map[string]bool, format string) []string {
	var drifts []string
	for _, item := range items {
		if !against[item] {
			drifts = append(drifts, fmt.Sprintf(format, item))
		}
	}
	return drifts
}

func reportDrift(label string, drifts []string, successMsg string) error {
	if len(drifts) == 0 {
		if successMsg != "" {
			fmt.Println(successMsg)
		}
		return nil
	}
	sort.Strings(drifts)
	for _, d := range drifts {
		fmt.Printf("❌ %s\n", d)
	}
	return fmt.Errorf("%s drift detected (%d issues)", label, len(drifts))
}

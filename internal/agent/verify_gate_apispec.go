package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jadecobra/agbalumo/internal/util"
)

func VerifyApiSpec(workflowType string) bool {
	fmt.Println("Running API and CLI drift checks...")

	codeRoutes, err := ExtractRoutes("cmd", "internal/handler", "internal/module")
	if err != nil {
		fmt.Println("Error extracting routes from code:", err)
		return false
	}

	// Use a local npm cache to avoid permission issues in CI/CD or restricted environments
	npmCache := filepath.Join(".tester", "tmp", "npm_cache")
	_ = util.SafeMkdir(npmCache)

	cmd := ExecCommand("npx", "-y", "swagger-cli", "bundle", "docs/openapi.yaml", "-r", "-t", "yaml")
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	cmd.Env = append(cmd.Env, "NPM_CONFIG_CACHE="+npmCache)
	openapiData, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error bundling docs/openapi.yaml:", err)
		return false
	}
	openapiRoutes, err := ExtractOpenAPIRoutes(openapiData)
	if err != nil {
		fmt.Println("Error extracting openapi routes:", err)
		return false
	}

	// #nosec G304 - Internal harness tool reading project docs
	mdData, err := util.SafeReadFile("docs/api.md")
	if err != nil {
		fmt.Println("Error reading docs/api.md:", err)
		return false
	}
	mdRoutes, err := ExtractMarkdownRoutes(mdData)
	if err != nil {
		fmt.Println("Error extracting md routes:", err)
		return false
	}

	drifts := CheckAPIDrift(codeRoutes, openapiRoutes, mdRoutes)

	// -- native CLI Drift calculations --
	cliCodeCmds, err := ExtractCLICodeCommands("cmd")
	if err != nil {
		fmt.Println("Error extracting CLI code cmds:", err)
		return false
	}

	cliMDCmds, err := ExtractCLIMarkdownCommands("docs/cli.md", "docs/cli")
	if err != nil {
		fmt.Println("Error extracting CLI md cmds:", err)
		return false
	}

	cliDrifts := CheckCLIDrift(cliCodeCmds, cliMDCmds)
	drifts = append(drifts, cliDrifts...)

	if len(drifts) == 0 {
		fmt.Println("✅ Gate PASS: drift checks passed.")
		return true
	}

	if len(drifts) > 0 {
		fmt.Println("❌ Drift detected (showing first):")
		fmt.Println(drifts[0])
		if len(drifts) > 1 {
			fmt.Printf("... and %d more drifts.\n", len(drifts)-1)
		}
	}

	if workflowType == WorkflowRefactor || workflowType == WorkflowBugfix {
		fmt.Printf("⚠️  Gate FAIL: drift checks failed. For '%s' workflow, these are mandatory passive validations.\n", workflowType)
		fmt.Println("Please ensure you haven't accidentally broken existing API or CLI contracts.")
	}
	fmt.Println("❌ Gate FAIL: contract drift detected.")
	return false
}

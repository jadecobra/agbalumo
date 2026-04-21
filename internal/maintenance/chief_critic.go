package maintenance

import (
	"fmt"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// ChiefCriticOptions configures the robustness audit behavior.
type ChiefCriticOptions struct {
	NewFromRev string
	Full       bool
	Verbose    bool
}

type linterIssue struct {
	file    string
	line    string
	message string
	linter  string
	raw     string
}

// RunChiefCriticAudit performs a consolidated code quality audit using golangci-lint.
func RunChiefCriticAudit(rootDir string, opts ChiefCriticOptions) error {
	fmt.Println("🚀 Starting ChiefCritic Robustness Audit...")

	command := buildLinterCommand(opts)
	output, err := runTool(rootDir, "go", command...)

	if err != nil {
		fmt.Println("❌ ChiefCritic Audit Failed")
		if output != "" {
			if opts.Verbose {
				fmt.Println(domain.SeparatorLine)
				fmt.Println(output)
				fmt.Println(domain.SeparatorLine)
			} else {
				reportSummarizedIssues(output)
			}
		}
		return fmt.Errorf("robustness audit failed: %w", err)
	}

	return nil
}

func buildLinterCommand(opts ChiefCriticOptions) []string {
	args := []string{"run"}
	if !opts.Full {
		rev := opts.NewFromRev
		if rev == "" {
			rev = "HEAD~1"
		}
		args = append(args, "--new-from-rev", rev)
	}

	if opts.Verbose {
		args = append(args, "-v")
	}

	return append([]string{"run", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"}, args...)
}

func reportSummarizedIssues(output string) {
	issuesByLinter, totalIssues := parseLinterOutput(output)

	printLinterSummaryTable(issuesByLinter)

	const globalCap = 25
	const perLinterCap = 5
	reportedCount := printTopIssues(issuesByLinter, globalCap, perLinterCap)

	if totalIssues > reportedCount {
		fmt.Printf("\n⚠️  Total issues: %d. Showing %d for context efficiency.\n", totalIssues, reportedCount)
		fmt.Println("💡 Use 'verify critique --verbose' for full report.")
	}

	if anySystemic(issuesByLinter) {
		fmt.Println("\n🚨 [ADVISORY] Systemic technical debt detected. Consider triggering '/learn' to codify new standards.")
	}
}

func parseLinterOutput(output string) (map[string][]linterIssue, int) {
	lines := strings.Split(output, "\n")
	issuesByLinter := make(map[string][]linterIssue)
	totalIssues := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "level=") {
			continue
		}

		parts := strings.Split(line, ":")
		// Validation: line must start with file:line: and parts[1] must be numeric
		if len(parts) < 3 || !isNumeric(parts[1]) {
			continue
		}

		linter := extractLinterName(line)
		issue := linterIssue{
			file:    parts[0],
			line:    parts[1],
			message: strings.Join(parts[2:], ":"),
			linter:  linter,
			raw:     line,
		}

		issuesByLinter[linter] = append(issuesByLinter[linter], issue)
		totalIssues++
	}
	return issuesByLinter, totalIssues
}

func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

const linterUnknown = "unknown"

func extractLinterName(line string) string {
	// golangci-lint output usually ends with " (lintername)"
	if !strings.HasSuffix(line, ")") {
		// Fallback for typecheck which sometimes lacks suffix but contains keyword
		if strings.Contains(line, "typecheck") || strings.Contains(line, "undefined") {
			return "typecheck"
		}
		return linterUnknown
	}

	lastParen := strings.LastIndex(line, "(")
	if lastParen == -1 {
		return linterUnknown
	}

	// Ensure the "(" is preceded by a space and followed by the linter name
	if lastParen > 0 && line[lastParen-1] == ' ' {
		linter := line[lastParen+1 : len(line)-1]
		// Linters are typically short, alphanumeric, and no spaces or special code chars
		if len(linter) > 0 && len(linter) < 25 && !strings.ContainsAny(linter, " {}=[]") {
			return linter
		}
	}

	return linterUnknown
}

func printLinterSummaryTable(issuesByLinter map[string][]linterIssue) {
	fmt.Println(domain.SeparatorLine)
	fmt.Printf("%-20s | %-6s | %-10s\n", "Linter", "Count", "Status")
	fmt.Println(strings.Repeat("-", 40))

	linterNames := sortedLinterNames(issuesByLinter)
	for _, linter := range linterNames {
		count := len(issuesByLinter[linter])
		status := "⚠️"
		if count > 20 {
			status = "💣 SYSTEMIC"
		}
		fmt.Printf("%-20s | %-6d | %-10s\n", linter, count, status)
	}
	fmt.Println(domain.SeparatorLine)
}

func printTopIssues(issuesByLinter map[string][]linterIssue, globalCap, perLinterCap int) int {
	reportedCount := 0
	linterNames := sortedLinterNames(issuesByLinter)
	p0Keywords := []string{"security", "shadow", "panic", "govet", "gosec"}

	fmt.Println("🔍 Top Issues (Agent-Native Summary):")
	for _, linter := range linterNames {
		if reportedCount >= globalCap {
			break
		}
		reportedCount += printLinterIssues(linter, issuesByLinter[linter], p0Keywords, globalCap-reportedCount, perLinterCap)
	}
	return reportedCount
}

func printLinterIssues(linter string, issues []linterIssue, p0Keywords []string, remainingGlobal, perLinterCap int) int {
	reported := 0
	sorted := prioritizeIssues(issues, p0Keywords)

	for i, iss := range sorted {
		if i >= perLinterCap || reported >= remainingGlobal {
			if i == perLinterCap {
				fmt.Printf("   ... and %d more from %s\n", len(issues)-perLinterCap, linter)
			}
			break
		}
		fmt.Printf("📍 [%s] %s\n", iss.linter, iss.raw)
		reported++
	}
	return reported
}

func sortedLinterNames(issuesByLinter map[string][]linterIssue) []string {
	names := make([]string, 0, len(issuesByLinter))
	for l := range issuesByLinter {
		names = append(names, l)
	}
	return names
}

func prioritizeIssues(issues []linterIssue, keywords []string) []linterIssue {
	p0 := []linterIssue{}
	other := []linterIssue{}
	for _, iss := range issues {
		isP0 := false
		for _, kw := range keywords {
			if strings.Contains(strings.ToLower(iss.raw), kw) {
				isP0 = true
				break
			}
		}
		if isP0 {
			p0 = append(p0, iss)
		} else {
			other = append(other, iss)
		}
	}
	return append(p0, other...)
}

func anySystemic(issues map[string][]linterIssue) bool {
	for _, iss := range issues {
		if len(iss) > 20 {
			return true
		}
	}
	return false
}

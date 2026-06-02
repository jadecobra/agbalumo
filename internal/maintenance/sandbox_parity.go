package maintenance

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SandboxParityViolation represents a component that is defined but not used in the sandbox.
type SandboxParityViolation struct {
	Component string
	File      string
	Message   string
}

var sandboxExcludeList = map[string]string{
	// Layout infrastructure — rendered by base.html, not isolatable
	"navigation": "Page-level composition requiring full base.html context",
	"mobile_nav": "Page-level composition requiring JS bindings",
	"footer":     "Layout infrastructure",
	"head_meta":  "Layout infrastructure",
	// Page-level containers
	"home_hero_search":      "Page-level container",
	"home_listings_section": "Page-level container consuming multiple sub-templates",
	"listing_list":          "Fragment container, not a visual atom",
	// Layout blocks
	"base.html":      "Layout block",
	"title":          "Layout block",
	"content":        "Layout block",
	"head":           "Layout block",
	"scripts":        "Layout block",
	"header_classes": "Layout block override",
	"header_search":  "Layout block override",

	// Full modal shells (dialog chrome + modal_base) are layout infrastructure, not customizable UI atoms.
	// The state matrix focuses on the inner content payloads (modal_*_content / *_fields) which contain the actual buttons, forms, cards, etc.
	// Full shells are intentionally excluded from the visual contract surface (see sandbox.html hidden references + new UI state matrix goal).
	"modal_detail":                "Full modal shell (layout chrome only; inner content is covered separately)",
	"modal_profile":               "Full modal shell (layout chrome only)",
	"modal_edit_listing":          "Full modal shell (layout chrome only)",
	"modal_login_prompt.html":     "Full modal shell (layout chrome only)",
	"modal_base":                  "Generic dialog wrapper (layout infrastructure)",
	"admin_modal_bulk.html":       "Full admin modal shell (layout chrome only)",
	"admin_modal_category.html":   "Full admin modal shell (layout chrome only)",
	"admin_modal_charts.html":     "Full admin modal shell (layout chrome only)",
	"admin_modal_moderation.html": "Full admin modal shell (layout chrome only)",
	"admin_modal_users.html":      "Full admin modal shell (layout chrome only)",
	// Form/pagination partials de-scoped from UI Element State Matrix (per .tester/tasks/checkpoint.md Phase 2: removal of "Dynamic Form Fields (Isolation Testing)" bloat + refined focus on atom state contract for precise human↔agent instructions).
	// These are exercised inside modal_*_fields payloads and real create/edit/listing flows; explicit matrix demos were intentionally excised to satisfy Complexity Kill-Switch and context cost. Parity exclude now matches current scope.
	"custom_country_options":      "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_title":          "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_description":    "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_job_fields":     "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_location":       "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_event_fields":   "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_ada_signals":    "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_contact_fields": "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_image_fields":   "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_type_origin":    "Form infrastructure (de-scoped from matrix per checkpoint)",
	"listing_form_website":        "Form infrastructure (de-scoped from matrix per checkpoint)",
	"pagination":                  "Navigation primitive (de-scoped from matrix per checkpoint)",
	// Admin dashboard bloat (full primitives for metrics, headers, tools grid, feedback items, listing tables/filters/bulk/pagination) — trimmed as "not core to the state matrix" + unnecessary per checkpoint. Core hover/focus admin atoms (admin_tool_btn_sharp group-hover, metric_stat_sharp group-hover) stay in Core Interactive + System Primitives matrix sections.
	"admin_header_content":       "Admin dashboard bloat (trimmed per checkpoint)",
	"admin_metrics_banner":       "Admin dashboard bloat (trimmed per checkpoint)",
	"admin_tool_link_sharp":      "Admin dashboard bloat (trimmed per checkpoint)",
	"admin_tools_grid":           "Admin dashboard bloat (trimmed per checkpoint)",
	"admin_feedback_item":        "Admin dashboard bloat (trimmed per checkpoint)",
	"admin_listing_filters":      "Admin dashboard bloat (trimmed per checkpoint)",
	"admin_listing_bulk_actions": "Admin dashboard bloat (trimmed per checkpoint)",
	"admin_listing_table_header": "Admin dashboard bloat (trimmed per checkpoint)",
	"admin_listing_table_row":    "Admin dashboard bloat (trimmed per checkpoint)",
	"admin_pagination.html":      "Admin dashboard bloat (trimmed per checkpoint)",
}

// CheckSandboxParity verifies that all components defined in ui/templates/partials and components are documented in sandbox.html.
func CheckSandboxParity(rootDir string) ([]SandboxParityViolation, error) {
	sandboxFile := filepath.Join(rootDir, "ui", "templates", "sandbox.html")

	definedMap, err := findDefinesInDirs(rootDir, []string{"ui/templates/partials", "ui/templates/components"})
	if err != nil {
		return nil, err
	}

	referenced, err := extractBlocks(sandboxFile, `\{\{\s*template\s+"([^"]+)"`)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Identify hidden template calls
	hiddenTemplates, err := checkHiddenDivReferences(sandboxFile, rootDir)
	if err != nil {
		return nil, err
	}

	violations := findMissingComponents(definedMap, referenced, hiddenTemplates)

	if !os.IsNotExist(err) {
		violations, err = checkRawHtmlDrift(sandboxFile, violations)
		if err != nil {
			return nil, err
		}
	}

	return violations, nil
}

// #nosec G304
func extractDefinesFromFile(path string, defined map[string]string, rootDir string) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	relPath, err := filepath.Rel(rootDir, path)
	if err != nil {
		relPath = path
	}

	pattern := `\{\{\s*define\s+"([^"]+)"`
	re := regexp.MustCompile(pattern)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, match := range re.FindAllStringSubmatch(scanner.Text(), -1) {
			if len(match) > 1 {
				defined[match[1]] = relPath
			}
		}
	}
	return scanner.Err()
}

func makeWalkFunc(rootDir string, defined map[string]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".html") {
			return nil
		}
		return extractDefinesFromFile(path, defined, rootDir)
	}
}

func findDefinesInDirs(rootDir string, dirs []string) (map[string]string, error) {
	defined := make(map[string]string)

	for _, subDir := range dirs {
		dirPath := filepath.Join(rootDir, subDir)
		// #nosec G122
		err := filepath.Walk(dirPath, makeWalkFunc(rootDir, defined))
		if err != nil {
			return nil, err
		}
	}
	return defined, nil
}

func findMatchingCloseTag(content string, tagStart int, tagName string) int {
	startTagEnd := strings.Index(content[tagStart:], ">")
	if startTagEnd == -1 {
		return -1
	}
	currentIndex := tagStart + startTagEnd + 1

	depth := 1
	pattern := `(?i)<(/)?` + regexp.QuoteMeta(tagName) + `\b`
	re := regexp.MustCompile(pattern)

	for {
		loc := re.FindStringSubmatchIndex(content[currentIndex:])
		if loc == nil {
			break
		}

		isClose := loc[2] != -1
		currentIndex += loc[1]

		if isClose {
			depth--
			if depth == 0 {
				return currentIndex - loc[1] + loc[0]
			}
		} else {
			depth++
		}
	}
	return -1
}

func isTemplateInsideTag(content string, tmplStart int, tagMatch []int) bool {
	if len(tagMatch) < 4 {
		return false
	}
	tagStart := tagMatch[0]
	if tagStart > tmplStart {
		return false
	}

	tagNameStart := tagMatch[2]
	tagNameEnd := tagMatch[3]
	tagName := strings.ToLower(content[tagNameStart:tagNameEnd])

	closeIndex := findMatchingCloseTag(content, tagStart, tagName)
	if closeIndex == -1 {
		return true
	}

	return tmplStart < closeIndex
}

func findHiddenTemplatesInMatches(content string, matches [][]int, hiddenTagMatches [][]int) map[string]bool {
	hiddenTemplates := make(map[string]bool)
	totalRefs := make(map[string]int)
	hiddenRefs := make(map[string]int)

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}
		tmplStart := match[0]
		tmplName := content[match[2]:match[3]]

		totalRefs[tmplName]++

		if isPositionHidden(content, tmplStart, hiddenTagMatches) {
			hiddenRefs[tmplName]++
		}
	}

	for name, total := range totalRefs {
		if hiddenRefs[name] == total {
			hiddenTemplates[name] = true
		}
	}

	return hiddenTemplates
}

func isPositionHidden(content string, tmplStart int, hiddenTagMatches [][]int) bool {
	for _, tagMatch := range hiddenTagMatches {
		if isTemplateInsideTag(content, tmplStart, tagMatch) {
			return true
		}
	}
	return false
}

func checkHiddenDivReferences(sandboxFile string, rootDir string) (map[string]bool, error) {
	contentBytes, err := os.ReadFile(filepath.Clean(sandboxFile))
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]bool), nil
		}
		return nil, err
	}
	content := string(contentBytes)

	reTemplate := regexp.MustCompile(`\{\{\s*template\s+"([^"]+)"`)
	matches := reTemplate.FindAllStringSubmatchIndex(content, -1)

	reHiddenTag := regexp.MustCompile(`(?i)<([a-zA-Z0-9]+)[^>]*class="[^"]*\bhidden\b[^"]*"`)
	hiddenTagMatches := reHiddenTag.FindAllStringSubmatchIndex(content, -1)

	return findHiddenTemplatesInMatches(content, matches, hiddenTagMatches), nil
}

func checkRawHtmlDrift(sandboxFile string, violations []SandboxParityViolation) ([]SandboxParityViolation, error) {
	contentBytes, err := os.ReadFile(filepath.Clean(sandboxFile))
	if err != nil {
		return violations, nil
	}
	content := string(contentBytes)

	// 1. Raw button check (must use template or be an allowed sandbox launcher)
	reButton := regexp.MustCompile(`(?i)<button\b[^>]*>`)
	reAllowed := regexp.MustCompile(`(?i)data-testid="sandbox-launch-`)
	matches := reButton.FindAllString(content, -1)
	for _, btnTag := range matches {
		if !reAllowed.MatchString(btnTag) {
			violations = append(violations, SandboxParityViolation{
				Component: "Raw Button",
				File:      "ui/templates/sandbox.html",
				Message:   fmt.Sprintf("Raw <button> element detected: %q in sandbox.html. Use button_sharp template instead.", btnTag),
			})
		}
	}

	// 2. Raw listing card check
	reCard := regexp.MustCompile(`\blisting-card\b`)
	if reCard.MatchString(content) {
		violations = append(violations, SandboxParityViolation{
			Component: "Raw Listing Card",
			File:      "ui/templates/sandbox.html",
			Message:   "Raw listing-card class detected in sandbox.html. Use listing_card template instead.",
		})
	}

	return violations, nil
}

func findMissingComponents(definedMap map[string]string, referenced []string, hiddenTemplates map[string]bool) []SandboxParityViolation {
	refMap := make(map[string]bool)
	for _, r := range referenced {
		refMap[r] = true
	}

	var violations []SandboxParityViolation
	for d, file := range definedMap {
		if sandboxExcludeList[d] != "" {
			continue
		}

		if !refMap[d] {
			violations = append(violations, SandboxParityViolation{
				Component: d,
				File:      file,
				Message:   fmt.Sprintf("Component %q is defined but not documented in sandbox.html", d),
			})
			continue
		}

		if hiddenTemplates[d] {
			violations = append(violations, SandboxParityViolation{
				Component: d,
				File:      "ui/templates/sandbox.html",
				Message:   fmt.Sprintf("Component %q is referenced but inside a hidden element in sandbox.html (cannot be visually verified)", d),
			})
		}
	}
	return violations
}

func extractBlocks(path string, pattern string) ([]string, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	re := regexp.MustCompile(pattern)
	var names []string
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, match := range re.FindAllStringSubmatch(scanner.Text(), -1) {
			if len(match) > 1 && !seen[match[1]] {
				names = append(names, match[1])
				seen[match[1]] = true
			}
		}
	}

	return names, scanner.Err()
}

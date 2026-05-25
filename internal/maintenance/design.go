package maintenance

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// DesignViolation represents a deviation from the design system.
type DesignViolation struct {
	Content string
	Reason  string
	File    string
	Line    int
}

// LineFinder locates sequential tokens in the original content to report accurate line numbers.
type LineFinder struct {
	content string
	offset  int
}

// NewLineFinder creates a new LineFinder.
func NewLineFinder(content string) *LineFinder {
	return &LineFinder{content: content, offset: 0}
}

// FindLine returns the line number of a raw string inside the file.
func (lf *LineFinder) FindLine(tokenRaw string) int {
	idx := strings.Index(lf.content[lf.offset:], tokenRaw)
	if idx == -1 {
		return 1
	}

	absIdx := lf.offset + idx
	lf.offset = absIdx + len(tokenRaw)

	line := 1
	for i := 0; i < absIdx; i++ {
		if lf.content[i] == '\n' {
			line++
		}
	}
	return line
}

// CheckDesignStandards scans templates for violations of the UI Dialect protocol.
func CheckDesignStandards(dir string) ([]DesignViolation, error) {
	var violations []DesignViolation

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		v, err := checkFileStandards(path)
		if err != nil {
			return err
		}
		violations = append(violations, v...)

		uv, err := checkUppercaseDensity(path)
		if err != nil {
			return err
		}
		violations = append(violations, uv...)

		return nil
	})

	return violations, err
}

func checkFileStandards(path string) ([]DesignViolation, error) {
	// #nosec G304 -- verification utility scans local templates only
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)
	hasLabel := strings.Contains(content, "<label")

	roundedRegex := regexp.MustCompile(`\brounded-(sm|md|lg|xl|2xl|3xl)\b`)
	hexRegex := regexp.MustCompile(`(?i)#([0-9a-f]{3}|[0-9a-f]{6})\b`)
	adHocRegex := regexp.MustCompile(`(bg|text)-(white|black|stone-\d+)/\d+`)
	fontSizeRegex := regexp.MustCompile(`text-\[(\d+)`)
	opacityRegex := regexp.MustCompile(`text-text-sub/(\d+)`)

	var violations []DesignViolation
	violations = append(violations, tokenizeAndAudit(content, path, hasLabel, roundedRegex, hexRegex, adHocRegex, fontSizeRegex, opacityRegex)...)
	violations = append(violations, checkFileSpacers(content, path)...)
	violations = append(violations, checkFileCloseButtons(path)...)

	return violations, nil
}

func checkFileSpacers(content, path string) []DesignViolation {
	var violations []DesignViolation
	lines := strings.Split(content, "\n")
	for idx, line := range lines {
		violations = append(violations, checkDeadSpacers(path, idx+1, line)...)
	}
	return violations
}

func checkFileCloseButtons(path string) []DesignViolation {
	fv, err := checkRedundantCloseButtons(path)
	if err == nil {
		return fv
	}
	return nil
}

func tokenizeAndAudit(content, path string, hasLabel bool, roundedRegex, hexRegex, adHocRegex, fontSizeRegex, opacityRegex *regexp.Regexp) []DesignViolation {
	var violations []DesignViolation
	lf := NewLineFinder(content)
	reader := strings.NewReader(content)
	z := html.NewTokenizer(reader)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			token := z.Token()
			rawToken := z.Raw()
			lineNum := lf.FindLine(string(rawToken))
			tagStr := token.String()

			attrs := make(map[string]string)
			for _, attr := range token.Attr {
				attrs[attr.Key] = attr.Val
			}

			violations = append(violations, auditHTMLToken(token, tagStr, attrs, path, lineNum, hasLabel, roundedRegex, hexRegex, adHocRegex, fontSizeRegex, opacityRegex)...)
		}
	}
	return violations
}

func auditHTMLToken(token html.Token, tagStr string, attrs map[string]string, path string, lineNum int, hasLabel bool, roundedRegex, hexRegex, adHocRegex, fontSizeRegex, opacityRegex *regexp.Regexp) []DesignViolation {
	var violations []DesignViolation
	violations = append(violations, auditTailwindClasses(attrs, path, lineNum, tagStr, roundedRegex, fontSizeRegex, opacityRegex, adHocRegex)...)
	violations = append(violations, auditAttributesAndA11y(token, attrs, path, lineNum, tagStr, hasLabel, hexRegex)...)
	return violations
}

func auditTailwindClasses(attrs map[string]string, path string, lineNum int, tagStr string, roundedRegex, fontSizeRegex, opacityRegex, adHocRegex *regexp.Regexp) []DesignViolation {
	var violations []DesignViolation
	_, hasClass := attrs["class"]
	if !hasClass {
		return nil
	}

	violations = append(violations, checkRounding(path, lineNum, tagStr, roundedRegex)...)
	violations = append(violations, checkMinFontSize(path, lineNum, tagStr)...)
	violations = append(violations, checkLowContrastOpacity(path, lineNum, tagStr)...)
	violations = append(violations, checkHardcodedModalBg(path, lineNum, tagStr)...)
	violations = append(violations, checkAdHocColors(path, lineNum, tagStr)...)

	return violations
}

func auditAttributesAndA11y(token html.Token, attrs map[string]string, path string, lineNum int, tagStr string, hasLabel bool, hexRegex *regexp.Regexp) []DesignViolation {
	var violations []DesignViolation

	violations = append(violations, checkHexCodes(path, lineNum, tagStr, hexRegex)...)
	violations = append(violations, checkInlineHandlers(path, lineNum, tagStr)...)
	violations = append(violations, checkInlineStyles(path, lineNum, tagStr)...)
	violations = append(violations, checkMissingTestIDs(path, lineNum, tagStr)...)
	violations = append(violations, checkA11ySemantics(path, lineNum, tagStr, hasLabel)...)

	return violations
}

func checkRounding(path string, lineNum int, line string, re *regexp.Regexp) []DesignViolation {
	var v []DesignViolation
	if match := re.FindString(line); match != "" {
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  fmt.Sprintf("Forbidden rounding class '%s' (Brutalist standard requires sharp edges everywhere)", match),
		})
	}
	return v
}

func checkHexCodes(path string, lineNum int, line string, re *regexp.Regexp) []DesignViolation {
	var v []DesignViolation
	matches := re.FindAllStringIndex(line, -1)
	for _, matchIdx := range matches {
		start := matchIdx[0]
		// Skip HTML entities (&#123;).
		if start > 0 && line[start-1] == '&' {
			continue
		}
		match := line[matchIdx[0]:matchIdx[1]]
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  fmt.Sprintf("Hardcoded hex value '%s' (use Tailwind tokens instead)", match),
		})
	}
	return v
}

func checkMinFontSize(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	re := regexp.MustCompile(`text-\[(\d+)`)
	matches := re.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		if len(m) > 1 {
			var size int
			_, err := fmt.Sscanf(m[1], "%d", &size)
			if err == nil && size < 10 {
				v = append(v, DesignViolation{
					File:    path,
					Line:    lineNum,
					Content: line,
					Reason:  "Font size below 10px minimum (ADR: surface-theme-unification)",
				})
				break
			}
		}
	}
	return v
}

func checkLowContrastOpacity(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	re := regexp.MustCompile(`text-text-sub/(\d+)`)
	matches := re.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		if len(m) > 1 {
			var opacity int
			_, err := fmt.Sscanf(m[1], "%d", &opacity)
			if err == nil && opacity < 70 {
				v = append(v, DesignViolation{
					File:    path,
					Line:    lineNum,
					Content: line,
					Reason:  "Text opacity below 70% minimum for contrast",
				})
				break
			}
		}
	}
	return v
}

func checkHardcodedModalBg(path string, lineNum int, line string) []DesignViolation {
	base := filepath.Base(path)
	if !strings.HasPrefix(base, "modal_") && base != "ui_components.html" {
		return nil
	}

	if strings.Contains(strings.ReplaceAll(line, "dark:bg-earth-dark", ""), "bg-earth-dark") {
		return []DesignViolation{{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Hardcoded dark background bypasses light/dark theme sync (ADR: surface-theme-unification)",
		}}
	}
	return nil
}

func checkInlineHandlers(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	re := regexp.MustCompile(`\bon(click|change|submit|mouseover)\s*=`)
	if re.MatchString(line) {
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Inline event handler violates CSP policy (use data-* attributes + event delegation)",
		})
	}
	return v
}

func checkInlineStyles(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	if strings.Contains(line, "background-image:") {
		return nil
	}

	re := regexp.MustCompile(`\bstyle\s*=\s*"`)
	if re.MatchString(line) {
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Inline style attribute (use Tailwind classes)",
		})
	}
	return v
}

func checkUppercaseDensity(path string) ([]DesignViolation, error) {
	// #nosec G304 -- verification utility scans local templates only
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	count := 0
	re := regexp.MustCompile(`\buppercase\b`)

	for scanner.Scan() {
		if re.MatchString(scanner.Text()) {
			count++
		}
	}

	if count > 4 {
		return []DesignViolation{{
			File:   path,
			Line:   0, // File-level
			Reason: fmt.Sprintf("Uppercase density too high (%d occurrences, max 4 per template). Demote secondary elements to capitalize.", count),
		}}, scanner.Err()
	}

	return nil, scanner.Err()
}

func checkMissingTestIDs(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	// Target <button or <a that has hx-get|hx-post|hx-delete|hx-put
	if (strings.Contains(line, "<button") || strings.Contains(line, "<a ")) &&
		(strings.Contains(line, "hx-get") || strings.Contains(line, "hx-post") || strings.Contains(line, "hx-delete") || strings.Contains(line, "hx-put")) {
		if !strings.Contains(line, "data-testid=") {
			v = append(v, DesignViolation{
				File:    path,
				Line:    lineNum,
				Content: line,
				Reason:  "Interactive HTMX element missing data-testid (required for deterministic E2E testing)",
			})
		}
	}
	return v
}

func checkDeadSpacers(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	// Single-line check for empty section: <section...>\s*</section>
	re := regexp.MustCompile(`<section[^>]*>\s*</section>`)
	if re.MatchString(line) {
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Empty section element creates dead vertical space",
		})
	}
	return v
}

func checkRedundantCloseButtons(path string) ([]DesignViolation, error) {
	// Only relevant for modal templates, excluding components definitions
	if !strings.Contains(path, "modal_") || strings.Contains(path, "ui_components.html") {
		return nil, nil
	}

	// #nosec G304 -- verification utility scans local templates only
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	closeActions := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "data-modal-action=\"close\"") || strings.Contains(line, "template \"btn_close\"") {
			closeActions++
		}
	}

	if closeActions > 1 {
		return []DesignViolation{{
			File:   path,
			Reason: fmt.Sprintf("Multiple close actions detected (%d). Modals MUST have exactly one primary close mechanism to prevent UI clutter.", closeActions),
		}}, nil
	}
	return nil, nil
}

func checkAdHocColors(path string, lineNum int, line string) []DesignViolation {
	var v []DesignViolation
	// Catch things like bg-white/5 or hover:bg-black/5 which are often used for "fake" transparency that breaks themes
	re := regexp.MustCompile(`(bg|text)-(white|black|stone-\d+)/\d+`)
	matches := re.FindAllString(line, -1)
	for _, match := range matches {
		// Exceptions for backdrop or specific gradients if needed, but generally these should be themed tokens
		if strings.Contains(line, "backdrop:") || strings.Contains(line, "bg-gradient") {
			continue
		}
		v = append(v, DesignViolation{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  fmt.Sprintf("Ad-hoc color transparency '%s' detected. Use semantic tokens (e.g., bg-surface-dark/50) or theme-aware classes.", match),
		})
	}
	return v
}

func checkA11ySemantics(path string, lineNum int, line string, hasLabel bool) []DesignViolation {
	var v []DesignViolation

	v = append(v, checkIconButtonA11y(path, lineNum, line)...)
	v = append(v, checkLabelA11y(path, lineNum, line)...)
	v = append(v, checkInputA11y(path, lineNum, line, hasLabel)...)
	v = append(v, checkImgA11y(path, lineNum, line)...)

	return v
}

func checkIconButtonA11y(path string, lineNum int, line string) []DesignViolation {
	if strings.Contains(line, "<button") && strings.Contains(line, "class=") &&
		strings.Contains(line, "icon") && !strings.Contains(line, "aria-label=") {
		if regexp.MustCompile(`class="[^"]*icon[^"]*"`).MatchString(line) {
			return []DesignViolation{{
				File:    path,
				Line:    lineNum,
				Content: line,
				Reason:  "Icon-only button missing aria-label (required for screen readers)",
			}}
		}
	}
	return nil
}

func checkLabelA11y(path string, lineNum int, line string) []DesignViolation {
	if strings.Contains(line, "<label") && !strings.Contains(line, "for=") {
		return []DesignViolation{{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Label missing 'for' attribute to associate with input",
		}}
	}
	return nil
}

func checkInputA11y(path string, lineNum int, line string, hasLabel bool) []DesignViolation {
	if !hasLabel {
		return nil
	}
	if strings.Contains(line, "<input") || strings.Contains(line, "<textarea") || strings.Contains(line, "<select") {
		if !strings.Contains(line, "id=") {
			return []DesignViolation{{
				File:    path,
				Line:    lineNum,
				Content: line,
				Reason:  "Form input missing 'id' attribute (required when labels are present in file)",
			}}
		}
	}
	return nil
}

func checkImgA11y(path string, lineNum int, line string) []DesignViolation {
	if strings.Contains(line, "<img") && !strings.Contains(line, "alt=") {
		return []DesignViolation{{
			File:    path,
			Line:    lineNum,
			Content: line,
			Reason:  "Image missing 'alt' attribute (use alt=\"\" for decorative images)",
		}}
	}
	return nil
}

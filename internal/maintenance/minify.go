package maintenance

import (
	"os"
	"regexp"
	"strings"
)

var (
	commentRegex  = regexp.MustCompile(`(?s)<!--.*?-->`)
	newlinesRegex = regexp.MustCompile(`\n{2,}`)
)

// CompileAgentBundle reads the source files, strips HTML comments and excessive whitespace,
// and writes the result to dest.
func CompileAgentBundle(sources []string, dest string) error {
	var sb strings.Builder

	for _, src := range sources {
		content, err := os.ReadFile(src) //nolint:gosec // maintenance tool reads trusted core agent files
		if err != nil {
			return err
		}
		sb.Write(content)
		sb.WriteString("\n")
	}

	result := sb.String()

	// 1. Strip HTML comments
	result = commentRegex.ReplaceAllString(result, "")

	// 2. Reduce multiple consecutive newlines to a single newline
	result = newlinesRegex.ReplaceAllString(result, "\n")

	// 3. Trim leading/trailing whitespace
	result = strings.TrimSpace(result)

	return os.WriteFile(dest, []byte(result), 0600)
}

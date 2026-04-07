package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ExtractOpenAPIRoutes extracts HTTP methods and paths from docs/openapi.yaml
func ExtractOpenAPIRoutes(content []byte) ([]Route, error) {
	lines := strings.Split(string(content), "\n")
	var routes []Route
	var currentPath string

	pathRe := regexp.MustCompile(`^\s*'?(/.*?)'?:$`)
	methodRe := regexp.MustCompile(`(?i)^\s*(get|post|put|delete|patch|options|head):.*$`)

	for _, line := range lines {
		if matches := pathRe.FindStringSubmatch(line); len(matches) > 1 {
			currentPath = strings.TrimSpace(matches[1])
			continue
		}

		if matches := methodRe.FindStringSubmatch(line); len(matches) > 1 {
			method := strings.ToUpper(matches[1])
			if currentPath != "" {
				routes = append(routes, NewRoute(method, currentPath))
			}
		}
	}

	return uniqueAndSort(routes), nil
}

// ExtractMarkdownRoutes extracts HTTP methods and paths from docs/api.md
func ExtractMarkdownRoutes(content []byte) ([]Route, error) {
	lines := strings.Split(string(content), "\n")
	var routes []Route

	re := regexp.MustCompile(`(?i)^\|\s*(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD)\s*\|\s*` + "`?" + `([^` + "`" + `|\s]+)` + "`?" + `\s*\|`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			method := strings.ToUpper(matches[1])
			path := matches[2]
			routes = append(routes, NewRoute(method, path))
		}
	}

	return uniqueAndSort(routes), nil
}

// ExtractCLICodeCommands extracts CLI subcommands from Go source files.
func ExtractCLICodeCommands(dir string) ([]string, error) {
	var cmds []string
	regexes := []*regexp.Regexp{
		regexp.MustCompile(`(?m)Use:\s*"([^ "\n]+)`),
		regexp.MustCompile(`(?m)makeSimpleCmd\("([^"]+)"`),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		found, walkErr := extractCommandsFromCode(path, info, err, regexes)
		if walkErr != nil {
			return walkErr
		}
		cmds = append(cmds, found...)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return uniqueStrings(cmds), nil
}

func extractCommandsFromCode(path string, info os.FileInfo, err error, regexes []*regexp.Regexp) ([]string, error) {
	if err != nil || info.IsDir() || filepath.Ext(path) != ".go" {
		return nil, err
	}
 
	data, readErr := readFileOrErr(path, "code file")
	if readErr != nil {
		return nil, readErr
	}
 
	var found []string
	for _, re := range regexes {
		matches := re.FindAllStringSubmatch(string(data), -1)
		for _, match := range matches {
			if len(match) > 1 {
				cmd := strings.TrimSpace(match[1])
				if cmd != "" && cmd != "agbalumo" {
					found = append(found, cmd)
				}
			}
		}
	}
	return found, nil
}

// ExtractCLIMarkdownCommands extracts CLI commands from Markdown documentation.
func ExtractCLIMarkdownCommands(paths ...string) ([]string, error) {
	var cmds []string
	for _, path := range paths {
		found, err := extractCommandsFromPath(path)
		if err == nil {
			cmds = append(cmds, found...)
		}
	}
	return uniqueStrings(cmds), nil
}

func extractCommandsFromPath(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return walkMarkdownDir(path)
	}
	return extractFromMarkdownFile(path)
}

func walkMarkdownDir(dir string) ([]string, error) {
	var cmds []string
	err := filepath.Walk(dir, func(p string, i os.FileInfo, e error) error {
		if e == nil && !i.IsDir() && filepath.Ext(p) == ".md" {
			found, _ := extractFromMarkdownFile(p)
			cmds = append(cmds, found...)
		}
		return nil
	})
	return cmds, err
}

func extractFromMarkdownFile(path string) ([]string, error) {
	if filepath.Ext(path) != ".md" {
		return nil, nil
	}

	data, err := readFileOrErr(path, "markdown file")
	if err != nil {
		return nil, err
	}

	headerRe := regexp.MustCompile(`(?m)^###+\s+(.*)`)
	ignored := map[string]bool{
		"subcommands": true, "commands": true, "flags": true, "example": true,
		"quick reference": true, "environment variables": true, "global flags": true,
	}

	var cmds []string
	matches := headerRe.FindAllStringSubmatch(string(data), -1)
	for _, match := range matches {
		if len(match) > 1 {
			cmd := strings.TrimSpace(strings.ToLower(match[1]))
			if cmd != "" && !ignored[cmd] {
				cmds = append(cmds, cmd)
			}
		}
	}
	return cmds, nil
}

// CompareRoutes returns a list of differences between two route sets.
func CompareRoutes(sourceName, targetName string, source, target []Route) []string {
	targetMap := make(map[string]bool)
	for _, r := range target {
		key := fmt.Sprintf("%s %s", r.Method, r.Path)
		targetMap[key] = true
	}
	var missing []string
	for _, r := range source {
		key := fmt.Sprintf("%s %s", r.Method, r.Path)
		if !targetMap[key] {
			missing = append(missing, fmt.Sprintf("Missing in %s: %s (found in %s)", targetName, key, sourceName))
		}
	}
	return missing
}

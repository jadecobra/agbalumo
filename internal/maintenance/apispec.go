package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
				routes = append(routes, Route{
					Method: method,
					Path:   NormalizePath(currentPath),
				})
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
			routes = append(routes, Route{
				Method: method,
				Path:   NormalizePath(path),
			})
		}
	}

	return uniqueAndSort(routes), nil
}

// ExtractCLICodeCommands extracts CLI subcommands from Go source files.
func ExtractCLICodeCommands(dir string) ([]string, error) {
	var cmds []string
	useRe := regexp.MustCompile(`(?m)Use:\s*"([^ "\n]+)`)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
		// G304: Maintenance utility reads the OpenAPI source file
		data, err := os.ReadFile(path) //nolint:gosec // maintenance utility
			if err != nil {
				return err
			}
			matches := useRe.FindAllStringSubmatch(string(data), -1)
			for _, match := range matches {
				if len(match) > 1 {
					cmd := strings.TrimSpace(match[1])
					if cmd != "" && cmd != "agbalumo" {
						cmds = append(cmds, cmd)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return uniqueStrings(cmds), nil
}

// ExtractCLIMarkdownCommands extracts CLI commands from Markdown documentation.
func ExtractCLIMarkdownCommands(paths ...string) ([]string, error) {
	var cmds []string
	headerRe := regexp.MustCompile(`(?m)^###+\s+(.*)`)
	ignored := map[string]bool{
		"subcommands": true, "commands": true, "flags": true, "example": true,
		"quick reference": true, "environment variables": true, "global flags": true,
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.IsDir() {
			walkErr := filepath.Walk(path, func(p string, i os.FileInfo, e error) error {
				if e == nil && !i.IsDir() && filepath.Ext(p) == ".md" {
					// G304: Maintenance utility reads partial openapi files
					data, _ := os.ReadFile(p) //nolint:gosec // maintenance utility
					matches := headerRe.FindAllStringSubmatch(string(data), -1)
					for _, match := range matches {
						if len(match) > 1 {
							cmd := strings.TrimSpace(strings.ToLower(match[1]))
							if cmd != "" && !ignored[cmd] {
								cmds = append(cmds, cmd)
							}
						}
					}
				}
				return nil
			})
			if walkErr != nil {
				// We don't fail the whole command if one MD file fails to read
				fmt.Fprintf(os.Stderr, "Warning: failed to walk %s: %v\n", path, walkErr)
			}
		} else if filepath.Ext(path) == ".md" {
		// G304: Maintenance utility reads API markdown file
		data, err := os.ReadFile(path) //nolint:gosec // maintenance utility
			if err == nil {
				matches := headerRe.FindAllStringSubmatch(string(data), -1)
				for _, match := range matches {
					if len(match) > 1 {
						cmd := strings.TrimSpace(strings.ToLower(match[1]))
						if cmd != "" && !ignored[cmd] {
							cmds = append(cmds, cmd)
						}
					}
				}
			}
		}
	}

	return uniqueStrings(cmds), nil
}

func uniqueAndSort(routes []Route) []Route {
	seen := make(map[string]bool)
	var unique []Route
	for _, r := range routes {
		key := r.Method + " " + r.Path
		if !seen[key] {
			seen[key] = true
			unique = append(unique, r)
		}
	}
	sort.Slice(unique, func(i, j int) bool {
		if unique[i].Path == unique[j].Path {
			return unique[i].Method < unique[j].Method
		}
		return unique[i].Path < unique[j].Path
	})
	return unique
}

func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	var unique []string
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			unique = append(unique, s)
		}
	}
	sort.Strings(unique)
	return unique
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

package agent

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

	routes = uniqueAndSort(routes)
	return routes, nil
}

// ExtractMarkdownRoutes extracts HTTP methods and paths from docs/api.md
func ExtractMarkdownRoutes(content []byte) ([]Route, error) {
	lines := strings.Split(string(content), "\n")
	var routes []Route

	// Regex to match markdown table rows: `| GET | /path | ... |` or `| GET | `+/path+` | ... |`
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

	routes = uniqueAndSort(routes)
	return routes, nil
}

// CompareRoutes finds routes present in source but missing in target
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
			missing = append(missing, fmt.Sprintf("❌ Missing in %s: %s (found in %s)", targetName, key, sourceName))
		}
	}

	return missing
}

// CheckAPIDrift orchestrates the comparisons and returns all detected drift
func CheckAPIDrift(codeRoutes, openapiRoutes, mdRoutes []Route) []string {
	var allDrift []string

	// Compare Code vs OpenAPI
	allDrift = append(allDrift, CompareRoutes("Code (cmd/server.go)", "OpenAPI (docs/openapi.yaml)", codeRoutes, openapiRoutes)...)
	allDrift = append(allDrift, CompareRoutes("OpenAPI (docs/openapi.yaml)", "Code (cmd/server.go)", openapiRoutes, codeRoutes)...)

	// Compare Code vs MD
	allDrift = append(allDrift, CompareRoutes("Code (cmd/server.go)", "API Docs (docs/api.md)", codeRoutes, mdRoutes)...)
	allDrift = append(allDrift, CompareRoutes("API Docs (docs/api.md)", "Code (cmd/server.go)", mdRoutes, codeRoutes)...)

	// Compare OpenAPI vs MD
	allDrift = append(allDrift, CompareRoutes("OpenAPI (docs/openapi.yaml)", "API Docs (docs/api.md)", openapiRoutes, mdRoutes)...)
	allDrift = append(allDrift, CompareRoutes("API Docs (docs/api.md)", "OpenAPI (docs/openapi.yaml)", mdRoutes, openapiRoutes)...)

	return uniqueStrings(allDrift)
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
	return unique
}

func ExtractCLICodeCommands(dir string) ([]string, error) {
	var cmds []string
	useRe := regexp.MustCompile(`(?m)Use:\s*"([^ "\n]+)`)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// #nosec G304 G122 - Internal harness tool reading project files
			data, err := os.ReadFile(path)
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
	cmds = uniqueStrings(cmds)
	sort.Strings(cmds)
	return cmds, nil
}

func ExtractCLIMarkdownCommands(paths ...string) ([]string, error) {
	var cmds []string
	headerRe := regexp.MustCompile(`(?m)^###+\s+(.*)`)
	ignored := map[string]bool{
		"subcommands":               true,
		"commands":                  true,
		"flags":                     true,
		"example":                   true,
		"quick reference":           true,
		"environment variables":     true,
		"global flags":              true,
		"agent harness and testing": true,
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			continue // skip missing files like docs/cli/*.md if dir not there
		}
		if info.IsDir() {
			_ = filepath.Walk(path, func(p string, i os.FileInfo, e error) error {
				if e == nil && !i.IsDir() && filepath.Ext(p) == ".md" {
					// #nosec G304 G122 - Internal harness tool reading project files
					data, _ := os.ReadFile(p)
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
		} else if filepath.Ext(path) == ".md" {
			// #nosec G304 - Internal harness tool reading project files
			data, err := os.ReadFile(path)
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

	cmds = uniqueStrings(cmds)
	sort.Strings(cmds)
	return cmds, nil
}

func CheckCLIDrift(codeCmds, mdCmds []string) []string {
	var diffs []string
	codeMap := make(map[string]bool)
	for _, c := range codeCmds {
		codeMap[c] = true
	}
	mdMap := make(map[string]bool)
	for _, c := range mdCmds {
		mdMap[c] = true
	}

	for _, c := range codeCmds {
		if !mdMap[c] {
			diffs = append(diffs, fmt.Sprintf("❌ Missing in CLI Docs: %s (found in Code)", c))
		}
	}
	for _, c := range mdCmds {
		if !codeMap[c] {
			diffs = append(diffs, fmt.Sprintf("❌ Missing in Code: %s (found in CLI Docs)", c))
		}
	}

	sort.Strings(diffs)
	return diffs
}

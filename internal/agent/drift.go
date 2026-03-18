package agent

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ExtractOpenAPIRoutes extracts HTTP methods and paths from docs/openapi.yaml
func ExtractOpenAPIRoutes(content []byte) ([]Route, error) {
	lines := strings.Split(string(content), "\n")
	var routes []Route
	var currentPath string

	pathRe := regexp.MustCompile(`^\s*(/.*?):$`)
	methodRe := regexp.MustCompile(`(?i)^\s*(get|post|put|delete|patch|options|head):$`)

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
					Path:   normalizePath(currentPath),
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
				Path:   normalizePath(path),
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

func normalizePath(p string) string {
	// Our bash logic:
	// 1. replace :id with {id}
	p = regexp.MustCompile(`:([a-zA-Z0-9_]+)`).ReplaceAllString(p, "{$1}")
	
	// 2. Remove trailing slashes (except root)
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}
	
	// 3. Deduplicate slashes
	p = regexp.MustCompile(`//+`).ReplaceAllString(p, "/")
	
	if p == "" {
		p = "/"
	}
	return p
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

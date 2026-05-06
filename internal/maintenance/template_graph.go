package maintenance

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type TemplateNode struct { //nolint:govet // maintenance utility
	Name       string
	DefinedIn  string
	CalledBy   []TemplateCall
	Calls      []string
	References []string
}

type TemplateCall struct { //nolint:govet // maintenance utility
	File     string
	DictKeys []string
	Line     int
}

var (
	defineRegex   = regexp.MustCompile(`{{\s*define\s*"([^"]+)"\s*}}`)
	templateRegex = regexp.MustCompile(`{{\s*template\s*"([^"]+)"\s*(?:dict\s+(.+?))?\s*}}`)
	refRegex      = regexp.MustCompile(`(?:\$|)\.(\w+)`)
	keyRegex      = regexp.MustCompile(`"(\w+)"`)
	actionRegex   = regexp.MustCompile(`{{(.*?)}}`)
)

func BuildTemplateGraph(dir string) (map[string]*TemplateNode, error) {
	graph := make(map[string]*TemplateNode)

	if err := scanDefinitions(dir, graph); err != nil {
		return nil, err
	}

	if err := scanInvocations(dir, graph); err != nil {
		return nil, err
	}

	return graph, nil
}

func scanDefinitions(dir string, graph map[string]*TemplateNode) error {
	// #nosec G122 -- local maintenance utility
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		// #nosec G304,G122 -- local templates only
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		sContent := string(content)
		matches := defineRegex.FindAllStringSubmatchIndex(sContent, -1)

		for i, match := range matches {
			processDefinition(path, sContent, matches, i, match, graph)
		}
		return nil
	})
}

func processDefinition(path, sContent string, matches [][]int, i int, match []int, graph map[string]*TemplateNode) {
	name := sContent[match[2]:match[3]]
	end := len(sContent)
	if i+1 < len(matches) {
		end = matches[i+1][0]
	}
	block := sContent[match[0]:end]

	node, ok := graph[name]
	if !ok {
		node = &TemplateNode{Name: name}
		graph[name] = node
	}
	node.DefinedIn = path
	node.References = extractNodeReferences(block)
	node.Calls = extractNodeCalls(block)
}

func extractNodeReferences(block string) []string {
	refs := make(map[string]bool)
	// First, find all {{ ... }} blocks
	actionMatches := actionRegex.FindAllStringSubmatch(block, -1)
	for _, am := range actionMatches {
		actionContent := am[1]
		// Then, find references only within the action content
		refMatches := refRegex.FindAllStringSubmatch(actionContent, -1)
		for _, rm := range refMatches {
			refs[rm[1]] = true
		}
	}
	var sortedRefs []string
	for r := range refs {
		sortedRefs = append(sortedRefs, r)
	}
	sort.Strings(sortedRefs)
	return sortedRefs
}

func extractNodeCalls(block string) []string {
	var calls []string
	callMatches := templateRegex.FindAllStringSubmatch(block, -1)
	for _, cm := range callMatches {
		calls = append(calls, cm[1])
	}
	return calls
}

func scanInvocations(dir string, graph map[string]*TemplateNode) error {
	// #nosec G122 -- local maintenance utility
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		// #nosec G304,G122 -- local templates only
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")
		for lineIdx, line := range lines {
			processInvocationLine(path, lineIdx+1, line, graph)
		}
		return nil
	})
}

func processInvocationLine(path string, lineNum int, line string, graph map[string]*TemplateNode) {
	matches := templateRegex.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		tmplName := m[1]
		dictContent := ""
		if len(m) > 2 {
			dictContent = m[2]
		}

		node, ok := graph[tmplName]
		if !ok {
			node = &TemplateNode{Name: tmplName}
			graph[tmplName] = node
		}

		call := TemplateCall{
			File:     path,
			Line:     lineNum,
			DictKeys: extractDictKeys(dictContent),
		}
		node.CalledBy = append(node.CalledBy, call)
	}
}

func extractDictKeys(dictContent string) []string {
	if dictContent == "" {
		return nil
	}
	keyMatches := keyRegex.FindAllStringSubmatch(dictContent, -1)
	keys := make(map[string]bool)
	for _, km := range keyMatches {
		keys[km[1]] = true
	}
	var sortedKeys []string
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

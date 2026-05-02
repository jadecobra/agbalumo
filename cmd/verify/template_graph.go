package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var templateGraphCmd = &cobra.Command{
	Use:   "template-graph [name]",
	Short: "Show dependency tree for a named template",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		graph, err := maintenance.BuildTemplateGraph("ui/templates")
		if err != nil {
			return fmt.Errorf("failed to build template graph: %w", err)
		}

		if len(args) == 1 {
			name := args[0]
			node, ok := graph[name]
			if !ok {
				return fmt.Errorf("template '%s' not found", name)
			}
			printTree(graph, node, 0, make(map[string]bool))
			return nil
		}

		// Print summary
		fmt.Printf("%-30s %-10s %-10s %-30s\n", "TEMPLATE", "CALLERS", "CALLS", "FILE")
		fmt.Println(strings.Repeat("-", 80))

		var names []string
		for name := range graph {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			node := graph[name]
			fmt.Printf("%-30s %-10d %-10d %-30s\n",
				name,
				len(node.CalledBy),
				len(node.Calls),
				node.DefinedIn,
			)
		}

		return nil
	},
}

func printTree(graph map[string]*maintenance.TemplateNode, node *maintenance.TemplateNode, indent int, visited map[string]bool) {
	prefix := strings.Repeat("  ", indent)
	if indent > 0 {
		prefix += "└─ "
	}

	fmt.Printf("%s%s (%s)\n", prefix, node.Name, node.DefinedIn)

	if len(node.References) > 0 {
		fmt.Printf("%s   Refs: %s\n", strings.Repeat("  ", indent), strings.Join(node.References, ", "))
	}

	if visited[node.Name] {
		fmt.Printf("%s   (circular dependency)\n", strings.Repeat("  ", indent))
		return
	}
	visited[node.Name] = true

	for _, call := range node.Calls {
		if next, ok := graph[call]; ok {
			printTree(graph, next, indent+1, visited)
		} else {
			fmt.Printf("%s   └─ %s (EXTERNAL/UNKNOWN)\n", strings.Repeat("  ", indent+1), call)
		}
	}
}

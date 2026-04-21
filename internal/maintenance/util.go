package maintenance

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/jadecobra/agbalumo/internal/util"
)

func readFileOrErr(path, label string) ([]byte, error) {
	data, err := os.ReadFile(path) //nolint:gosec // maintenance utility
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", label, err)
	}
	return data, nil
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
	return util.UniqueStrings(strs)
}

func runTool(dir, name string, args ...string) (string, error) {
	//nolint:gosec // G204: Maintenance utility running trusted audit tools
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

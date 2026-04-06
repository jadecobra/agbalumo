package maintenance

import (
	"sort"
)

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

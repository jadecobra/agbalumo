package util

import "sort"

// UniqueStrings returns a sorted slice of unique strings.
func UniqueStrings(input []string) []string {
	if len(input) == 0 {
		return input
	}
	m := make(map[string]bool)
	var result []string
	for _, s := range input {
		if !m[s] {
			m[s] = true
			result = append(result, s)
		}
	}
	sort.Strings(result)
	return result
}

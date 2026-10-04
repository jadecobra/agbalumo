package maintenance

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/uudashr/gocognit"
)

// ComplexityProfileOptions configures cognitive complexity profiling.
type ComplexityProfileOptions struct {
	Threshold   int
	Top         int
	Diagnostics bool
}

// ProfileComplexity scans rootDir for Go source files and computes cognitive complexity per function.
func ProfileComplexity(rootDir string, opts ComplexityProfileOptions) ([]gocognit.Stat, error) {
	var allStats []gocognit.Stat

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if shouldSkipComplexityDir(path, d) {
			return filepath.SkipDir
		}
		if !isCandidateGoFile(d) {
			return nil
		}

		fileStats := analyzeGoFile(path, opts)
		allStats = appendMatchingStats(allStats, fileStats, opts.Threshold)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory for complexity: %w", err)
	}

	sortStats(allStats)
	return limitStats(allStats, opts.Top), nil
}

func shouldSkipComplexityDir(path string, d fs.DirEntry) bool {
	if !d.IsDir() {
		return false
	}
	name := d.Name()
	if name == "." {
		return false
	}
	return strings.HasPrefix(name, ".") || name == vVendor || name == vNodeModules
}

func isCandidateGoFile(d fs.DirEntry) bool {
	return !d.IsDir() && strings.HasSuffix(d.Name(), ".go")
}

func analyzeGoFile(path string, opts ComplexityProfileOptions) []gocognit.Stat {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil
	}

	if opts.Diagnostics {
		return gocognit.ComplexityStatsWithDiagnostic(node, fset, nil, true)
	}
	return gocognit.ComplexityStats(node, fset, nil)
}

func appendMatchingStats(all []gocognit.Stat, incoming []gocognit.Stat, threshold int) []gocognit.Stat {
	for _, stat := range incoming {
		if stat.Complexity >= threshold {
			all = append(all, stat)
		}
	}
	return all
}

func sortStats(stats []gocognit.Stat) {
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Complexity != stats[j].Complexity {
			return stats[i].Complexity > stats[j].Complexity
		}
		return stats[i].FuncName < stats[j].FuncName
	})
}

func limitStats(stats []gocognit.Stat, top int) []gocognit.Stat {
	if top > 0 && len(stats) > top {
		return stats[:top]
	}
	return stats
}

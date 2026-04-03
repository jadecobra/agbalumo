package agent

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

// claudeSonnetWindow is the binding worst-case context window.
// Satisfying this automatically satisfies larger windows (e.g. Gemini Flash 1M).
const claudeSonnetWindow = 200_000

type FileCost struct {
	FilePath string
	Lines    int
	Tokens   int // estimated via cl100k_base tokenizer
}

type CostReport struct {
	TotalFiles       int
	TotalLines       int
	TotalTokens      int     // NEW: sum of token counts across all files
	RMS              float64 // existing: LOC-based RMS (kept for parallel transition)
	TokenRMS         float64 // NEW: token-based RMS
	ContextWindowPct float64 // NEW: TotalTokens / claudeSonnetWindow * 100
	TopFiles         []FileCost
}

func CalculateContextCost(dir string) (*CostReport, error) {
	var fileCosts []FileCost

	ignoredDirs := map[string]bool{
		".git":         true,
		"vendor":       true,
		"node_modules": true,
		"dist":         true,
		"build":        true,
		".tester":      true,
		".agents":      true,
		".agent":       true,
		"scripts":      true,
	}

	validExts := map[string]bool{
		".go":   true,
		".html": true,
		".css":  true,
		".js":   true,
		".json": true,
		".md":   true,
		".sh":   true,
		".yml":  true,
		".yaml": true,
		".sql":  true,
	}

	ignoredFiles := map[string]bool{
		"package-lock.json": true,
		"pnpm-lock.yaml":    true,
		"yarn.lock":         true,
		"go.sum":            true,
		"go.mod":            true,
		"Taskfile.yml":      true,
	}

	ignoredPathSegments := []string{
		"ui/static",
	}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if ignoredDirs[d.Name()] {
				return filepath.SkipDir
			}

			// Handle multi-segment paths (e.g., "ui/static")
			slashPath := "/" + filepath.ToSlash(path)
			for _, seg := range ignoredPathSegments {
				if strings.Contains(slashPath, "/"+seg+"/") || strings.HasSuffix(slashPath, "/"+seg) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if ignoredFiles[d.Name()] {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !validExts[ext] {
			return nil
		}

		// #nosec G304 G122 - Internal harness tool reading project files
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		lines := bytes.Count(content, []byte{'\n'})
		if len(content) > 0 && content[len(content)-1] != '\n' {
			lines++
		}

		// Token counting — cl100k_base is a close approximation for Gemini/Claude (~5% error)
		enc, encErr := tiktoken.GetEncoding("cl100k_base")
		tokenCount := 0
		if encErr == nil {
			tokenCount = len(enc.Encode(string(content), nil, nil))
		} else {
			// Fallback: ~4 chars per token is a widely-used heuristic
			tokenCount = len(content) / 4
		}

		fileCosts = append(fileCosts, FileCost{
			FilePath: path,
			Lines:    lines,
			Tokens:   tokenCount,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	report := &CostReport{}
	if len(fileCosts) == 0 {
		return report, nil
	}

	report.TotalFiles = len(fileCosts)

	var sumSquares float64
	for _, fc := range fileCosts {
		report.TotalLines += fc.Lines
		sumSquares += float64(fc.Lines) * float64(fc.Lines)
	}

	meanSquares := sumSquares / float64(len(fileCosts))
	report.RMS = math.Sqrt(meanSquares)

	var tokenSumSquares float64
	for _, fc := range fileCosts {
		report.TotalTokens += fc.Tokens
		tokenSumSquares += float64(fc.Tokens) * float64(fc.Tokens)
	}
	tokenMeanSquares := tokenSumSquares / float64(len(fileCosts))
	report.TokenRMS = math.Sqrt(tokenMeanSquares)
	report.ContextWindowPct = float64(report.TotalTokens) / claudeSonnetWindow * 100

	sort.SliceStable(fileCosts, func(i, j int) bool {
		return fileCosts[i].Tokens > fileCosts[j].Tokens
	})

	report.TopFiles = fileCosts
	if len(report.TopFiles) > 10 {
		report.TopFiles = report.TopFiles[:10]
	}

	return report, nil
}
// Fast-path verification final attempt

package maintenance

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

const claudeSonnetWindow = 200_000

type FileCost struct {
	FilePath string
	Lines    int
	Tokens   int
}

type CostReport struct {
	TopFiles         []FileCost
	TotalFiles       int
	TotalLines       int
	TotalTokens      int
	RMS              float64
	TokenRMS         float64
	ContextWindowPct float64
}

var ignoredDirs = map[string]bool{
	".git": true, "vendor": true, "node_modules": true, "dist": true, "build": true,
	".tester": true, ".agents": true, ".agent": true, "scripts": true,
}

var validExts = map[string]bool{
	".go": true, ".html": true, ".css": true, ".js": true, ".json": true,
	".md": true, ".sh": true, ".yml": true, ".yaml": true, ".sql": true,
}

var ignoredFiles = map[string]bool{
	"package-lock.json": true, "pnpm-lock.yaml": true, "yarn.lock": true,
	"go.sum": true, "go.mod": true,
}

type costCollector struct {
	costs []FileCost
}

func (c *costCollector) walk(path string, d os.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
		if ignoredDirs[d.Name()] {
			return filepath.SkipDir
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

	cost, costErr := calculateFileCost(path)
	if costErr == nil {
		c.costs = append(c.costs, cost)
	}
	return nil
}

func CalculateContextCost(dir string) (*CostReport, error) {
	collector := &costCollector{}
	err := filepath.WalkDir(dir, collector.walk)
	if err != nil {
		return nil, err
	}

	return generateReport(collector.costs), nil
}

func calculateFileCost(path string) (FileCost, error) {
	// G304: Maintenance utility reads source files for token counting
	content, err := os.ReadFile(path) //nolint:gosec // maintenance utility
	if err != nil {
		return FileCost{}, err
	}

	lines := bytes.Count(content, []byte{'\n'})
	if len(content) > 0 && content[len(content)-1] != '\n' {
		lines++
	}

	enc, encErr := tiktoken.GetEncoding("cl100k_base")
	tokenCount := 0
	if encErr == nil {
		tokenCount = len(enc.Encode(string(content), nil, nil))
	} else {
		tokenCount = len(content) / 4
	}

	return FileCost{
		FilePath: path,
		Lines:    lines,
		Tokens:   tokenCount,
	}, nil
}

func generateReport(fileCosts []FileCost) *CostReport {
	report := &CostReport{}
	if len(fileCosts) == 0 {
		return report
	}

	report.TotalFiles = len(fileCosts)
	var sumSquares, tokenSumSquares float64
	for _, fc := range fileCosts {
		report.TotalLines += fc.Lines
		report.TotalTokens += fc.Tokens
		sumSquares += float64(fc.Lines) * float64(fc.Lines)
		tokenSumSquares += float64(fc.Tokens) * float64(fc.Tokens)
	}

	report.RMS = math.Sqrt(sumSquares / float64(len(fileCosts)))
	report.TokenRMS = math.Sqrt(tokenSumSquares / float64(len(fileCosts)))
	report.ContextWindowPct = float64(report.TotalTokens) / claudeSonnetWindow * 100

	sort.SliceStable(fileCosts, func(i, j int) bool {
		return fileCosts[i].Tokens > fileCosts[j].Tokens
	})

	report.TopFiles = fileCosts
	if len(report.TopFiles) > 10 {
		report.TopFiles = report.TopFiles[:10]
	}

	return report
}

package agent

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileCost struct {
	FilePath string
	Lines    int
}

type CostReport struct {
	TotalFiles int
	TotalLines int
	RMS        float64
	TopFiles   []FileCost
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
	}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		lines := bytes.Count(content, []byte{'\n'})
		if len(content) > 0 && content[len(content)-1] != '\n' {
			lines++
		}

		fileCosts = append(fileCosts, FileCost{
			FilePath: path,
			Lines:    lines,
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

	sort.SliceStable(fileCosts, func(i, j int) bool {
		return fileCosts[i].Lines > fileCosts[j].Lines
	})

	report.TopFiles = fileCosts
	if len(report.TopFiles) > 10 {
		report.TopFiles = report.TopFiles[:10]
	}

	return report, nil
}

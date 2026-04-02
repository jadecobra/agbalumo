package agent

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jadecobra/agbalumo/internal/util"
)

func VerifyCoverage() bool {
	fmt.Println("Verifying test coverage...")

	covFile := filepath.Join(".tester", "coverage", "coverage.out")

	if _, statErr := util.SafeStat(covFile); util.SafeIsNotExist(statErr) {
		fmt.Println("❌ Gate FAIL: coverage profile not generated.")
		return false
	}

	// #nosec G304 - Internal harness tool reading coverage profile
	f, err := util.SafeOpen(covFile)
	if err != nil {
		fmt.Println("❌ Gate FAIL: unable to read coverage profile.")
		return false
	}
	defer func() { _ = f.Close() }()

	coverage, err := ParseCoverageProfile(f)
	if err != nil {
		fmt.Println("❌ Gate FAIL: unable to parse coverage profile.")
		return false
	}

	// #nosec G304 - Internal harness tool reading thresholds
	thresholdsData, err := util.SafeReadFile(filepath.Join(".agents", "coverage-thresholds.json"))
	var thresholds map[string]float64
	if err == nil {
		var parseErr error
		thresholds, parseErr = ParseThresholds(thresholdsData)
		if parseErr != nil {
			fmt.Println("❌ " + parseErr.Error())
			return false
		}
	} else {
		globalThreshold := 90.0
		// #nosec G304 - Internal harness tool reading threshold
		legacyData, err := util.SafeReadFile(filepath.Join(".agents", "coverage-threshold"))
		if err == nil {
			parsed, err := strconv.ParseFloat(strings.TrimSpace(string(legacyData)), 64)
			if err == nil {
				globalThreshold = parsed
			}
		}
		thresholds = map[string]float64{"default": globalThreshold}
	}

	violations := EnforceCoverage(coverage, thresholds)

	if len(violations) > 0 {
		// #nosec G204 - Internal harness tool executing cover tool
		out, _ := exec.Command("go", "tool", "cover", "-func="+covFile).CombinedOutput()
		totalLine := ""
		for _, line := range strings.Split(string(out), "\n") {
			if strings.HasPrefix(line, "total:") {
				totalLine = line
			}
		}

		fmt.Printf("❌ Gate FAIL: %s. Thresholds not met.\n", totalLine)
		if len(violations) > 0 {
			fmt.Println("  " + violations[0])
			if len(violations) > 1 {
				fmt.Printf("  ... and %d more violations.\n", len(violations)-1)
			}
		}
		return false
	}

	fmt.Printf("✅ Gate PASS: %s meets thresholds.\n", GateCoverage)
	return true
}

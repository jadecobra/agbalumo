package maintenance

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const errNotFoundFmt = "could not find %s: %w"

// RunPerformanceAudit executes the performance audit checks.
func RunPerformanceAudit(rootDir string) error {
	fmt.Println("🚀  Starting Performance Audit...")
	fmt.Println("---------------------------------")

	checks := getPerfChecks()
	errs := runPerfChecks(rootDir, checks)

	fmt.Println("---------------------------------")
	return reportAuditErrors(errs)
}

type perfCheck struct {
	fn   func(string) error
	name string
}

func getPerfChecks() []perfCheck {
	return []perfCheck{
		{checkFileSizes, "Static Asset Sizes"},
		{checkSQLitePragmas, "SQLite Configuration"},
		{runSearchBenchmark, "Search Smoke Benchmark"},
		{runBulkInsertBenchmark, "Bulk Insert Benchmark (10k items)"},
	}
}

func runPerfChecks(rootDir string, checks []perfCheck) []string {
	var errs []string
	isCI := os.Getenv("GITHUB_ACTIONS") == "true"

	for _, c := range checks {
		if shouldSkipCheck(isCI, c.name) {
			fmt.Printf("[?] Skipping %s (CI environment detected)... ✅\n", c.name)
			continue
		}

		fmt.Printf("[?] Checking %s... ", c.name)
		if err := c.fn(rootDir); err != nil {
			fmt.Println("❌ Failed")
			errs = append(errs, err.Error())
		} else {
			fmt.Println("✅ Passed")
		}
	}
	return errs
}

func shouldSkipCheck(isCI bool, name string) bool {
	return isCI && strings.Contains(name, "Benchmark")
}

func reportAuditErrors(errs []string) error {
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("❌ %s\n", e)
		}
		return fmt.Errorf("performance audit failed with %d issues", len(errs))
	}

	fmt.Println("✅ Performance Audit Passed Successfully!")
	return nil
}


func checkFileSizes(rootDir string) error {
	cssPath := filepath.Join(rootDir, "ui/static/css/output.css")
	jsPath := filepath.Join(rootDir, "ui/static/js/app.js")

	cssStat, err := os.Stat(cssPath)
	if err != nil {
		return fmt.Errorf(errNotFoundFmt, cssPath, err)
	}
	if cssStat.Size() > 150*1024 {
		return fmt.Errorf("output.css is too large: %d bytes (limit: 150KB)", cssStat.Size())
	}

	jsStat, err := os.Stat(jsPath)
	if err != nil {
		return fmt.Errorf(errNotFoundFmt, jsPath, err)
	}
	if jsStat.Size() > 50*1024 {
		return fmt.Errorf("app.js is too large: %d bytes (limit: 50KB)", jsStat.Size())
	}

	return nil
}

func checkSQLitePragmas(rootDir string) error {
	sqlitePath := filepath.Join(rootDir, "internal/repository/sqlite/sqlite.go")
	data, err := readFileOrErr(sqlitePath, "SQLite repository file")
	if err != nil {
		return err
	}

	content := string(data)
	checks := []string{"journal_mode=WAL", "busy_timeout", "MaxOpenConns"}
	for _, c := range checks {
		if !strings.Contains(content, c) {
			return fmt.Errorf("missing critical SQLite pragma/config: %s", c)
		}
	}
	return nil
}

func runSearchBenchmark(rootDir string) error {
	testFile := filepath.Join(rootDir, "internal/repository/sqlite/search_performance_test.go")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		// Skip if test file doesn't exist (like in unit tests)
		return nil
	}

	cmd := exec.Command("go", "test", "-bench=BenchmarkSearchPerformance/FindAll_Default_Page1", "-benchtime=100ms", "./internal/repository/sqlite/search_performance_test.go")
	cmd.Dir = rootDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("search benchmark failed: %w\nOutput: %s", err, string(out))
	}
	return nil
}
func runBulkInsertBenchmark(rootDir string) error {
	benchFile := filepath.Join(rootDir, "internal/repository/sqlite/sqlite_listing_bench_test.go")
	if _, err := os.Stat(benchFile); os.IsNotExist(err) {
		return nil
	}

	benchName := "BenchmarkSQLiteRepository_BulkInsertListings"
	// We run with -benchtime=1x because the benchmark itself does 10,000 items,
	// which is enough for a performance verification sample.
	cmd := exec.Command("go", "test", "-v", "-bench="+benchName, "-run=^#", "-benchtime=1x", "./internal/repository/sqlite")
	cmd.Dir = rootDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("bulk insert benchmark failed: %w\nOutput: %s", err, string(out))
	}
	return nil
}

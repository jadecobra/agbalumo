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

	errs := []string{}
	const failMsg = "❌ Failed"

	// 1. Static Asset Sizes
	fmt.Print("[?] Checking Static Asset Sizes... ")
	if err := checkFileSizes(rootDir); err != nil {
		fmt.Println("❌ Failed")
		errs = append(errs, err.Error())
	} else {
		fmt.Println("✅ Passed")
	}

	// 2. SQLite Configuration
	fmt.Print("[?] Checking SQLite Pragmas... ")
	if err := checkSQLitePragmas(rootDir); err != nil {
		fmt.Println("❌ Failed")
		errs = append(errs, err.Error())
	} else {
		fmt.Println("✅ Passed")
	}

	// 3. Search Smoke Benchmark
	fmt.Print("[?] Running Search Smoke Benchmark... ")
	if err := runSearchBenchmark(rootDir); err != nil {
		fmt.Println(failMsg)
		errs = append(errs, err.Error())
	} else {
		fmt.Println("✅ Passed")
	}

	fmt.Println("---------------------------------")
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

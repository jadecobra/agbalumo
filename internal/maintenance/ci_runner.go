package maintenance

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// CITask represents a single step in the CI pipeline.
type CITask struct {
	Fn   func() error
	Name string
}

// RunParallelCI executes the CI pipeline using dynamic concurrency based on system resources.
func RunParallelCI(ctx context.Context, tasks []CITask) error {
	start := time.Now()
	fmt.Printf("🚀 Starting Parallel CI Pipeline (CPUs: %d)\n", runtime.NumCPU())

	// Use an errgroup with a limited number of concurrent workers.
	// We allow high concurrency for light tasks, but for CI we'll limit to NumCPU.
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(runtime.NumCPU())

	var mu sync.Mutex
	completed := 0
	total := len(tasks)

	for _, task := range tasks {
		t := task // capture range variable
		g.Go(func() error {
			innerStart := time.Now()
			fmt.Printf("\n[RUN] %s\n", t.Name)

			err := t.Fn()

			mu.Lock()
			completed++
			pct := (float64(completed) / float64(total)) * 100
			status := "✅"
			if err != nil {
				status = "❌"
			}
			fmt.Printf("\n[%s] %s (%.2fs) [%d/%d - %.0f%%]\n", status, t.Name, time.Since(innerStart).Seconds(), completed, total, pct)
			mu.Unlock()

			return err
		})
	}

	err := g.Wait()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("\n❌ CI Pipeline Failed after %s: %v\n", duration.Round(time.Second), err)
		return err
	}

	fmt.Printf("\n✅ CI Pipeline Passed Successfully in %s!\n", duration.Round(time.Second))
	return nil
}

// QuietRunCmd runs a command and only returns error, suppressing stdout/stderr unless there's an error.
func QuietRunCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...) //nolint:gosec // G204: Maintenance utility runs trusted CI tools
	// We can buffer output and only show if error, but for CI we often want to see it.
	// However, to avoid interleaving, we'll let RunParallelCI handle the start/end logs.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunPlaywrightInDocker executes Playwright tests inside a Linux Docker container to ensure platform parity.
func RunPlaywrightInDocker(cwd string) error {
	fmt.Println("🐳 Running Playwright E2E Tests in Linux Container...")

	// 1. Cross-compile the server for Linux (matching host architecture for Docker parity)
	fmt.Println("🦀 Cross-compiling server for Linux...")
	buildCmd := exec.Command("go", "build", "-o", "server-linux", "main.go")
	buildCmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH="+runtime.GOARCH)
	buildCmd.Dir = cwd
	if out, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to build linux binary: %v\n%s", err, out)
	}
	defer func() { _ = os.Remove(filepath.Join(cwd, "server-linux")) }()

	// 2. Resolve Docker tag dynamically from package.json
	image := getPlaywrightDockerTag(cwd)

	// 3. Run Playwright in Docker using the pre-built binary
	// We use -e AGBALUMO_TEST_SERVER_COMMAND to point to the linux binary.
	// We run 'npm ci' first to ensure node_modules parity.
	return QuietRunCmd("docker", "run", "--rm",
		"-v", cwd+":/app",
		"-w", "/app",
		"-e", "AGBALUMO_TEST_SERVER_COMMAND=./server-linux serve",
		"-e", "AGBALUMO_ENV=test",
		image,
		"sh", "-c", "npm ci && npx playwright test",
	)
}

// getPlaywrightDockerTag resolves the appropriate Playwright Docker image tag from package.json.
func getPlaywrightDockerTag(cwd string) string {
	fallback := "mcr.microsoft.com/playwright:v1.59.1-noble"
	data, err := os.ReadFile(filepath.Join(cwd, "package.json")) //nolint:gosec // G304: Maintenance utility reads project config
	if err != nil {
		return fallback
	}

	var pkg struct {
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fallback
	}

	ver, ok := pkg.DevDependencies["@playwright/test"]
	if !ok {
		return fallback
	}

	// Strip caret/tilde
	ver = strings.TrimLeft(ver, "^~")
	if ver == "" {
		return fallback
	}

	return fmt.Sprintf("mcr.microsoft.com/playwright:v%s-noble", ver)
}

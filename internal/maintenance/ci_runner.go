package maintenance

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

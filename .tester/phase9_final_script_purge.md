# Phase 9: Final Script Purge & Gosec Rationale Migration

## Objective
Eradicate the final obsolete verification scripts from the `scripts/` directory by porting the last linting rule (`check-gosec-rationale`) into our native compiled architecture and deleting unused lint wrappers. 

## Context
You have a script `verify-golangci-config.sh` that checks if the linter config is valid. This is completely obsolete now because our native `ci` loop runs the pinned binary directly and will fail automatically if the config is bad. 
Additionally, `scripts/utils/check-gosec-rationale.sh` strictly enforces that `// #nosec` comments have a justification using grep. This can easily be mapped directly into Go.

## Steps for Execution
1. Open `cmd/verify/main.go`.
2. Create a `gosecCmd`:
```go
var gosecCmd = &cobra.Command{
	Use:   "gosec-rationale",
	Short: "Verify that all // #nosec directives include a rationale comment",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Checking for mandatory rationale in // #nosec directives...")
		
		var invalid []string
		
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}
			
			// Exclude common directories to simulate bash excluded paths
			if strings.Contains(path, "/vendor/") || strings.Contains(path, "/.tester/") || strings.Contains(path, "/tmp/") {
				return nil
			}

			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return nil 
			}

			lines := strings.Split(string(content), "\n")
			for i, line := range lines {
				if strings.Contains(line, "// #nosec") || strings.Contains(line, "//#nosec") {
					// Check if a rationale hyphen exists
					if !strings.Contains(line, " - ") && !strings.Contains(line, " -- ") {
						invalid = append(invalid, fmt.Sprintf("%s:%d -> %s", path, i+1, strings.TrimSpace(line)))
					}
				}
			}
			return nil
		})

		if err != nil {
			return err
		}

		if len(invalid) > 0 {
			fmt.Println("❌ Error: Found // #nosec directives without a mandatory rationale comment.")
			fmt.Println("Rationale must be preceded by a hyphen (-) or double-hyphen (--).")
			for _, issue := range invalid {
				fmt.Println("  ", issue)
			}
			return fmt.Errorf("mandatory rationale missing for %d occurrences", len(invalid))
		}

		fmt.Println("✅ All // #nosec directives have rationales.")
		return nil
	},
}
```
3. Add `gosecCmd` to `rootCmd.AddCommand()` in `main()`.
4. Delete `scripts/verify-golangci-config.sh`.
5. Delete `scripts/utils/check-gosec-rationale.sh`.
6. Delete the now-empty `scripts/utils/` directory: `rmdir scripts/utils/`.
7. Commit cleanly natively.

## Verification
- Both shell scripts are deleted.
- Running `go run cmd/verify/main.go gosec-rationale` executes efficiently in pure Go.

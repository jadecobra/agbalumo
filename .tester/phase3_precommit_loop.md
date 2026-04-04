# Phase 3: Fast Pre-Commit Git Hook Engine

## Objective
Extend `cmd/verify/main.go` with a `precommit` command that runs highly optimized, parallelized checks restricted only to staged files for maximum local feedback speed.

## Context
The legacy pre-commit hook relies on complex bash matching using `git diff`. By doing this in Go, we guarantee native cross-platform matching.

## Steps for Execution
1. Open `cmd/verify/main.go`.
2. Add a helper function to fetch staged files:
```go
func getStagedFiles(extension string) ([]string, error) {
	out, err := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACMR").Output()
	if err != nil {
		return nil, err
	}
	var files []string
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		if strings.HasSuffix(l, extension) {
			files = append(files, l)
		}
	}
	return files, nil
}
```
3. Create a `precommitCmd` cobra command.
4. Inside the `RunE` method, implement the fast loop steps:
   - **Mod Tidy**: Run `go mod tidy`, followed by `git diff --exit-code go.mod go.sum` to ensure no drift.
   - **Fmt**: Retrieve staged `.go` files. Run `gofmt -w` passing the fetched files.
   - **Build**: Retrieve staged `.go` files. Run `go build -o /dev/null ./...` (a fast syntax and dependency check).
   - **Lint Stage**: Run `go run github.com/golangci/golangci-lint/cmd/golangci-lint run -c scripts/.golangci.yml --new-from-rev=HEAD`.
5. Register the `precommitCmd` to the root `cobra.Command`.
6. Commit changes with message: `feat(ci): implement go native precommit optimizations`.

## Verification
- Stage some modified `.go` files.
- Run `go run cmd/verify/main.go precommit` locally. Confirm it formats the code and runs the linter exclusively on the diff.

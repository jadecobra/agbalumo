# Phase 2: CI Pipeline Orchestration in Go

## Objective
Extend `cmd/verify/main.go` to have a dedicated `ci` subcommand that natively orchestrates our CI pipeline checks in pure Go.

## Context
Our CI pipeline heavily relies on YAML orchestration which causes dev/prod drift. We are moving this orchestration into our native `cmd/verify` application.

## Steps for Execution
1. Open `cmd/verify/main.go`.
2. Import `"os/exec"` if not already imported. Add a helper function to run shell commands efficiently:
```go
func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
```
3. Add a new `cobra.Command` variable called `ciCmd`.
4. Inside the `RunE` function of `ciCmd`, sequentially execute the following commands using the helper. Print clear demarcations (e.g., `fmt.Println("=== Running Lint ===")`) between stages.
   - **Lint**: `runCmd("go", "run", "github.com/golangci/golangci-lint/cmd/golangci-lint", "run", "-c", "scripts/.golangci.yml")`
   - **Test**: `runCmd("go", "test", "-race", "-cover", "-count=1", "./...")`
   - **Vulncheck**: `runCmd("go", "run", "golang.org/x/vuln/cmd/govulncheck", "./...")`
   - **API/CLI Contract Drift**: Call the `apiSpecCmd.RunE(cmd, args)` manually or invoke `runCmd("go", "run", "cmd/verify/main.go", "api-spec")`.
   - **Template Drift**: Call `templateDriftCmd.RunE(cmd, args)`.
5. Add `ciCmd` to `rootCmd.AddCommand(...)` in the `main()` function.
6. Commit the changes natively with message: `feat(ci): implement native go ci pipeline orchestration`.

## Verification
- Run `go run cmd/verify/main.go ci` and confirm that all checks execute sequentially and correctly pass (or appropriately fail if a test breaks).

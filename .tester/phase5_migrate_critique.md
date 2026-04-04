# Phase 5: ChiefCritic Migration to Pure Go

## Objective
Migrate the functionality of `scripts/critique.sh` strictly into `cmd/verify/main.go` using the `tools.go` pinned binaries, and permanently delete the bash version.

## Context
The legacy bash script attempts to resolve archaic `GOBIN` paths. Since we implemented `tools.go` in Phase 1, we can invoke these tools accurately without any shell dependencies. 

## Steps for Execution
1. Open `cmd/verify/main.go`.
2. Add a new `critiqueCmd` variable:
```go
var critiqueCmd = &cobra.Command{
	Use:   "critique",
	Short: "Run ChiefCritic robustness audit natively",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Starting ChiefCritic Robustness Audit...")

		fmt.Println("\n[1/4] Checking Cognitive Complexity (gocognit)...")
		if err := runCmd("go", "run", "github.com/uudashr/gocognit/cmd/gocognit", "-over", "10", "./cmd", "./internal"); err != nil {
			fmt.Println("❌ Complexity threshold exceeded!")
			return err
		}

		fmt.Println("\n[2/4] Checking Repeated Strings (goconst)...")
		_ = runCmd("go", "run", "github.com/jgautheron/goconst/cmd/goconst", "./cmd/...", "./internal/...")

		fmt.Println("\n[3/4] Checking Struct Alignment (fieldalignment)...")
		_ = runCmd("go", "run", "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment", "./internal/...", "./cmd/...")

		fmt.Println("\n[4/4] Checking Code Duplication (dupl)...")
		_ = runCmd("go", "run", "github.com/mibk/dupl", "-threshold", "15", "-t", "./cmd", "./internal")

		fmt.Println("\n✅ ChiefCritic Audit Complete!")
		return nil
	},
}
```
3. Add `critiqueCmd` to `rootCmd.AddCommand(...)` in `main()`.
4. Run `go run cmd/verify/main.go critique` to confirm it executes correctly.
5. Create an artifact `critique_report.md` manually to log any initial findings.
6. Delete `scripts/critique.sh`.
7. Commit changes natively: `refactor(ci): migrate critique audit to pure go orchestration`.

## Verification
- `scripts/critique.sh` should no longer exist.
- `go run cmd/verify/main.go critique` should run all four tools seamlessly.

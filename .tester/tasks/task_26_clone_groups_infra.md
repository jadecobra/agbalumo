# Task 26: Clone Group Reduction — Batch A (verify tooling & infra)

## Context

`go run cmd/verify/main.go critique` [4/4] reports **358 clone groups**. This batch
targets the lowest-hanging, highest-impact groups in the **verify toolchain and
maintenance infrastructure**. These are non-test, production-adjacent files where
clone extraction has zero regression risk.

---

## Target Clone Groups

### 1. `cmd/verify/main.go` — 6-clone `cobra.Command` boilerplate (lines 157–353)

**Pattern**: Every single-delegate command (`verifyShasCmd`, `ciToolsCmd`, `gitleaksCmd`,
`ignoredFilesCmd`, `critiqueCmd`, `perfCmd`) is a 7-line `var xCmd = &cobra.Command{...}`
block that wraps exactly one `maintenance.Xxx(...)` call.

**Fix**: Define a helper `makeSimpleCmd(use, short string, fn func() error) *cobra.Command`
that returns the cobra.Command. Replace all 6+ identical command structs with one-liner
calls to this helper.

```go
// Before (repeated 6 times):
var verifyShasCmd = &cobra.Command{
    Use:   "verify-shas",
    Short: "...",
    RunE: func(cmd *cobra.Command, args []string) error {
        return maintenance.VerifyActionSHAs(".")
    },
}

// After:
var verifyShasCmd = makeSimpleCmd("verify-shas", "...", func() error {
    return maintenance.VerifyActionSHAs(".")
})
```

### 2. `cmd/verify/main.go:72–81` — 2-clone map-building loop

**Pattern**: Two identical `for _, c := range list { map[c] = true }` loops building
`codeMap` and `mdMap`.

**Fix**: Extract `toSet(items []string) map[string]bool` helper function.

### 3. `internal/maintenance/apispec.go:28–55` & `ast.go:115–117` — 3-clone error pattern

**Pattern**: `if err != nil { return fmt.Errorf("...") }` repeated for the same
shape of "could not read file X" error wrapping.

**Fix**: Extract `readFileOrErr(path, label string) ([]byte, error)` helper in `maintenance`.

### 4. `internal/maintenance/util.go:29–34` vs `internal/util/slices.go:12–17`

**Pattern**: Identical slice dedup logic exists in two packages.

**Fix**: Keep the implementation in `internal/util/slices.go`. Update
`internal/maintenance/util.go` to call `util.UniqueStrings(...)` instead of duplicating.

### 5. `cmd/admin.go:82–85`, `cmd/aglog/main.go:15–18`, `cmd/listing_update.go:90–93`

**Pattern**: Three `cmd/*.go` files repeat the same 4-line "load db + handle error" block.

**Fix**: Extract `mustOpenDB(cfg Config) Repository` or equivalent helper into
`cmd/shared.go` (new file). All three callers adopt it.

---

## Verification

- [ ] `go build ./...`
- [ ] `go run cmd/verify/main.go critique` — confirm clone count for these files is reduced
- [ ] `go test ./...`
- [ ] Commit: `refactor(verify): extract cobra command helpers and reduce clone groups`

# Phase 1: Tool Pinning (Pure Go Build Architecture)

## Objective
Convert all global `go install` CI tools into a unified `tools.go` module-pinned paradigm. This natively guarantees dev/prod version parity without relying on bash scripts or `Taskfile.yml`.

## Context
Currently, the pipeline uses archaic `go install` targets with custom `GOBIN` hacks. We want to lock tool versions in `go.mod` using the blank import pattern.

## Steps for Execution
1. Create a new file at `tools/tools.go`.
2. Add the following content:
```go
//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/jgautheron/goconst/cmd/goconst"
	_ "github.com/mibk/dupl"
	_ "github.com/uudashr/gocognit/cmd/gocognit"
	_ "github.com/zricethezav/gitleaks/v8"
	_ "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
```
3. Run `go mod tidy` in the terminal to resolve and lock the dependencies in `go.mod` and `go.sum`.
4. Verify that running `go run github.com/golangci/golangci-lint/cmd/golangci-lint --version` works successfully (it should automatically download and run the binary version pinned in `go.mod`).
5. Commit the changes to `tools/tools.go`, `go.mod`, and `go.sum` with a message: `build(tools): pin CI toolchain versions natively via tools.go`.

## Verification
- Running `go mod tidy` should complete without removing the dependencies.
- You should be able to run `go run golang.org/x/vuln/cmd/govulncheck ./...` natively.

---
name: "Verify Authoring"
description: "Boilerplate pattern for creating new verify CLI subcommands with TDD."
triggers:
  - "add verify subcommand"
  - "new verify tool"
  - "automate this check"
  - "create verify command"
mutating: true
---

# Verify Authoring Skill

## Trigger
When a manual check should be automated as a `verify` subcommand.

## File Structure (5 files, always the same)

### 1. `internal/maintenance/<name>.go`
```go
package maintenance

// Check<Name> performs deterministic verification of <what>.
func Check<Name>(rootDir string) ([]<Name>Violation, error) {
    // Implementation
}

type <Name>Violation struct {
    File    string
    Line    int
    Message string
}
```

### 2. `internal/maintenance/<name>_test.go`
```go
package maintenance

func TestCheck<Name>(t *testing.T) {
    tmpDir := t.TempDir()
    // Create fixture files
    // Call Check<Name>(tmpDir)
    // Assert violations
}
```
Use table-driven tests. Use `t.TempDir()` for fixtures. Minimum 3 test cases.

### 3. `cmd/verify/<name>.go`
```go
package main

import (
    "fmt"
    "github.com/jadecobra/agbalumo/internal/maintenance"
    "github.com/spf13/cobra"
)

var <name>Cmd = makeSimpleCmd("<name>", "<description>", func() error {
    // Call maintenance.Check<Name>(".")
    // Print violations or success
})
```

### 4. Register in `cmd/verify/main.go`
Add `rootCmd.AddCommand(<name>Cmd)` in the `init()` function.

### 5. Add to `.agents/verify-manifest.yaml`
```yaml
- name: <name>
  trigger: <when_to_run>
  description: "<what it checks>"
```

## Verification Checklist
```bash
go test ./internal/maintenance/ -run TestCheck<Name> -v
go run ./cmd/verify <name>
go run ./cmd/verify skill-conformance
go run ./cmd/verify check-resolvable
go build ./...
```

## Commit Convention
`feat(maintenance): add verify <name> for <purpose>`

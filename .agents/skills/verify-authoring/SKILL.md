---
name: "Verify Authoring"
description: "Boilerplate pattern for creating new verify CLI subcommands with TDD."
triggers:
  - "add verify subcommand"
  - "new verify tool"
  - "automate this check"
  - "create verify command"
  - "/skill-audit"
  - "skill audit"
mutating: true
---

# Verify Authoring Skill

## Trigger
When a manual check should be automated as a `verify` subcommand.

**Pre-check (cheapest route first)**: Before authoring, confirm no existing `verify` subcommand or `.golangci.yml` linter (e.g. `unused`, `gocognit`) already covers the class. Enabling a linter beats a new subcommand.

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
Registration surface (miss one = precommit/api-spec failure): `cmd/verify/<name>.go`, `cmd/verify/main.go`, `cmd/verify/main_test.go`, `.agents/verify-manifest.yaml`, `docs/cli/verify.md`, then `go run ./cmd/verify minify-context` to refresh `.agents/bundle.min.md`.

```bash
go test ./internal/maintenance/ -run TestCheck<Name> -v
go run ./cmd/verify <name>
go run ./cmd/verify skill-conformance
go run ./cmd/verify check-resolvable
go run ./cmd/verify api-spec
go build ./...
```

**Gate proof (post-flight)**: A check not invoked by `precommit` or `ci` enforces nothing. Wire it into the precommit or ci gate in `cmd/verify/ci.go` (or document why it is manual-only) and confirm it appears in `go run ./cmd/verify ci` output.

## Skill Completeness Audit (7-Item Checklist)
When authoring or auditing a skill (`/skill-audit <name>`), confirm all seven requirements.
1. `SKILL.md` exists with valid YAML frontmatter containing `name`, `description`, `triggers`, and `mutating`.
2. Procedural steps or executable instructions are clearly specified.
3. Any deterministic logic has corresponding unit tests in `internal/maintenance/`.
4. Resolver entry exists in `.agents/skills/RESOLVER.md` under `## Procedural Skills`.
5. `go run ./cmd/verify skill-conformance` passes cleanly.
6. `go run ./cmd/verify check-resolvable` passes cleanly.
7. Registration entry exists in `.agents/verify-manifest.yaml` under `skills:`.

## Commit Convention
`feat(maintenance): add verify <name> for <purpose>` or `chore(agents): skill-audit <name> to 7/7`

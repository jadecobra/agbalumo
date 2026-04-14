---
description: Coding Standards and Guidelines (Go, HTMX, Tailwind)
---

# Coding Standards & Guidelines

This workflow is referenced for specific edge-case rules regarding code style, imports, error handling, and testing. Use this as a reference when writing new features.

## Code Style Guidelines

### Imports (Group Order)
1. Standard library (blank line)
2. Third-party packages (blank line)
3. Local packages

### Naming Conventions
- Packages: lowercase single word (`domain`, `handler`, `service`)
- Types: PascalCase (`ListingHandler`, `UserStore`)
- Interfaces: `XxxStore`, `XxxService` (`ListingStore`)
- Errors: `ErrSomething` as package-level vars (`ErrInvalidDeadline`)

### Structs & Types
- Use struct tags: `json`, `form` for API binding
- Constructor pattern: `NewXxx()` functions

### Error Handling
- Use `RespondError(c, err)` for HTTP handlers - logs internally, renders friendly error page
- Wrap with `echo.NewHTTPError` for specific codes.
- Handlers should generally return `error`.

## Context Cost Awareness (Tokens)

To maintain agentic efficiency, we monitor **Token Density**.
- **Advisory TokenRMS**: Target **< 110**.
- **Context Window**: Monitor `ContextWindowPct` relative to Claude Sonnet (200k tokens).
- **Efficiency Pattern**: If a file exceeds **500 tokens**, consider if splitting into sub-packages or smaller files would improve logical cohesion and "Agentic Attention."
- **Janitor Run**: Use `/janitor` to clean up high-cost or high-entropy files when the TokenRMS exceeds thresholds significantly.
- **Verification**: Verify token density at any time by running: `go run cmd/verify/main.go context-cost`.

## Security & Linter Suppression

- **Suppression Justification:** Any time you use a `#nosec` directive, you must include a valid justification comment. Validate this project constraint by running `go run cmd/verify/main.go gosec-rationale`.

## Testing Conventions

### Test Structure
Use table-driven tests:
```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        expectErr bool
    }{
        {name: "valid case", input: "value", expectErr: false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) { ... })
    }
}
```

### Mocks
Use `github.com/stretchr/testify/mock`. Place mocks in `internal/mock/`.

### Test Helpers & Anti-Duplication
- **Use Standards First:** Always check `internal/testutil/` for centralized UI, authentication, or seeding helpers (like `setupAdminTestContext`, `testutil.NewMainTemplate`) before creating custom ones inline.
- **Prevent Copy-Paste:** To avoid failing the `cmd/verify/main.go critique` toolchain, any repetitive setup boilerplate or repeated logic between subtests *must* be extracted into explicit helper functions.
- **Cognitive Complexity limits:** High-complexity test functions (with deep nesting or serial sequential assertions) will fail the project's quality gates. Extract large blocks into private helpers.

### Coverage & Rules
- **TDD:** Write failing tests (RED) FIRST.
- **Coverage:** Thresholds are defined in `.agents/coverage-thresholds.json`. NEVER lower these value.
- **Persona Sync:** Changes to rules MUST be mirrored in `GEMINI.md` and `.agents/rules/`.
- **Functions:** Keep functions small and single-purpose (SRP).
- **Comments:** Code should be self-documenting; avoid unnecessary comments.

> **UI Component Constraint**: Do NOT write raw HTML form, button, or modal elements directly into major page templates. You MUST utilize or propose existing component definitions located inside `ui/templates/components/`. All colors, spacing, and visual effects must be derived exclusively from `tailwind.config.js` logic natively—no arbitrary hex codes or custom CSS outside of the Tailwind engine.

## Strict Lessons

This section contains corrections and constraints derived from the `[/learn]` workflow. These rules take precedence over existing style guidelines.

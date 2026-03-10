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

### Coverage & Rules
- **TDD:** Write tests first. A feature isn't done until tests pass.
- **Coverage:** Threshold is enforced from `.agent/coverage-threshold`. NEVER lower this value — write more tests instead.
- **Functions:** Keep functions small and single-purpose (SRP).
- **Comments:** Code should be self-documenting; avoid unnecessary comments.

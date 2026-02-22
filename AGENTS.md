# Agbalumo - Agent Guidelines

## Project Overview

Agbalumo is a Go web application for the West African diaspora community, featuring a business directory, job board, event listings, and community requests. Built with Echo framework, SQLite, HTMX, and Tailwind CSS.

---

## Build/Lint/Test Commands

### Run All Tests
```bash
go test ./...
```

### Run Single Test
```bash
go test -v -run TestFunctionName ./internal/package/
go test -v -run TestFunctionName/SubtestName ./internal/package/
```

### Run Tests with Race Detection
```bash
go test -race ./...
```

### Pre-Commit Quality Gate
```bash
./scripts/pre-commit.sh
```
Runs: `gofmt`, `go mod tidy`, `go vet`, race tests, coverage (>=89.2%), secret scanning.

### Build & Restart Server
```bash
./scripts/verify_restart.sh
```
Builds CSS, compiles binary, restarts server on :8443.

### Build CSS
```bash
npm run build:css
```

---

## Code Style Guidelines

### Imports (Group Order)
1. Standard library (blank line)
2. Third-party packages (blank line)
3. Local packages

```go
import (
    "context"
    "net/http"

    "github.com/labstack/echo/v4"

    "github.com/jadecobra/agbalumo/internal/domain"
)
```

### Naming Conventions
- **Packages**: lowercase single word (`domain`, `handler`, `service`)
- **Types**: PascalCase (`ListingHandler`, `UserStore`)
- **Interfaces**: `XxxStore`, `XxxService` (`ListingStore`)
- **Errors**: `ErrSomething` as package-level vars (`ErrInvalidDeadline`)

### Structs & Types
- Use struct tags: `json`, `form` for API binding
- Constructor pattern: `NewXxx()` functions

```go
type ListingHandler struct {
    Repo         domain.ListingStore
    ImageService service.ImageService
}

func NewListingHandler(repo domain.ListingStore, is service.ImageService) *ListingHandler {
    return &ListingHandler{Repo: repo, ImageService: is}
}
```

### Error Handling
- Use `RespondError(c, err)` for HTTP handlers - logs internally, renders friendly error page
- Wrap with `echo.NewHTTPError` for specific codes:

```go
return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Validation Error: "+err.Error()))
```

### Handlers
- Return `error` from all handlers
- Use `c.Render()` for page templates, `c.JSON()` for API responses
- Get user from context: `c.Get("User")`

---

## Testing Conventions

### Test Structure
```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        expectErr bool
    }{
        {name: "valid case", input: "value", expectErr: false},
        {name: "invalid case", input: "", expectErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### Mocks
Use `github.com/stretchr/testify/mock`. Mocks go in `internal/mock/`.

```go
mockRepo.On("FindByID", ctx, "123").Return(domain.Listing{}, errors.New("not found"))
```

---

## Architecture

```
cmd/           CLI commands (Cobra)
internal/
  config/      Configuration
  domain/      Core types, interfaces, business rules
  handler/     HTTP handlers (Echo)
  middleware/  Auth, sessions, rate limiting
  mock/        Test mocks
  repository/  Data access interfaces
  service/     Business logic layer
  ui/          Template renderer
ui/
  templates/   HTML templates (Go templates)
  static/      CSS, JS, images
```

---

## Feature Development Workflow

**Follow the 3-layer verification** (`.agent/workflows/feature-implementation.md`):

1. **Layer 1 - Unit Tests**: Write failing test (Red) → Implement (Green) → Refactor
2. **Layer 2 - CLI/Integration**: Run `./scripts/pre-commit.sh` + `./scripts/verify_restart.sh`
3. **Layer 3 - Browser**: Use browser subagent to verify UI works

---

## Key Rules

- **TDD**: Write tests first. A feature isn't done until tests pass.
- **Coverage**: NEVER lower the 89.2% threshold - write more tests instead
- **Commits**: Short, imperative mood ("add user auth" not "added user auth")
- **Functions**: Small, single-purpose (SRP)
- **NO comments**: Code should be self-documenting
- **NO committing**: `ARCHITECTURE_CRITIQUE.md` (it's in `.gitignore`)
- **Always restart**: Run server restart workflow after changes pass verification

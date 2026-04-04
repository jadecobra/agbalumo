# Phase 13: The Agent-Optimized Architecture (AppEnv)

## Objective
Eradicate Dependency Injection (DI) constructor bloat and the rigid "Lego-brick" pattern. Consolidate all database stores, configuration, and loggers into a single unified `AppEnv` (Application Environment) struct.

## Context
Currently, `server.go` wires together handlers by passing massive lists of individual interfaces (e.g., `UserStore`, `ListingStore`, `CategoryStore`, `ClaimStore`). This forces Agents to traverse and modify 5-6 files (Interfaces, Mocks, Structs, Routes, and Tests) just to add a single DB column or query. 
By wrapping our `*sqlite.SQLiteRepository` and `*config.Config` into a master `AppEnv` struct, handlers inherit the entire database universe instantly. This makes AI Agent iteration frictionless.

## Steps for Execution
1. **Create the `AppEnv` struct** (e.g., in `internal/infra/env/env.go` or `internal/domain/env.go`):
   ```go
   package env
   
   import (
       "log/slog"
       "github.com/jadecobra/agbalumo/internal/config"
       "github.com/jadecobra/agbalumo/internal/repository/sqlite"
   )
   
   type AppEnv struct {
       DB     *sqlite.SQLiteRepository // Concrete type. No mocks.
       Config *config.Config
       Logger *slog.Logger
   }
   ```
2. **Refactor `server.go`**:
   - Initialize the `AppEnv`:
     ```go
     app := &env.AppEnv{
         DB:     repo,
         Config: cfg,
         Logger: slog.Default(),
     }
     ```
   - Update your module initialization to replace `admin.AdminDependencies{...}` with simply `admin.NewAdminHandler(app)`
3. **Refactor the Handlers (e.g. `internal/module/admin/admin.go`)**:
   - Replace the `AdminDependencies` struct. Handlers now hold `app *env.AppEnv`.
   - Update HTTP routing logic to use `app.DB.FindUser(...)` directly.
4. **Purge the Mocks**:
   - If tests are currently using mock generated files for `UserStore` or `ListingStore`, delete them. 
   - Update the tests to use `:memory:` SQLite connections passed into a test `AppEnv`. 

## Verification
- Handlers require precisely 1 argument in their constructor (`app *AppEnv`).
- Running `go run cmd/verify/main.go ci` passes, proving that real SQLite tests work flawlessly without mock interfaces.
- Agent context overhead drops significantly as interface files are deleted.

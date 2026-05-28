# Project Atlas: Intent → Code Map

## Request Lifecycle
User Request → `cmd/serve.go` → `cmd/server.go::SetupServer()` → `internal/infra/server/server.go::Setup()` → `setupRoutes()` → Module `RegisterRoutes()`

## DI Container
`internal/infra/env/env.go::AppEnv` — centralizes DB, Config, Logger, and all service interfaces. All module handlers receive `*AppEnv`.

## Intent Map

| User Intent | Handler | Data Access | Template |
|---|---|---|---|
| Search for food | `internal/module/listing/listing.go` | `internal/repository/sqlite/sqlite_listing_read.go` | `ui/templates/index.html` → `ui/templates/partials/listing_list.html` |
| View listing detail | `internal/module/listing/listing.go` | `internal/repository/sqlite/sqlite_listing_read.go` | `ui/templates/partials/modal_detail.html` |
| Save a listing | `internal/module/listing/listing_save.go` | `internal/repository/sqlite/sqlite_saved_listing.go` | `ui/templates/partials/save_button.html` |
| Filter saved listings | `internal/module/listing/listing_save.go` | `internal/repository/sqlite/sqlite_saved_listing.go` | `ui/templates/partials/listing_list.html` |
| Create listing | `internal/module/listing/listing_mutations.go` | `internal/repository/sqlite/sqlite_listing_write.go` | `ui/templates/partials/modal_create_listing.html` |
| Edit listing | `internal/module/listing/listing_mutations.go` | `internal/repository/sqlite/sqlite_listing_write.go` | `ui/templates/partials/modal_edit_listing.html` |
| Admin dashboard | `internal/module/admin/dashboard.go` | `internal/repository/sqlite/sqlite_stats.go` | `ui/templates/admin_dashboard.html` |
| Auth (Google OAuth) | `internal/module/auth/` | `internal/repository/sqlite/sqlite_user.go` | redirect flow |
| Submit feedback | `internal/module/feedback/` | `internal/repository/sqlite/feedback.go` | `ui/templates/partials/modal_feedback.html` |
| About page | `internal/common/page_handler.go` | none | `ui/templates/about.html` |

## Module Structure (actual handler locations)
- `internal/module/listing/` — Public listing CRUD, search, save, featured
- `internal/module/admin/` — Admin dashboard, moderation, user management  
- `internal/module/auth/` — Google OAuth + dev login
- `internal/module/feedback/` — User feedback submission
- `internal/common/` — Static page handlers (about)

## Key Interfaces (domain layer)
- `domain.ListingStore` = `ListingReader` + `ListingWriter` (CQRS split)
- `domain.ListingRepository` = composed super-interface (all stores)
- Implementation: `internal/repository/sqlite/sqlite.go::SQLiteRepository`

## Verification Tools
- Run `go run ./cmd/verify map --depth 2` for directory tree.
- Run `go run ./cmd/verify map --symbols` for exported type/func index.
- Run `go run ./cmd/verify map --templates` for template dependency graph.
- Run `go run ./cmd/verify schema` for SQLite schema.
- Run `go run ./cmd/verify trace` for request lifecycle observability.

## Agent Cold-Start Path (for .agents/AGENT-BOOTSTRAP.md consumers)

1. Read `.agents/AGENT-BOOTSTRAP.md` (the canonical 2200-2600 token first prose).
2. Read this ATLAS + run (or request operator paste of) `go run ./cmd/verify map --depth 2 --symbols --templates`.
3. Run (or request paste of) one domain-specific verify (e.g. `context-cost`, `trace`, `template-contract`).
4. You now have enough to sketch a medium feature change (new filter + cached invalidation) with exact files + the 2 verifies that would catch violations — without opening source.

**Non-local / cloud agents** (gemini cli, antigravity, opencode, grok build, etc.): Never guess paths. Always request the operator to run the exact `verify` command(s) and paste the full output. Treat pasted output as the only source of truth.

This section exists so the bootstrap + one verify paste satisfies the high autonomy bar.

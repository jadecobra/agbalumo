# Project Atlas: Intent → Code Map

## Request Lifecycle
User Request → `cmd/serve.go` → `cmd/server.go::SetupServer()` → `internal/infra/server/server.go::Setup()` → `setupRoutes()` → Module `RegisterRoutes()`

## DI Container
`internal/infra/env/env.go::AppEnv` — centralizes DB, Config, Logger, and all service interfaces. All module handlers receive `*AppEnv`.

## Intent Map

| User Intent | Handler | Data Access | Template |
|---|---|---|---|
| Search for food | `internal/module/listing/listing.go::HandleHome` | `internal/repository/sqlite/sqlite_listing_read.go::FindAll` | `ui/templates/index.html` → `partials/listing_list.html` → `partials/listing_card.html` |
| View listing detail | `internal/module/listing/listing.go::HandleDetail` | `sqlite_listing_read.go::FindByID` | `partials/modal_detail.html` |
| Save a listing | `internal/module/listing/listing_save.go::HandleSaveToggle` | `sqlite_saved_listing.go` | `partials/save_button.html` |
| Filter saved listings | `internal/module/listing/listing_save.go::HandleSavedListings` | `sqlite_saved_listing.go` | `partials/listing_list.html` |
| Create listing | `internal/module/listing/listing_mutations.go::HandleCreate` | `sqlite_listing_write.go::Save` | `partials/modal_create_listing.html` |
| Edit listing | `internal/module/listing/listing_mutations.go::HandleUpdate` | `sqlite_listing_write.go::Save` | `partials/modal_edit_listing.html` |
| Admin dashboard | `internal/module/admin/dashboard.go::HandleDashboard` | `sqlite_stats.go`, `sqlite_listing_read.go` | `admin_dashboard.html` |
| Auth (Google OAuth) | `internal/module/auth/` | `sqlite_user.go` | redirect flow |
| Submit feedback | `internal/module/feedback/` | `sqlite/feedback.go` | `partials/modal_feedback.html` |
| About page | `internal/common/page_handler.go::HandleAbout` | none | `about.html` |

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

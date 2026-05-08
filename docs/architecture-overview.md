# Architecture Overview

This document summarizes the core architecture, data flow, and layers of the Agbalumo application.

## 1. Directory Structure

- **cmd/**: Contains the main application entry points. Used heavily by the Cobra CLI framework.
- **cmd/serve.go**: Cobra command that boots the server. Delegates to `SetupServer()` in `cmd/server.go`.
- **cmd/server.go**: Calls `internal/infra/server.Setup(cfg)` to initialize Echo, DB, and all module handlers.
- **internal/domain/**: Core types, structs, interfaces, validations, business concepts.
- **internal/module/**: HTTP handler modules organized by domain (listing/, admin/, auth/, feedback/). Each module implements `domain.Registrar` to register its own routes.
- **internal/handler/**: Contains search latency tests only. Handler logic lives in `internal/module/`.
- **internal/infra/env/**: Contains `AppEnv` — the central dependency injection container passed to all module handlers.
- **internal/infra/server/**: Server bootstrap, route wiring, middleware setup, and background service initialization.
- **internal/middleware/**: Custom Echo middleware functions (Auth parsing, rate limiting).
- **internal/repository/sqlite/**: Database implementations interacting with the SQL driver, split into smaller scoped files.
- **internal/service/**: Logic layer handling external business components that span across multiple repositories.
- **internal/ui/**: Contains Go `html/template` initialization, dynamic template parsing logic, and template path handling.
- **ui/templates/**: Raw HTMX templates.
- **ui/static/**: Pre-compiled CSS and raw JS assets.

## 2. Server Architecture (Data Retrieval & Render Flow)
1. User requests `GET /listings/123`.
2. Echo Router matches the route registered in the module's RegisterRoutes() method (wired in internal/infra/server/server.go::setupRoutes) and calls the handler.
3. Handler utilizes the `Repo` (dependency-injected SQLite store) to fetch the `domain.Listing` entity `repo.FindByID()`.
4. Logic/Formatting is applied as needed by the Handler.
5. The Handler creates a strictly typed ViewModel struct (per the Strict ViewModel Mandate) and passes it to c.Render().
6. The UI package parses `modal_detail.html` combined with `base.html` (if applicable) and writes output to `c.Response()`.
7. Results are pushed to browser. Note that partial HTMX loads return partial HTML fragments rather than complete HTML bodies.

## 3. Storage / Database Architecture
Agbalumo uses a single-file SQLite database with Write-Ahead Logging (WAL) enabled in production for optimal read concurrency.

- **Storage Location**: Root or Docker volume at `agbalumo.db`.
- **FTS5 Full Text Search**: `listings_fts` virtual table runs underneath the primary `listings` table to provide blazingly fast full-text searching functionality without external services like Elasticsearch. Triggers automatically update this virtual table.
- **Pooling**: Connections are governed tightly. Max open connections are typically set to `1` because SQLite writes must be serialized.

## 4. UI Architecture
- **Templating**: Standard Go `text/template` libraries with partial templating logic. A `base.html` defines the layout, and fragments are injected via blocks.
- **HTMX**: Used for dynamic updates (loading states, paging, modals) without writing complex Javascript.
- **Tailwind CSS**: Utility classes drive component aesthetics rather than separate CSS files. `npm run build:css` shrinks standard Tailwind into a minified payload.

## 5. Security Gates
- Rate Limiting implemented globally in Echo Config.
- Secret parsing checks done via shell validation in `go run ./cmd/verify precommit`.
- CSRF verification middleware runs on specific POST operations.

## 6. Dependency Injection
- All application dependencies are centralized in `internal/infra/env/env.go`.
- `AppEnv` holds: DB (ListingRepository), Config, Logger, and all service interfaces (CSV, Geocoding, Image, Listing, Categorization, Metrics).
- Module handlers receive `*AppEnv` via their constructor (e.g., `listing.NewListingHandler(app)`).
- See `docs/ATLAS.md` for the complete intent-to-file mapping.


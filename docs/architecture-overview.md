# Architecture Overview

This document summarizes the core architecture, data flow, and layers of the Agbalumo application.

## 1. Directory Structure

- **cmd/**: Contains the main application entry points. Used heavily by the Cobra CLI framework.
- **cmd/server.go**: Application startup sequence, Echo framework routing, Middleware attaching, Environment loading.
- **internal/domain/**: Core types, structs, interfaces, validations, business concepts.
- **internal/handler/**: HTTP logic layer mapping Requests -> DB/Service -> HTML/JSON response. All HTTP routing is connected here.
- **internal/middleware/**: Custom Echo middleware functions (Auth parsing, rate limiting).
- **internal/repository/sqlite/**: Database implementations interacting with the SQL driver, split into smaller scoped files.
- **internal/service/**: Logic layer handling external business components that span across multiple repositories.
- **internal/ui/**: Contains Go `html/template` initialization, dynamic template parsing logic, and template path handling.
- **ui/templates/**: Raw HTMX templates.
- **ui/static/**: Pre-compiled CSS and raw JS assets.

## 2. Server Architecture (Data Retrieval & Render Flow)
1. User requests `GET /listings/123`.
2. Echo Router looks for a match in `cmd/server.go` and calls the handler: `ListingHandler.HandleDetail(c)`.
3. Handler utilizes the `Repo` (dependency-injected SQLite store) to fetch the `domain.Listing` entity `repo.FindByID()`.
4. Logic/Formatting is applied as needed by the Handler.
5. The Handler creates a generic map context (e.g. `map[string]interface{}{"Listing": listing, "User": user}`).
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
- Secret parsing checks done via shell validation in `scripts/pre-commit.sh`.
- CSRF verification middleware runs on specific POST operations.

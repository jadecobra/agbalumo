# AGENT-BOOTSTRAP.md

**Canonical first prose read** (high success bar). Read *before* full AGENTS.md, source, or other docs.

**Mission (60s North Star)**: ruthlessly simple systems for finding African food in <60 seconds. Prioritize User Value + Min Latency. Complexity Kill-Switch: delete additions increasing UI steps/DB latency without 2x utility.

## Mandatory Session Start
1. `go run ./cmd/verify preflight` (active rules + invariants).
2. Read `.agents/invariants.json` (db_engine=sqlite, csp_policy="script-src 'self'", etc.). Construct all env references from these values only.
3. Read `.agents/skills/RESOLVER.md` — match intent to triggers.
4. Read `.agents/verify-manifest.yaml` — command inventory (never guess paths/contracts).
5. Read this file.

*Rule*: Read matched `.agents/skills/*.md` before edits. Run relevant verify *before* every mutation; cite output in reasoning.

## Hexagonal Architecture Boundaries
- `internal/domain/`: pure types, structs, interfaces. Zero external deps.
- `internal/handler/`: thin. Bind, call service, respond error, render ViewModel. No raw maps.
- `internal/service/`: business logic spanning repositories.
- `internal/repository/`: SQLite cached reads, WAL mode, FTS5 search, transactions.

Composition Root: `main.go` -> `cmd/` -> `internal/infra/server.Setup()` -> `internal/infra/env.AppEnv` injected into module handlers.

## Intent -> File Map (Live: `go run ./cmd/verify map --depth 2 --symbols --templates`)
| Intent | Module | Data Access (sqlite) | Template example |
|---|---|---|---|
| Search | listing/listing.go | sqlite_listing_read.go + FTS5 | index.html + list |
| Details | listing/listing.go | sqlite_listing_read.go | modal_detail.html |
| Save | listing/listing_save.go | sqlite_saved_listing.go | save_button.html |
| Create/Edit | listing/listing_mutations.go | sqlite_listing_write.go | modal_create/edit |
| Admin | admin/dashboard.go | sqlite_stats.go | admin_dashboard.html |

## Verify Quickref (Cite output in reasoning)
- `preflight` — active rules & constraints.
- `context-cost` — TokenRMS (target <110), high-cost files.
- `map --depth 2 --symbols --templates` — architecture truth.
- `trace` — request lifecycle (middleware -> DB -> UI).
- `schema` — SQLite schema.
- `template-contract` — ViewModel enforcement.
- `design` / `design-evidence` — UI Dialect compliance & raw violations.
- `browser` — Playwright matrix.
- `minify-context` — current bundle state.

## Rules & Constraints
- **Agent/Host Boundary**: Do not write Go application code to solve LLM token/context limits.
- **Timer Floors**: Background run_command = reactive (no timer); `gh run watch` or network = 300s floor. Never guess.
- **Fail-Fast Budget**: Max 3 compilation/test loops. 2 consecutive failures -> HALT + output traceback.
- **CSP**: `script-src 'self'` only. No inline scripts.
- **No Paperwork**: No progress files. Commits and green remote CI are proof.
- **CI Watch**: Run `./scripts/pushw.sh` or `gh run watch` after push (300s timer floor).
- **Edit Batching**: Plan all hunks before editing. Batch edits in the same file using one `multi_replace_file_content` call.
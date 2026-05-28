# AGENT-BOOTSTRAP.md

**Canonical first prose read** (target 2200-2600 tokens). Read this *before* full AGENTS.md, source, or other docs. Enables high success bar: after this + ATLAS + one `verify` output, sketch a rule-compliant medium feature change (e.g. new filter + cache invalidation) with exact file paths and verify steps — citing only these artifacts.

**Mission (60s North Star)**: Build ruthlessly simple, high-utility systems so a user finds African food in any city in <60 seconds. Prioritize User Value + Minimal Latency over architectural purity. Complexity Kill-Switch: any addition that increases UI steps or DB latency without 2x utility must be deleted.

## Mandatory Session Start (5 steps)

1. `go run ./cmd/verify preflight` (active rules + modified domains + invariants dump).
2. Read `.agents/invariants.json` (protocol=https, port=8443, db_engine=sqlite, csp_policy="script-src 'self'", default_coverage=90, max_featured_listings=3). Construct all env references from these values only.
3. Read `.agents/skills/RESOLVER.md` — match intent to triggers. No exact slash for "improve agent onboarding"; treat as repository rule change + doc improvement.
4. Read `.agents/verify-manifest.yaml` — 25+ deterministic tools are ground truth. Never guess paths or contracts.
5. Read this file (AGENT-BOOTSTRAP.md) as the single primary prose on-ramp.

**Rule**: Read any matched `.agents/skills/*.md` before edits. Run relevant verify *before* every mutation; cite output in reasoning.

**Diverse agents note** (gemini cli, antigravity, opencode, grok build, Cursor, etc.): If you cannot execute local `go run ./cmd/verify`, ask the operator to run it (with flags) and paste the full stdout/JSON. Treat the pasted output as source of truth. Do not hallucinate file names.

## Hexagonal Architecture + DI (condensed)

Boundaries (never violate):
- `internal/domain/`: pure types, structs, interfaces, errors, validations. Zero external deps.
- `internal/handler/`: thin. Only bind, call service, `RespondError(c, err)`, render typed ViewModel. No raw maps in c.Render, no business logic.
- `internal/service/`: the Product Engine. Business logic spanning repos.
- `internal/repository/`: data access only (sqlite + cached). Production: WAL mode, FTS5 for search, `last_verified_at` zombie filter on all reads, bounding-box + Haversine for spatial, no `fmt.Sprintf` SQL concat, transactions for multi-step writes, deep-copy on cache return.

Composition root: `main.go` → `cmd.Execute()` (Cobra) → `cmd/{root,serve,server,shared}.go` → `internal/infra/server.Setup()` → `internal/infra/env.AppEnv` (holds all repos + services) injected into module handlers.

Request lifecycle: Echo route (registered via module `RegisterRoutes()`) → handler (HX-Request aware for fragments) → service/repo → strictly typed ViewModel → `c.Render()` (Go templates + HTMX OOB swaps).

Key modules (feature code lives here, not handler/):
- `internal/module/listing/` — search, save, create, edit, featured, geocode, UI cards/modals, cache busting.
- `internal/module/admin/`, `auth/` (Google OAuth + dev), `feedback/`, `user/`, `common/`.

Key interfaces (domain):
- `domain.ListingStore` = Reader + Writer (CQRS).
- `domain.ListingRepository` = super-interface.
- Impl: `internal/repository/sqlite/sqlite.go` + cached wrapper.

## Intent → File Map (live truth via verify)

| Intent                  | Primary Module                  | Data Access (sqlite)                  | Template example                  |
|-------------------------|---------------------------------|---------------------------------------|-----------------------------------|
| Search listings        | listing/listing.go             | sqlite_listing_read.go + FTS5        | index.html + listing_list.html   |
| View detail / modal    | listing/listing.go             | sqlite_listing_read.go               | modal_detail.html                |
| Save / filter saved    | listing/listing_save.go        | sqlite_saved_listing.go              | save_button.html + list          |
| Create / Edit listing  | listing/listing_mutations.go   | sqlite_listing_write.go              | modal_create/edit_*.html         |
| Admin dashboard        | admin/dashboard.go             | sqlite_stats.go                      | admin_dashboard.html             |
| Auth (OAuth)           | auth/                          | sqlite_user.go                       | redirect                         |

**Live map**: `go run ./cmd/verify map --depth 2 --symbols --templates`. Use output, never hardcode paths from memory.

## Verify Quickref (use these; cite output)

Core (run before edits):
- `go run ./cmd/verify preflight` — session rules + constraints.
- `go run ./cmd/verify context-cost` — TokenRMS (target <110), high-cost files.
- `go run ./cmd/verify map --depth 2 --symbols --templates` — architecture truth.
- `go run ./cmd/verify trace` — full request lifecycle (middleware → DB → UI).
- `go run ./cmd/verify schema` — DB contract (may require operator to supply fixture output if no local DB).
- `go run ./cmd/verify template-contract` — ViewModel enforcement.
- `go run ./cmd/verify deprecated` — scan map[string]interface{} etc.
- `go run ./cmd/verify agents-coverage` — 100% package docs gate.
- `go run ./cmd/verify minify-context` — current agent bundle state.

Others: `design`, `browser` (UI only), `cache-buster`, `doc-drift`, `sweep`.

**For non-local agents**: paste operator output. Use it to answer "which file?" questions.

## Token / Context Rules (from coding-standards)

- Advisory: TokenRMS < 110. Current baseline measured via `verify context-cost`.
- Rule: File >500 tokens → consider split for "Agentic Attention".
- This bootstrap exists to minimize first-read cost while hitting the autonomy bar.
- Run `context-cost` + `minify-context` after any doc change. Janitor high-entropy files when needed.
- Goal: one 2200-2600 token read + 1-2 verify outputs → medium feature sketch without source dive.

## Pre-Edit Ritual (required for compliance + success bar)

1. `preflight` + `context-cost` + `map --symbols` (or area-specific) + any domain verify.
2. Cite exact output lines in reasoning.
3. Sketch using *only* bootstrap + ATLAS + the cited verify output.
4. High bar achieved when you can name: exact module file(s), repo file(s), template(s), cache invalidation points, the 2 verifies that would catch violations, without opening source.

## Condensed Constraints (highest impact)

From AGENTS.md + coding-standards (full versions only when editing rules):
- No writing `cmd/` or `internal/` Go to solve agent token/context problems (AGENT/HOST BOUNDARY). Use docs + rule changes only.
- Timer floors: background run_command = reactive (no timer); `gh run watch` or network = 300s floor only. Never guess.
- High-tier (this model): max 3 serial compile/test loops per task. 2 consecutive failures → HALT + traceback.
- Reactive yield: trust harness wakeup on run_command; do not poll.
- UI changes (not this task): full `browser` locally before push; use semantic partials, no raw utilities for atoms, min-h-0 + overflow-hidden for modals, etc.
- All dynamic template fields require {{ else }} fallback.
- CSP: script-src 'self' only. No inline scripts.
- Cache bust via `verify cache-buster` only.
- "No Paperwork": Git commits (and green remote CI) are proof. No progress files.
- Remote CI Guard: push → immediately `./scripts/pushw.sh` or `gh run watch` (300s floor). Work incomplete until green.

## When to Read More

- Full AGENTS.md: tone, full SESSION START, git rules, anti-sycophancy.
- `.agents/workflows/coding-standards.md`: edge-case lessons + full strict UI/quota/testing rules (only on actual code change).
- `.agents/verify-manifest.yaml`: complete command inventory.
- Run `go run ./cmd/verify minify-context` for the current concatenated agent bundle.

**Self-measure**: After ingesting this + ATLAS + one verify paste, you should produce a correct, boundary-respecting sketch for a medium change (new filter affecting listing query + cached repo invalidation path) naming files + required verifies.

This file is the optimized on-ramp. Edit it only when it measurably reduces first-read cost or raises the autonomy bar. Re-verify context-cost after edits.

---
Maintained as repository rule change + doc. See AGENTS.md for full protocol. Last updated: implementation of 2025 agent onboarding review.
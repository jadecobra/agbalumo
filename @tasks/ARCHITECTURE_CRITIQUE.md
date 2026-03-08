# Architecture Critique — Agbalumo

## Overall Grade: **A+**
The application has matured significantly. The foundational architecture is clean, highly testable, and strictly typed. P0–P3 performance work is complete and a formal automated performance audit (`scripts/performance-audit.sh`) now runs on every commit. Two new fixes landed on 2026-03-04: SQLite connection pool capped at `MaxOpenConns(1)` and logo served as WebP (296KB → 64KB, 78% reduction). Remaining open items are low-priority warnings only.

---

## What's Working Well ✅

| Area | Notes |
| :--- | :--- |
| **Domain isolation** | `internal/domain` has zero external dependencies — pure Go structs, validation, and errors. |
| **TDD culture** | 90.3% coverage, pre-commit enforcement, race detection — exceptional for an MVP. |
| **Centralized Config & Logs** | Config is strictly typed (`internal/config`); `slog` issues structural JSON vs Text based on environment. |
| **Standardized Responses** | `RespondError` ensures consistent error logging and rendering, preventing data leaks. |
| **Interface Segregation** | Focused interfaces (`ListingStore`, `UserStore`, `FeedbackStore`, `AdminStore`) via composition. |
| **Middleware design** | Clean separation: rate limiting, sessions, security headers, CSRF — all composable. |
| **DB Tuning** | WAL mode, `busy_timeout=5000`, `synchronous=NORMAL`, `MaxOpenConns(1)`, compound index on `(is_active, status, type)`. |
| **FTS5 Search** | Trigram tokenizer enables indexed substring search across title/description/city. |
| **Caching** | In-memory TTL cache for `GetCounts`/`GetLocations` with `RWMutex` and value-copy returns. |
| **WebP Logo** | `logo.webp` (64KB) served via `<picture>` element with PNG fallback; saves ~232KB per page load. |
| **Automated Perf Audit** | `scripts/performance-audit.sh` enforces asset sizes, DB config, caching, a11y, and N+1 patterns on every commit. |

---

## Scorecard

| Dimension | Score | Notes |
| :--- | :--- | :--- |
| **Layer separation** | 9/10 | Auth middleware extracted, domain stays pure. |
| **Interface design** | 9/10 | ISP-compliant focused interfaces. |
| **Testability** | 9/10 | Interface mocks, httptest, >90% coverage. |
| **Security** | 8/10 | CSP, CSRF, HSTS, rate limiting, session hardening. Minor `unsafe-inline` in CSP. |
| **Scalability** | 9/10 | Paginated endpoints, parallel queries, FTS5 trigram search, cached counts, `TitleExists` EXISTS query. |
| **Frontend Perf** | 9/10 | Gzip, immutable cache headers, Chart.js admin-only, deferred scripts, font/CSS preloads, WebP logo (Lighthouse Perf: 98). |
| **Maintainability** | 9/10 | Clean decomposition, consistent error handling. |
| **Error handling** | 9/10 | Centralized via `RespondError` and safe `error.html` templates. |

---

## Performance Audit — Task Tracker

### 🔥 P0 — Critical (Immediate Impact)

- [x] **P0.1** Enable gzip compression — `e.Use(middleware.Gzip())` in `cmd/server.go`
- [x] **P0.2** Add `Cache-Control` headers to static files in `cmd/server.go`
- [x] **P0.3** Move `chart.umd.min.js` to admin-only pages — remove from `base.html`
- [x] **P0.4** Add `defer` to all `<script>` tags in `base.html`

### 🟠 P1 — High Priority

- [x] **P1.1** Add `LIMIT` to `FindAllByOwner` and add pagination
- [x] **P1.2** Add `rows.Err()` checks after all `for rows.Next()` loops in `sqlite.go`
- [x] **P1.3** Parallel home page queries — goroutines for `FindAll`, `GetCounts`, `GetFeaturedListings`
- [x] **P1.4** Add `TitleExists` with `EXISTS` query for duplicate title checks
- [x] **P1.5** Move `GOOGLE_MAPS_API_KEY` to `Config` struct (avoid per-request `os.Getenv`)

### 🟡 P2 — Medium Priority

- [x] **P2.1** Add SQLite FTS5 with trigram tokenizer for full-text search (replace LIKE queries)
- [x] **P2.2** Add `preload` hints for critical fonts/CSS in `base.html`
- [x] **P2.3** Add in-memory caching for `GetCounts` (refresh every 60s)
- [x] **P2.4** Paginate `GetAllUsers` admin endpoint

### ⚪ P3 — Nice to Have

- [x] **P3.1** Add health check endpoint (`/healthz`)
- [x] **P3.2** Remove `'unsafe-inline'` from CSP `script-src`
- [x] **P3.3** Batch `ExpireListings` with chunk size for WAL safety

---

## Performance Audit — 2026-03-04

> Audit run via `scripts/performance-audit.sh`. Lighthouse (mobile): **Perf 98 / A11y 89 / Best Practices 100 / SEO 92**.

### ✅ Fixed

- [x] **Perf.1** `db.SetMaxOpenConns(1)` — caps SQLite pool to match write concurrency; prevents goroutine contention at `database/sql` layer
- [x] **Perf.2** Convert `logo.png` (296KB) → `logo.webp` (64KB, 78% smaller); `<picture>` with PNG fallback in `base.html`
- [x] **Perf.3** Image Uploads: Pipeline now re-encodes all user uploads (JPEG/PNG/GIF) to lossy WebP format (CGo-free via `gen2brain/webp`), resulting in ~30% smaller storage and faster page loads.

### ✅ Resolved Warnings

- [x] **Warn.1** `output.css` is 100KB — Verified minified at 98KB; bumped threshold for Tailwind v4.
- [x] **Warn.2** Five icon-only `<button>` elements in mobile bottom nav lack `aria-label` — added accessible labels.
- [x] **Warn.3** Heuristic N+1 flags in `admin.go` lines 317 and 358 — documented safe bounded admin action.

### ⚙️ Tooling Added

- [x] `scripts/performance-audit.sh` — 7-check automated audit (assets, caching, DB, cache layer, a11y, N+1, live response times)
- [x] Integrated into `scripts/pre-commit.sh` as step 5 — exit 2 blocks commit, exit 1 (warnings) allows through

---

## UI Component & Aesthetic Critique — 2026-03-04

> UI components across all pages and modals evaluate strongly against the baseline "Sharp Editorial Earth" aesthetic established by the homepage.

### ✅ What's Working Well

* **Strict Sharpness:** Zero rounded corners (`rounded-none` implicitly) on modals, cards, buttons, and inputs.
* **Earth Palette:** Consistent use of custom dark earth tones (`earth-dark`, `earth-sand`, `earth-ochre`, `earth-clay`, `earth-cream`) with glassmorphism/blur effects.
* **Typography:** `font-serif` (Playfair Display) for dominant headers and `text-[10px] uppercase tracking-widest` for micro-copy and controls, maintaining a premium editorial feel.
* **Layout Patterns:** Shared modal container styles (`bg-earth-dark/95 backdrop-blur-xl border border-white/20 p-0 m-auto`) are used consistently across Create, Edit, Detail, and Profile modals.
* **Admin Dashboard:** Banners with `border-l-[6px] border-earth-ochre` and sharp components translate the brand perfectly into admin tooling.

### ✅ Resolved Deviations

- [x] **Inconsistent Modal Inputs:** `modal_edit_listing.html` and other partials have been updated to match the rich container inputs (`bg-earth-sand/10 border border-white/20 p-1`) established in `modal_create_listing.html`.
- [x] **Semantic Admin Badges:** Admin status badges now use branded dark tints (e.g., `bg-green-500/20 text-green-400` and `bg-earth-ochre/20`) rather than standard bright backgrounds, maintaining visual cohesion.

---

## Agentic Harness (10x Loop) — 2026-03-08

> Proposed enhancements to evolve the passive documentation into an active operational framework.

### 🚀 Planned Improvements

- [x] **Agent.1** Persona Shards — Implement `.agent/personas/` directory and `scripts/agent-exec.sh role <name>` to programmatically activate and enforce persona-specific constraints.
- [x] **Agent.2** Workflow State Machine — Implement `.agent/state.json` and `scripts/agent-exec.sh workflow` to track feature progress, active phases (Red/Green/Refactor), and "Gate" verification status.
- [x] **Agent.3** Automated Gate Verification — Implement `scripts/agent-gate.sh <gate_id>` to programmatically verify criteria (e.g., `gate verify red-test` ensures the test fails as expected).
- [ ] **Agent.4** Brand "Juice" Generation — Implement `scripts/generate-juice.sh` to parse `.agent/rules/brand.toon` and generate CSS tokens (`brand-tokens.css`) and Go constants for technical brand enforcement.
- [x] **Agent.5** Automated Drift Check — Integrate persona/coding-standard sync check into `scripts/pre-commit.sh` to enforce the "Double-Commit" rule for `agent.yaml`.
- [ ] **Agent.6** Loop-Optimized Task Template — Create `.agent/task-template.md` optimized for long-running agentic sessions with explicit browser recording slots.


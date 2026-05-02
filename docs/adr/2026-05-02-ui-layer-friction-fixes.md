# ADR: UI Layer Friction Fixes

## Status
Accepted

## Context
AI agent iteration on Go templates was 5-10x slower than necessary due to:
- Server restart required for every template change (no hot reload)
- CSS build step was manual and forgettable
- 3 handlers assembled listing data independently (divergence bugs)
- Templates were only testable via expensive browser verification
- Template dependency graph was opaque (85 cross-references, no tooling)

## Decisions
1. **Dev-mode hot reload**: `TemplateRenderer.Render()` recompiles templates
   on every request when `AGBALUMO_ENV != production`. ~5ms cost per request.
2. **CSS build in precommit**: `verify precommit` auto-runs `npm run build:css`
   when `.html` or `.css` files are staged. Output is re-staged automatically.
3. **Shared data builder**: `buildListingViewData()` ensures all listing render
   paths (HandleHome excepted for goroutine perf) pass identical base keys.
4. **Template render tests**: `renderPartial()` helper enables Go-driven testing
   of partials with real data. Replaces browser subagent for component verification.
5. **Template graph CLI**: `verify template-graph` shows callers, callees, and
   key gaps for any named template.

## Consequences
- (+) Agent iteration cycle drops from ~2 min to ~5 sec for template changes
- (+) SavedIDs-class bugs caught by Go tests, not browser
- (+) CSS drift eliminated from precommit gate
- (-) Dev-mode Render is ~5ms slower per request (acceptable)
- (-) Template render tests depend on real template files (not isolated)
- (-) HandleHome excluded from shared builder due to goroutine pattern

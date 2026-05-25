---
name: Browser Verification
description: Verify UI changes using browser subagent with proper environment detection
triggers:
  - "UI change"
  - "browser verification"
  - "layout check"
  - "viewport audit"
  - "visual regression"
mutating: true
---
# Browser Verification Skill

## Execution Strategy: Scoped vs. Full
- **Local Scoped Execution**: You are encouraged to run targeted, scoped Playwright commands (e.g., `go run ./cmd/verify browser -- --grep "YourTestName"` or `npx playwright test visual.spec.ts`) during local UI iteration to maintain near-instant feedback cycles.
- **Automated Verification**: Use `browser_subagent` for targeted manual review or aesthetic checks, and **always save a screenshot** of any mutated UI component to prove compliance.
- **Full-Suite Gate**: The complete Omni-Surface Matrix (`go run ./cmd/verify browser`) remains mandatory at the git pre-commit/CI boundary.

## Step 1.0: Static A11y Gate (NEW — run FIRST)
Run `go run ./cmd/verify a11y-check` — catches:
- Missing aria-label on icon-only buttons
- Missing for/id associations on form controls
- Contrast-suspect token combinations (token on surface pairs below 4.5:1)
  If violations exist, fix before proceeding to Step 1.5.

## Step 1.5: Server Uptime Check (MANDATORY)
Before running any browser-based tests, you MUST ensure the local server is running and reachable:
- Run `go run ./cmd/verify uptime`
- If it fails, start the server in the background using `go run ./cmd/verify watch &` (or equivalent) and wait for readiness.

## UI TDD Workflow (Aesthetic/Layout Tweaks)
When modifying templates, CSS, or client-side assets where standard Go unit tests do not apply:
1. **Identify Visual Issue**: Use `browser_subagent` to capture a baseline "Before" screenshot and isolate the layout flaw.
2. **Apply Changes**: Modify HTML templates or `input.css`.
3. **Compile & Reload**: Run `npm run build:css` (if applicable) and **RESTART** the server to clear the Go template cache.
4. **Visual Verification**: Check responsiveness, margins, and aesthetic quality across target viewports via `browser_subagent`.
5. **Regression Test**: Run `go run ./cmd/verify browser` before committing.

## Pre-flight (MANDATORY — always run before browser tasks)

1. **Read `.agents/invariants.json` — FIRST ACTION, NO EXCEPTIONS.** Extract `protocol` and `port`. Construct base URL as `{protocol}://localhost:{port}`. This file is the committed, canonical source derived from BASE_URL. Do NOT guess port numbers. Do NOT read `.env` (security risk).
2. Verify server is running: `go run ./cmd/verify uptime` — if it fails, start the server first using `go run ./cmd/verify watch`
## Verification Checklist
For EVERY UI element verified, you MUST check ALL of:
- [ ] **Exists**: `document.querySelector(selector)` returns non-null
- [ ] **Visible**: `element.offsetHeight > 0 && element.offsetWidth > 0`
- [ ] **Has Content**: `element.innerText.trim().length > 0`
- [ ] **Interactive**: Click/hover produces expected state change
- [ ] **Responsive**: Element is fully visible and usable across the Mandatory Viewports below.
- [ ] **Layout Integrity**: Sticky elements do not overlap content (verify `rect.top >= header.height`).
- [ ] **Fragment Sync**: Confirm OOB swap targets (e.g., `#featured-section`) updated independently of the main fragment.
- [ ] **State Sync**: Verify `window.filterState` or equivalent matches UI selection in the JS console.
- [ ] **Chaos Data Resilience**: The agent must manually inject 100-character strings into titles and delete image `src` attributes in the DOM before screenshotting, verifying that `line-clamp` and fallback backgrounds preserve the grid.
- [ ] **Touch Target Ergonomics**: Verify all clickable elements (filters, links, buttons) have a minimum physical interaction area of `44x44px` on mobile viewports.
## The Omni-Surface Verification Matrix (MANDATORY)
A holistic audit MUST cover these 4 domains across Mobile AND Desktop:
- 1. **Public Discovery**: `/` (Main Feed), Search Results.
- 2. **Detail Modals**: Modal Detail View (with long text).
- 3. **Mutation Modals**: Create Listing (Post), Feedback Modal.
- 4. **Admin Surfaces**: `/admin/login`, `/admin/dashboard`.


## Mandatory Viewports
For ANY layout change, you MUST verify at:
| Device | Resolution | Goal |
|--------|------------|------|
| Mobile | 375 x 812 | Check for overflow-x and menu accessibility |
| Tablet | 768 x 1024 | Check for column wrapping |
| Desktop| 1440 x 900 | Standard editorial layout check |
| Wide   | 1920 x 1080| Check for max-width constraints |
## Common Failure Patterns
| Symptom | Root Cause | Fix |
|---------|-----------|-----|
| Connection refused | Wrong port/protocol | Check `.agents/invariants.json` |
| Element exists but invisible | CSS `display: none` or `opacity: 0` | Check computed styles |
| Click has no effect | Duplicate event listeners | Check `initApp()` in `app.js` |
| Dropdown clipped | Viewport clearance | Add `open-upwards` logic |
| Stale content after fix | Browser cache | Increment `?v=N` in `head_meta.html` |
| Overlap on Mobile | Lack of dynamic padding | Use `calc(var(--nav-height) + padding)` |
| Layout breaks on real-world long text | No clamping | Apply `line-clamp-2` or `truncate` |
| Misclicks on mobile UI | Touch target < 44px | Add padding (e.g., `p-3`) to hit 44x44px minimum |


## Agent Targeting Rules
1. **Prefer data-testid**: Always use `[data-testid="..."]` for deterministic targeting in `browser_subagent`.
2. **Event Verification**: If an interaction fails, check `htmx.logger` or `htmx:configRequest` in the console to verify payload integrity.
3. **Internal State Audit**: If the DOM doesn't reflect a change, query `window.filterState` to determine if the logic layer is the bottleneck.
## Post-flight
1. Document each check result in `task.md` with pass/fail
2. For layout changes: capture before/after screenshots, save them as artifacts, and embed them in walkthrough
3. For interactive changes: describe the state transition verified
4. **Automated Fix Loop**: If a visual regression or design violation is found during the audit, you MUST NOT stop at reporting the error. Instead, you must automatically apply the minimal CSS/Tailwind fix, commit it atomically using the `style(design): <fix description>` conventional format, and capture an "After" screenshot to prove the resolution.


## Critical Failure Protocol (MANDATORY)
If the E2E verification fails during local development:
1. **Isolate and Debug Locally**: You are encouraged to run targeted, scoped Playwright assertions (`--grep`) to iteratively implement and verify the fix.
2. **Mandatory Post-Fix Sweep**: Before pushing or committing, you must execute the full, unfiltered pre-commit and local browser validation to ensure zero integration regressions.
3. **Report Failures**: If the full E2E suite fails and cannot be resolved locally, output the clinical failure log to the user and await direction.

## AXE Accessibility Parity Warning
Local Darwin rendering + warm browser cache can produce false-green AXE results that fail in the clean-room Linux CI container.
- For any change touching **color tokens, navigation, listing cards, or modal contrast**: run the full suite, then explicitly check `test-results/` for any `error-context.md` artifacts generated even on passing runs.
- If the previous push to CI failed on an AXE check that passed locally, the fix is NOT verified until CI is green. Do not declare victory.

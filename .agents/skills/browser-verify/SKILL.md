---
name: Browser Verification
description: Verify UI changes using browser subagent with proper environment detection
triggers:
  - "UI change"
  - "browser verification"
  - "layout check"
  - "viewport audit"
  - "visual regression"
mutating: false
---
# Browser Verification Skill

## Execution Strategy: Automated vs. Manual
- **Deterministic Checks**: You MUST ALWAYS run `go run ./cmd/verify browser` for deterministic layout and regression checks.
- **Manual/Exploratory Checks**: The `browser_subagent` tool should NOW ONLY be used for exploratory testing of new features, assessing visual/aesthetic quality, and dealing with external third-party integrations.
- **Continuous Automation**: If you find yourself repeating a manual subagent check, you MUST trigger the `/learn` workflow to extract that check into a Playwright test.

## UI TDD Workflow (Aesthetic/Layout Tweaks)
When modifying templates, CSS, or client-side assets where standard Go unit tests do not apply:
1. **Identify Visual Issue**: Use `browser_subagent` to capture a baseline "Before" screenshot and isolate the layout flaw.
2. **Apply Changes**: Modify HTML templates or `input.css`.
3. **Compile & Reload**: Run `npm run build:css` (if applicable) and **RESTART** the server to clear the Go template cache.
4. **Visual Verification**: Check responsiveness, margins, and aesthetic quality across target viewports via `browser_subagent`.
5. **Regression Test**: Run `go run ./cmd/verify browser` before committing.

## Pre-flight (MANDATORY — always run before browser tasks)

1. Read `.agents/invariants.json` — get `protocol` and `port` fields
2. Construct the base URL: `{protocol}://localhost:{port}`
3. Verify server is running: `lsof -i :{port}` — if no output, start server first
4. Read `.env` for `BASE_URL` as override — if set, use that instead
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

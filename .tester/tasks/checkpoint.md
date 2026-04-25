# Checkpoint: UI Filtering Logic & Layout Hardening

## Current State: STABLE
The filtering logic and "Editorial Brutalist" layout have been hardened and verified through both automated tests and manual browser audits.

### Completed Work
- [x] **Header Overlap Fix**: Resolved z-index and padding conflicts between the sticky header and hero section.
- [x] **HTMX Migration**: Replaced fragile JS event triggers with direct `hx-get`, `hx-include`, and `hx-vals` on search buttons.
- [x] **Component Hardening**: Fixed `button_sharp` to correctly render custom attributes and `data-testid`.
- [x] **OOB Swap Stabilization**: Ensured the `featured-section` is correctly cleared or updated during all discovery fragment requests.
- [x] **Automated Testing**: Added OOB swap assertions to Go integration tests and attribute rendering tests to the UI suite.
- [x] **Skill Codification**: Updated `.agents/skills/browser-verify/SKILL.md` with protocols for overlap audits and state sync.

---

## Errors & Regressions Encountered
| Issue | Root Cause | Resolution |
|-------|------------|------------|
| **Hero Obscured** | Insufficient padding and `header_content` overlap in `index.html`. | Implemented explicit header overrides and consistent padding. |
| **Filter Breakdown** | Native `dispatchEvent` not caught by HTMX; `window.filterState` sync failure. | Migrated to direct HTMX attributes on interactive elements. |
| **"Blind" Subagent** | Missing `data-testid` on templated components. | Updated `button_sharp` and registered `safeHTMLAttr` function. |
| **Stale Featured Data** | Fragment response lacked `#featured-section` OOB wrapper. | Added OOB swap markers to `listing_list.html` and mock templates. |
| **Mock Divergence** | `testutil.NewMainTemplate` didn't include HTMX markers. | Synchronized mock templates with production HTML structure. |

---

## Planned Next Steps

### 1. Client-Side Testing (High Priority)
- Introduce **Vitest** to unit test `filters.js` logic and `window.filterState` transitions.
- Eliminate "dark matter" in the client-side JavaScript layer.

### 2. Playwright Viewport Matrix
- Codify the responsive viewport audit (375px, 768px, 1440px, 1920px) into a Playwright suite.
- Automate the `rect.top >= header.height` overlap check.

### 3. Visual Regression
- Implement pixel-diffing for the Hero and Filters modal to prevent aesthetic drift from the "Editorial Brutalist" dialect.

### 4. Search Sensitivity Audit
- Review FTS5 SQLite logic to ensure partial matches in categories (e.g., "church") do not return broader unrelated results.

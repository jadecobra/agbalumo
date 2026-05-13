# ADR [2026-05-13-semantic-ui-atoms]: Abstracting Fragmented UI Atoms into Semantic Components

**Date**: 2026-05-13 **Status**: Accepted

## 1. Context & User Problem
Agents and engineers struggle with hardcoded Tailwind utilities across fragmented Go templates (e.g., `modal_profile.html`, `modal_feedback.html`). This fragmentation leads to:
1.  **Visual Cohesion Failures**: Slight variations in padding, hover states, or color tokens across different "close" buttons.
2.  **Maintenance Friction**: Updating a design token (like changing the close button hover background) requires a multi-file sweep.
3.  **Inconsistent Accessibility**: Aria-labels and focus rings are often omitted or inconsistent in hardcoded instances.

## 2. Decision
Abstract repeating UI atoms into semantic `@layer components` in `ui/static/css/input.css` and unified Go template partials in `ui/templates/partials/ui_components.html`.

Specifically, we will implement:
-   `.btn-close`: A standardized CSS class for close buttons.
-   `{{ define "btn_close" }}`: A reusable Go template partial for the close button atom.

## 3. The Complexity Kill-Switch (Rationale)
*   **User Value**: Ensures a premium, consistent interaction model across all modals. Users experience the same hover/click feel globally.
*   **Performance Budget**: < 1ms impact (CSS abstraction is processed at build time; template partial overhead is negligible).
*   **Minimalism Check**: Replaces ~10-15 lines of redundant Tailwind strings per modal with a single template call.

## 4. Consequences
*   **Technical Tradeoffs**: Deviates slightly from "pure" utility-first Tailwind by re-introducing semantic components, but the scale of the project justifies O(1) maintenance.
*   **Observability**: Verified via `go run ./cmd/verify design` and visual regression snapshots. Future iteration of `cmd/verify/design` must be updated to flag hardcoded literal colors (`bg-earth-cream`) outside of the `input.css` boundaries.
*   **SQLite Impact**: None.

## 5. Alternatives Considered
1.  **Continue with Raw Tailwind**: Rejected because it fails the "Maintenance Friction" gate.
2.  **Headless UI Component Library**: Rejected as it adds unnecessary JavaScript complexity for simple static atoms.

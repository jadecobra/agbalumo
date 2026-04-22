# ADR 2026-04-22: Hero-Guided Viewport Clearance for Discovery Filters
**Date**: 2026-04-22 **Status**: Accepted

## 1. Context & User Problem
The primary discovery flow for the "Ada" persona (finding African food locations in < 60s) was obstructed on desktop by UI clipping. The search-adjacent filter dropdown opened downward from a vertically-centered search bar on 1440x900 displays. This caused the dropdown's bottom half—container crucial category labels like "Food"—to extend beyond the browser's lower viewport boundary. Pre-accordion fixes addressed height but failed to account for the absolute Y-positioning of the parent trigger.

## 2. Decision
Transition the Home Hero layout from **Vertical Centering** to **Top-Aligned Density**.
1.  **Repositioning**: Use `justify-start` and explicit top padding (`5vh`) on the hero section to push the search-discovery interface into the upper third of the viewport.
2.  **Constraint Hardening**: Enforce a strict `320px` max-height on the filter accordion panel with an internal scrollbar fallback.
3.  **Critical CSS Priority**: Utilize inline styles (`!important`) for these positioning fixes to bypass stylesheet caching or framework build failures in transient development environments.

## 3. The Complexity Kill-Switch (Rationale)
* **User Value**: Guaranteed visibility of all filter options without requiring primary window scrolling. This is 2x better because it restores the "One-Click Discovery" promise.
* **Performance Budget**: Zero impact; purely layout-driven.
* **Minimalism Check**: Eliminated the need for complex "open-upwards" JavaScript collision detection logic by solving the problem via structural layout density.

## 4. Consequences
* **Technical Tradeoffs**: The hero section now has more "white space" at the bottom on very tall displays, though this is mitigated by the earth-dark background aesthetics.
* **Observability**: Monitored via absolute Y-coordinate checks in the browser subagent (`Verify Y < 400`).
* **SQLite Impact**: None.

## 5. Alternatives Considered
* **Dynamic Collision Detection**: Implementing a JS observer to flip the dropdown to `bottom-full`. Rejected as too brittle and prone to event-listener regressions (e.g. HTMX swaps). Structural layout is a "dumb" (simple) and robust solution.

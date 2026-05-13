# ADR 2026-05-13: Decoupling Mobile Filter Handle & Fixing Positioning Context

**Date**: 2026-05-13 **Status**: Accepted

## 1. Context & User Problem
The mobile filter drawer exhibited two primary interaction regressions:
1.  **Click Interception**: Clicks on the mobile handle area (top of the drawer) were being incorrectly processed as "outside clicks" by the global click listener, causing the drawer to close unexpectedly when users intended to interact with top-level filters.
2.  **Positioning Context Drift**: The drawer, intended to be a `fixed` bottom-sheet, was occasionally appearing at the top of the viewport or overlapping with the header. This was traced to `backdrop-blur-md` on the sticky header creating a new stacking/positioning context, which traps `fixed` children.

## 2. Decision
1.  **Context Decoupling**: Move `backdrop-blur-md` from the `<header>` element to an absolute, non-interactive child `div`. This preserves the visual effect without creating a transform context that breaks `fixed` children.
2.  **Interaction Decoupling**: 
    - Refactor the mobile filter header to separate the **visual handle** (strictly `pointer-events-none`) from the **drawer container** (which remains `pointer-events-auto` to catch clicks and prevent closure).
    - Harden the **close button** touch target to 48x48px (`p-3`) to exceed the 44px standard and improve reliability for one-handed thumb interactions.
3.  **Layout Hardening**: Add a subtle `border-b` and explicit padding to the drawer's top section to clearly delineate the interaction layer from the scrollable filter content.

## 3. The Complexity Kill-Switch (Rationale)
- **User Value**: Restores deterministic drawer behavior and navigation parity on mobile.
- **Performance Budget**: Zero impact; purely CSS/HTML structural fixes.
- **Minimalism Check**: Avoided adding complex JavaScript collision detection or portal logic by fixing the underlying CSS stacking context issue.

## 4. Consequences
- (+) Drawer positioning is now consistently viewport-relative.
- (+) Handle area is inert but part of the drawer's "inside" hit-box, preventing accidental closures.
- (-) Minimal increase in DOM depth in the header.

## 5. Alternatives Considered
- **JavaScript Portals**: Dynamically moving the drawer to `document.body`. Rejected as it complicates the HTMX/Go Template lifecycle and breaks local search-input focus inheritance.
- **Removing Blur**: Rejected as it would violate the "Premium Design" aesthetic standards.

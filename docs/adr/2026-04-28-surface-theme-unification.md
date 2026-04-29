# UI Surface Theme Unification

## Context
The application suffered from severe visual context switching: the home hero was dark, listing cards were light, and modals were hardcoded dark. The typography was aggressively small (8px) and low contrast (60% opacity), failing basic scanning utility and WCAG standards.

## Decision
1. We enforce `text-[10px] md:text-xs` as the absolute minimum font size.
2. We map all modal surfaces to the global light/dark theme (`bg-white dark:bg-surface-dark`) to eliminate the jarring flash effect when clicking a listing.
3. We strip redundant badges from the card headers to minimize visual bloat.

## Consequences
- The application will look much more coherent and less "pieced together."
- Readability and scanning speed (the <60s goal) is significantly improved.
- Some edge-case aesthetic choices in the modal are flattened for global consistency.

# Design Critique Scoring Rubric

## Information Density (0-10)
- Start at 10
- -2 if cards NOT above fold at 1440px
- -1 per 50px of nav-to-content dead space above 100px
- -1 if grid < 3 columns at 1440px

## Action Clutter (0-10)
- Start at 10
- -2 per dead/non-functional button
- -1 per duplicate nav element across viewports
- -1 per interactive element missing data-testid

## Typography (0-10)
- Start at 10
- -0.5 per uppercase violation above threshold (>4 per file)
- -2 per font-size below 10px
- -1 per hardcoded hex code

## State Completeness (0-10)
- Start at 10
- -2 per template key gap (missing data propagation)
- -1 per missing loading/hover/empty state
- -1 per inline style attribute

## Functional Ergonomics (0-10)
- Start at 10
- -2 per CSP violation (inline handler)
- -1 per touch target < 44px on mobile
- -1 per console error

## AI Slop (0-10)
- Start at 10
- -1 per inline style attribute
- -1 per forbidden rounding class
- -1 per hardcoded modal background
- -2 if generic placeholder images detected

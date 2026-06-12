# Design Critique Scoring Rubric

> Dimensions 1-6 score aesthetics. Dimension 7 scores correctness. A template MUST pass Dimension 7 before aesthetic scores are meaningful.

## Information Density (0-10)
- Start at 10
- -2 if cards NOT above fold at 1440px
- -1 per 50px of nav-to-content dead space above 100px
- -1 if grid < 3 columns at 1440px

### In Scope
- Card position relative to the fold at 1440px.
- Nav-to-content dead vertical space.
- Grid column counts at 1440px desktop view.

### Out of Scope (Non-scoring)
- Spacing on viewport widths below 1440px (mobile/tablet layouts).

## Action Clutter (0-10)
- Start at 10
- -2 per dead/non-functional button
- -1 per duplicate nav element across viewports
- -1 per interactive element missing data-testid

### In Scope
- Dead buttons lacking actions, hrefs, or hx-attributes.
- Redundant navigation links duplicated on the same viewport.
- Missing `data-testid` on interactive elements.

### Out of Scope (Non-scoring)
- Missing `data-testid` on static header text, labels, or paragraphs.

## Typography (0-10)
- Start at 10
- -0.5 per uppercase violation above threshold (>4 per file)
- -2 per font-size below 10px
- -1 per hardcoded hex code

### In Scope
- Uppercase text transformations beyond the 4-per-file limit.
- Hardcoded CSS/HTML font sizes below 10px.
- Hardcoded raw color hex codes (e.g. `#FFFFFF`).

### Out of Scope (Non-scoring)
- Uppercase letters in acronyms or brand names (e.g., "OAuth", "API").

## State Completeness (0-10)
- Start at 10
- -2 per template key gap (missing data propagation)
- -1 per missing loading/hover/empty state
- -1 per inline style attribute

### In Scope
- Missing dict parameters when templates invoke other sub-templates.
- Missing hover feedback on buttons, skeletons for loading, or clean empty-state text.
- Inline HTML `style="..."` attributes.

### Out of Scope (Non-scoring)
- Inline style for background images (`style="background-image: ..."`), which is explicitly whitelisted.

## Functional Ergonomics (0-10)
- Start at 10
- -2 per CSP violation (inline handler)
- -1 per touch target < 44px on mobile
- -1 per console error
- -2 if a split details layout mixes data nodes (e.g., Hours, Location) into the primary actions column
- -1 if interactive buttons do not share uniform height, padding, and alignment states

### In Scope
- Inline event handlers like `onclick` or `onsubmit`.
- Mobile buttons/links smaller than 44x44px physical area.
- Layout organization of details columns and action uniformities.

### Out of Scope (Non-scoring)
- Margin spacing issues not directly causing overlap or misalignment.

## AI Slop (0-10)
- Start at 10
- -1 per inline style attribute
- -1 per forbidden rounding class
- -1 per hardcoded modal background
- -2 if generic placeholder images detected
- -1 per redundant section label (e.g., "About" for card descriptions, "Contact" for button clusters)
- -2 if action buttons/CTAs and informational details are mixed/interlaced in a vertical stack

### In Scope
- Tailwind rounding classes (`rounded-*`) used under brutalist style guidelines.
- Hardcoded modal dark backgrounds that ignore theme context.
- Generic placeholders and redundant titles.

### Out of Scope (Non-scoring)
- Whitelisted rounding classes in third-party integrations or specifically excused templates.

## Accessibility & Semantic Integrity (0-10)
- Start at 10
- -2 per icon-only button missing `aria-label` (caught by `verify design`)
- -2 per form control missing explicit `for`/`id` association (caught by `verify design`)
- -2 per `<img>` missing `alt` attribute (caught by `verify design`)
- -1 per modal missing `role="dialog"` or `aria-labelledby`
- -1 per Axe color-contrast violation reported in `a11y.spec.ts`
- **Score < 8 = BLOCKER. Do not ship regardless of aesthetic scores.**

### In Scope
- Icon-only interactive elements lacking `aria-label`.
- Form inputs without explicit label association (`for`/`id`).
- Images without `alt` attribute.
- Missing modal role declarations and contrast failures.

### Out of Scope (Non-scoring)
- `aria-hidden="true"` on decorative icon spans inside elements that already carry `aria-label`.

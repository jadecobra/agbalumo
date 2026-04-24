# UI & Template Standards

The project follows a **Global Editorial Brutalist** standard. All templates MUST adhere to these rules without exception.

## 1. Visual Geometry (Zero Radius)
- **Rounded Edges are Forbidden**: DO NOT use `rounded-md`, `rounded-lg`, etc. 
- **Sharp Corners Only**: All cards, inputs, and buttons MUST be `rounded-none`.
- **Verify**: Always run `go run ./cmd/verify design` before committing any template change.

## 2. Typography Hierarchy
- **Headings**: Use `font-serif` (Playfair Display) for headers (`h1`, `h2`). 
- **Brand Emphasis**: Use `italic` on headings for a premium editorial feel.
- **Copy**: Use `font-sans` (Inter) for functional text and body copy.

## 3. The Earth Palette
- **Foundations**: Use `bg-earth-dark` for primary containers and `text-earth-cream` for contrast.
- **Accents**: Use `earth-ochre` for CTA buttons or emphasizes.
- **Source of Truth**: Read `docs/design_system/DIALECTS.md` for the current design tokens.

## 4. Components
- **Reuse**: Use partials from `partials/ui_components.html` (e.g. `button_sharp`).
- **No Hardcoded Hex**: Use Tailwind tokens (e.g. `bg-earth-dark`) instead of literal hex values.

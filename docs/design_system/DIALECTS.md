# UI Dialect: agbalumo Editorial Brutalist

The project has converged into a single visual system. All "Brand vs Sharp" contradictions are resolved in favor of a high-contrast, sharp-edged, editorial aesthetic.

## The Core Standard (Global)

- **Typography**: 
    - **Headings**: `font-serif` (Playfair Display) + `italic` for hero/primary headers.
    - **Body/UI**: `font-sans` (Inter) for functional text.
- **Edges**: `rounded-none` (0px) is mandatory for all structural elements (cards, inputs, buttons, containers).
    - *Exception*: `rounded-full` is allowed for pill-shaped status badges or circular icon containers.
- **Palette**: **Earth System** (`bg-earth-dark`, `text-earth-cream`, `accent-earth-ochre`).
- **Depth**: Glassmorphism (`backdrop-blur-md`, `bg-white/10`) + High-contrast borders (`border-white/20`).

## Enforcement Policy

1. **Deterministic Sharpness**: The `verify design` tool enforces `0px` radius on all `.html` templates.
2. **Token Integrity**: Hardcoded hex values are strictly prohibited.
3. **No Overrides**: Avoid inline `style` with `!important` to bypass the system.

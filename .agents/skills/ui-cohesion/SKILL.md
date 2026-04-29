---
name: "UI Cohesion Guard"
description: "Enforce visual consistency, legibility, and theme parity across templates"
triggers:
  - "template_change"
  - "ui_cohesion"
  - "design review"
  - "visual audit"
mutating: false
---

# UI Cohesion Guard

## Procedural Checklist

1. **Verify Design System Compliance**
   - Run `go run ./cmd/verify design` locally before submitting any UI changes.
   - Ensure all text elements meet minimum size and contrast thresholds.

2. **Font Size Guardrails**
   - Minimum font size is `10px`.
   - Never use arbitrary Tailwind classes like `text-[8px]` or `text-[9px]`.

3. **Contrast and Opacity**
   - For subtext (`text-text-sub`), ensure opacity is at least 70% (e.g., `text-text-sub/70`, `text-text-sub/80`).
   - Opacity values below 70% (e.g., `text-text-sub/60`) violate readability standards.

4. **Dark Mode & Theme Parity**
   - Modals and persistent UI components must respect light/dark mode.
   - Do not hardcode dark themes (e.g., using raw `bg-earth-dark` without `dark:` prefix).

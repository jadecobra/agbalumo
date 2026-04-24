# UI Design Dialects

The project uses two distinct visual "dialects" to differentiate between consumer-facing brand experiences and administrative/utility tools.

## 1. Dialect: Brand (Rounded)
Used for the public homepage, discovery, and user profiles. Focused on warmth, approachability, and high-quality branding.

*   **Context**: `index.html`, `profile.html`, `about.html`.
*   **Typography**: `font-sans` (Inter) for body, `font-display` (Inter) for headings.
*   **Border Radius**: Large (`rounded-3xl` for cards/modals, `rounded-xl` for inputs).
*   **Colors**: `primary`, `secondary`, `surface-light/dark`.
*   **Shadows**: `shadow-soft`, `shadow-md`.

## 2. Dialect: Sharp (Admin/Earth)
Used for the admin dashboard, bulk actions, and utility views. Focused on density, clarity, and a "brutalist" earth palette.

*   **Context**: `admin_*.html` and templates using `earth-*` colors.
*   **Typography**: `font-serif` (Playfair Display) for headings, `font-sans` (Inter) for utility text.
*   **Border Radius**: **NONE**. Components MUST have sharp corners (no `rounded-` classes) unless using a standard small rounding like `rounded-sm` for tiny chips.
*   **Colors**: `earth-clay`, `earth-ochre`, `earth-sand`.
*   **Implementation**: Use components with the `_sharp` suffix from `ui_components.html`.

---

## Hierarchy of Truth
1.  **Tailwind Config**: `tailwind.config.js` is the final authority for token values (colors, radius sizes).
2.  **Shared Components**: `ui/templates/partials/ui_components.html` is the source of truth for component implementation.
3.  **Dialect Rules**: These rules (and localized `AGENTS.md`) determine which components and classes are valid in a given context.

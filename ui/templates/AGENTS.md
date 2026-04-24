# UI & Template Intelligence

This package controls the visual presentation via two distinct dialects: **Brand** (Consumer) and **Sharp** (Admin).

## UI Dialects
*   **Dialect: Sharp**: Applied to all `admin_*.html` templates. 
    *   **Prohibition**: Never use `rounded-` classes (3xl, xl, etc.) in this context.
    *   **Standard**: Use `_sharp` components from `ui_components.html` and `earth-*` color tokens.
*   **Dialect: Brand**: Applied to public pages (`index.html`, `profile.html`).
    *   **Standard**: Use `rounded-3xl` for surfaces and `primary/secondary` color tokens.

## Component Architecture
*   **Source of Truth**: Read `tailwind.config.js` for token values and `ui/templates/partials/ui_components.html` for existing patterns before creating new UI.
*   **DRY Templates**: Do NOT repeat HTML structures for cards, buttons, or modals. Use defined components in `ui/templates/components/`.
*   **Tailwind Only**: Prohibit the use of inline styles or custom CSS classes. Use standard Tailwind utility classes only.
*   **HTMX Fragments**: When creating HTMX-specific partials, ensure they contain exactly ONE root element.


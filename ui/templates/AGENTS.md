# UI & Template Intelligence

This package controls the visual presentation and user experience.

## Component Architecture
*   **DRY Templates**: Do NOT repeat HTML structures for cards, buttons, or modals. Use defined components in `ui/templates/components/`.
*   **Tailwind Only**: Prohibit the use of inline styles or custom CSS classes in templates. Use standard Tailwind utility classes only.
*   **HTMX Fragments**: When creating HTMX-specific partials, ensure they contain exactly ONE root element to avoid swap ambiguities.

## Discovery & Debugging
*   **Test IDs**: Add `data-testid="xxx"` to all buttons, links, and forms to ensure reliable automated testing with the `browser_subagent`.
*   **Semantic Names**: Name template files logically based on their page or purpose (e.g., `listing_card.html` not `item.html`).

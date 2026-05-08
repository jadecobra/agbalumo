# Feedback Module Context

Small, focused module for collecting user feedback and bug reports.

## Implementation Rules
*   **Minimalism**: Keep the feedback schema flat. Only capture essential fields (Category, Comment, UserID).
*   **Rate Limiting**: Apply aggressive rate limiting to feedback submission to prevent spam.
*   **HTMX Integration**: Feedback forms should submit via HTMX and return a simple "Thank You" toast or fragment.

## Logic
*   **Anonymous Feedback**: Support anonymous submissions while tagging them with session metadata if possible.

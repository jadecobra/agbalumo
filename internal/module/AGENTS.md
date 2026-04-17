# Module Layer Intelligence

This package handles HTTP routing and business logic. It is the core of the application's functionality.

## Handler Constraints
*   **Error Handling**: ALWAYS use `RespondError(c, err)` for returning errors to the client. This ensures uniform logging and user-friendly error pages.
*   **Security**: Use the `csrf` middleware for all template-rendered forms. Ensure `{{ .CSRF }}` is included in the HTML input.
*   **Context usage**: Do not store business logic state in the Echo context. Use the `service` layer for all heavy lifting.

## Logic & Services
*   **Validation**: Perform input validation in the handler before passing data to the service. Use struct tags and the project's validator if available.
*   **Indemnity**: Handlers should be as "thin" as possible. If a handler exceeds 50 lines, extract the logic into a private helper or a service method.

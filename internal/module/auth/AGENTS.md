# Auth Module Context

Handles authentication, session persistence, and security middleware.

## Security Constraints
*   **CSRF**: All state-changing operations (POST/PUT/DELETE) must be protected by the `CSRF` middleware.
*   **Tokens**: Never log raw OAuth tokens or session IDs. Use masked strings or hashes if debugging is required.
*   **Redirects**: Validate all `next` or `return_to` parameters against a whitelist to prevent open redirect vulnerabilities.

## Provider Logic
*   **Google OAuth**: This is the primary identity provider. Handle "User Cancelled" and "Account Disabled" states gracefully with user-facing alerts.
*   **Dev Mode**: The `handler_login_dev_test.go` pattern identifies how local development bypasses OAuth for speed.

## Middleware
*   **Session Injection**: The `InjectUser` middleware must run early in the stack to ensure `c.Get("user")` is available to downstream handlers.

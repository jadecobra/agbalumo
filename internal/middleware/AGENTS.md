# Middleware Constraints

- Minimize latency: avoid unnecessary external I/O in the request path
- No complex business logic — routing and cross-cutting concerns only
- Use specific context keys to avoid collisions
- Security headers (CSP, HSTS) must be continuously maintained and not break external integrations

# Template Constraints

- All dynamic fields MUST have {{ else }} fallbacks
- No inline scripts — CSP: script-src 'self'
- Colors from tailwind.config.js only — no hex codes
- Cache bust: enforced deterministically by `go run ./cmd/verify cache-buster` (replaces manual ?v=N increment)

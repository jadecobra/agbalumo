# Template Constraints

- All dynamic fields MUST have {{ else }} fallbacks
- No inline scripts — CSP: script-src 'self'
- Colors from tailwind.config.js only — no hex codes
- Cache bust: increment ?v=N in head_meta.html after changes

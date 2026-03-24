# Claude Instructions

All core rules are in `AGENTS.md`. Follow them exactly.

## Claude-Specific Overrides
- Follow the `/feature-implementation` workflow (`.agents/workflows/feature-implementation.md`) for all features.
- Use `browser_subagent` (read the testing URL from `.agents/rules/browser-url.md`) to verify UI changes after tests pass.
- Update `docs/spec.md` after changes and tests pass.
- NEVER lower the coverage threshold in `.agents/coverage-threshold`.

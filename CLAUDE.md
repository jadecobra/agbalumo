# Claude Instructions

All core rules are in `AGENTS.md`. Follow them exactly.

## Claude-Specific Overrides
- Follow the `/feature-implementation` workflow (`.agent/workflows/feature-implementation.md`) for all features.
- Use `browser_subagent` to verify UI changes after tests pass.
- Update `docs/spec.md` after changes and tests pass.
- NEVER lower the coverage threshold in `.agent/coverage-threshold`.

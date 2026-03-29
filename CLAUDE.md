# Claude Instructions

All core rules are in `AGENTS.md`. Follow them exactly.

## Claude-Specific Overrides
- Follow the **`/build-feature`** workflow (`.agents/workflows/build-feature.md`) for all feature development.
- **Consult the Squad**:
    - Consult the **ProductOwner** for any changes affecting user value or cultural context.
    - Consult the **SystemsArchitect** for any core API, DB schema, or security changes.
- Use `browser_subagent` (read the testing URL from `.agents/rules/browser-url.md`) to verify UI changes after tests pass.
- Update `docs/spec.md` after changes and tests pass.
- NEVER lower the coverage threshold in `.agents/coverage-threshold`.
- **UI Standards**: Strictly adhere to the **WCAG** accessibility and **8pt grid** design system mandated by the **UIUXDesigner**.

---
description: Coding Standards and Guidelines (Go, HTMX, Tailwind)
---

# Coding Standards & Guidelines

This workflow is referenced for specific edge-case rules regarding code style, imports, error handling, and testing. Use this as a reference when writing new features.

## Code Style Guidelines

### Imports (Group Order)
1. Standard library (blank line)
2. Third-party packages (blank line)
3. Local packages

### Naming Conventions
- Packages: lowercase single word (`domain`, `handler`, `service`)
- Types: PascalCase (`ListingHandler`, `UserStore`)
- Interfaces: `XxxStore`, `XxxService` (`ListingStore`)
- Errors: `ErrSomething` as package-level vars (`ErrInvalidDeadline`)

### Structs & Types
- Use struct tags: `json`, `form` for API binding
- Constructor pattern: `NewXxx()` functions

### Error Handling
- Use `RespondError(c, err)` for HTTP handlers - logs internally, renders friendly error page
- Wrap with `echo.NewHTTPError` for specific codes.
- Handlers should generally return `error`.

## Context Cost Awareness (Tokens)

To maintain agentic efficiency, we monitor **Token Density**.
- **Advisory TokenRMS**: Target **< 110**.
- **Context Window**: Monitor `ContextWindowPct` relative to Claude Sonnet (200k tokens).
- **Efficiency Pattern**: If a file exceeds **500 tokens**, consider if splitting into sub-packages or smaller files would improve logical cohesion and "Agentic Attention."
- **Janitor Run**: Use `/janitor` to clean up high-cost or high-entropy files when the TokenRMS exceeds thresholds significantly.
- **Verification**: Verify token density at any time by running: `go run ./cmd/verify context-cost`.

## Security & Linter Suppression

- **Suppression Justification:** Any time you use a `#nosec` directive, you must include a valid justification comment. Validate this project constraint by running `go run ./cmd/verify gosec-rationale`.

## Testing Conventions

### Test Structure
Use table-driven tests:
```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        expectErr bool
    }{
        {name: "valid case", input: "value", expectErr: false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) { ... })
    }
}
```

### Mocks
Use `github.com/stretchr/testify/mock`. Place mocks in `internal/mock/`.

### Test Helpers & Anti-Duplication
- **Use Standards First:** Always check `internal/testutil/` for centralized UI, authentication, or seeding helpers (like `setupAdminTestContext`, `testutil.NewMainTemplate`) before creating custom ones inline.
- **Prevent Copy-Paste:** To avoid failing the `cmd/verify/main.go critique` toolchain, any repetitive setup boilerplate or repeated logic between subtests *must* be extracted into explicit helper functions.
- **Incremental Auditing:** To reduce noise from legacy technical debt, day-to-day development should use incremental auditing (`go run ./cmd/verify critique`). This only flags violations introduced in the current branch.
- **Auto-Healing:** Structural maintenance (like struct field alignment) should be automated. Run `go run ./cmd/verify heal` to resolve these issues automatically before committing.
- **Cognitive Complexity limits:** High-complexity test functions (with deep nesting or serial sequential assertions) will fail the project's quality gates. Extract large blocks into private helpers.

### Coverage & Rules
- **TDD:** Write failing tests (RED) FIRST.
- **Coverage:** Thresholds are defined in `.agents/coverage-thresholds.json`. NEVER lower these value.
- **Persona Sync:** Changes to rules MUST be mirrored in `GEMINI.md` and `.agents/rules/`.
- **Functions:** Keep functions small and single-purpose (SRP).
- **Comments:** Code should be self-documenting; avoid unnecessary comments.

> **UI Component Constraint**: Do NOT write raw HTML form, button, or modal elements directly into major page templates. You MUST utilize or propose existing component definitions located inside `ui/templates/components/`. All colors, spacing, and visual effects must be derived exclusively from `tailwind.config.js` logic natively—no arbitrary hex codes or custom CSS outside of the Tailwind engine.

## Strict Lessons

This section contains corrections and constraints derived from the `[/learn]` workflow. These rules take precedence over existing style guidelines.

* The agent MUST strictly enforce structural optimizations and code duplication limits by ensuring `cmd/verify/main.go critique` (incremental) is run during development and `cmd/verify/main.go ci` (full) is run before pushing.
* The agent SHOULD utilize `go run ./cmd/verify heal` to automatically resolve structural maintenance issues (like field alignment) before attempting manual refactoring.
* The agent MUST prioritize resolving "Code Duplication" (clone groups) over other technical debt from `critique`, as duplicated logic has the largest negative impact on agent context and maintainability. When addressing this debt, the agent MUST explicitly compare the total number of clone groups reported by `go run ./cmd/verify critique` before and after fixes are implemented to guarantee quantitative improvement.
* The agent MUST verify all GitHub Action SHAs using the `verify-shas` tool. Remote verification via `gh` CLI is a **local-only hard gate** and MUST be executed locally before pushing any changes to CI configurations. This check is omitted from production CI to avoid unnecessary secret dependencies and because GitHub's infrastructure natively validates SHAs during job initialization.
* **Gate Relocation Lesson**: Maintenance tools requiring third-party authentication or remote API access (e.g., `gh api`) should be prioritized as local verification gates rather than production CI blocks, unless the check is globally critical and has no local equivalent. This ensures development velocity while maintaining consistency across environments.
* **Server Verification Lesson**: The agent MUST always verify the active listener port (e.g., `https://localhost:8443` vs `http://localhost:8080`) from logs or `cmd/serve.go` logic before initiating browser verification, to avoid connection failures in secure or non-standard environments.
* **Environment Constant Lesson**: The agent MUST use `domain.EnvKeyAppEnv` and other domain-defined constants for environment variable keys instead of hardcoded strings (like `APP_ENV`), to maintain consistency with the project's `.env` schema.
* **CI Infrastructure Linting Lesson**: The agent MUST prioritize verifying CI YAML configuration (`.github/workflows/ci.yml`) as part of the local `audit` or `lint` phase. This ensures that infrastructure drift, such as invalid parameter names in GitHub Actions (e.g., `trivy-version` vs `version`), is caught locally before being pushed to production. The local verification suite is the source of truth for CI parameter validity.
* **Initialization Guard Lesson**: UI components MUST NOT initialize themselves via `DOMContentLoaded` listeners if they are already registered in the centralized `initApp()` sequence in `app.js`. This prevents double-binding the same event listeners to the document, which can cause interactions (e.g., toggles) to effectively cancel themselves out.
* **Static Asset Cache Invalidation Lesson**: When applying critical UI or functional fixes to JS or CSS files that are subject to aggressive caching (e.g., `Cache-Control: immutable`), the agent MUST increment the version query parameter (`?v=N`) in `ui/templates/components/head_meta.html`. This guarantees the fix is delivered to all clients immediately upon server restart.
* **JavaScript Reliability & Security Lesson**: The agent MUST ensure all client-side JavaScript assets are syntactically validated in the CI pipeline using `go run ./cmd/verify js-syntax`. Direct inline scripts MUST be avoided in templates to maintain a strict `script-src 'self'` Content Security Policy, preventing the risk of XSS and execution failures in hardened production environments.
* **UI Regression Verification Lesson**: The agent MUST always perform end-to-end verification of UI-facing regressions (like search filters or category switching) using the browser subagent. Relying solely on repository-level unit tests is insufficient to guarantee that the user's intent (e.g., Ada's 60-second discovery flow) is fully satisfied across the frontend stack.
* **Local Dynamic Audit Gate**: CLI-based CI pipelines MUST include a live server-startup check (`maintenance.VerifyServerStartup`) to catch nil-pointer regressions in routing and dependency injection layers that unit tests may miss.
* **Environment Parity**: Always verify CI failures in production by matching the exact environment variables (e.g., `AGBALUMO_ENV=development`) used in the remote workflow.
* **CI Orchestration Audit Lesson**: When refactoring core internal CLI entry points (e.g., `cmd/verify`), the agent MUST autonomously audit `.github/workflows/*.yml` and all local automation scripts (e.g., `setup-hooks.sh`) to ensure orchestration logic is updated to the new syntax. This prevents production pipeline failures caused by legacy path assumptions.
* **Production Failure Verification Lesson**: Before attempting to fix a production CI failure, the agent MUST explicitly verify the failure by viewing the actual run logs (e.g., using `gh run view --log`). Deducing the failure from local state is insufficient and prone to targeting the wrong symptoms or missing environmental nuances.
* **Post-Push Production Monitoring**: After successfully pushing changes to the remote repository, the agent MUST monitor the production CI/CD pipeline using the GitHub CLI (e.g., `gh run watch`). This ensures that the code passes all remote-only infrastructure gates and maintains a 'green' stable branch. The turn is not complete until production CI stability is confirmed.
* **Native UI Accordion Lesson**: The agent MUST prioritize native HTML5 tags (e.g., `<details>` and `<summary>`) for accordions and expansion UI. This ensures baseline functionality (expanding/collapsing) remains bulletproof regardless of JavaScript state or HTMX context.
* **UI Template Robustness Lesson**: UI templates MUST include fallback text for all dynamic data fields (e.g., using `{{ if .Field }}{{ .Field }}{{ else }}Fallback{{ end }}`) to prevent "invisible" or "empty" elements from breaking the layout or UX.
* **Browser Verification Deep-Dive**: Browser-based validation MUST NOT only check for the existence of an element but MUST explicitly verify visibility (`isVisible`), `innerText` length, and state transitions (e.g., "arrow rotated successfully") to catch silent regressions in rendering and interactivity.
* **Idempotent Global Listeners**: Global UI event listeners (e.g., "click-outside" or "Escape key" close handlers) MUST be bound to the `document` exactly once using an idempotent flag (e.g., `window._is_bound`) or a centralized `initApp` sequence that specifically prevents duplicate bindings during HTMX swaps.

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

## Product-Centric Performance (The 60-Second Goal)

The "North Star" for this project is for a user to find African food in any city in under 60 seconds.

### Complexity Kill-Switch
You MUST justify the existence of any new feature or abstraction. If it increases UI friction or DB latency without a 2x increase in user utility, it must be deleted or simplified.

### Bottleneck-Aware Growth
Features designed solely for user acquisition (sharing, referrals, social loops) are considered UI Bloat and MUST be rejected until listing quality (Accuracy/Verification) is no longer the primary bottleneck.

### Performance Budget
The target Time to First Result (TTFR) is < 500ms on a standard mobile connection.

### Latency Guardrail
Any change estimated to add >100ms to the critical search path requires a formal ADR and a justification of why the "User Value" outweighs the speed penalty.

## Data Integrity & Trust Mandate

Speed is irrelevant if the data is wrong. Trust is our most expensive asset.

### The Hours-to-Pulse Pipeline
We do not call blind. To minimize nuisance and maximize "Proof of Life" accuracy:
- **Scraper-First Hours**: The Menu URL scraper MUST prioritize extracting "Hours of Operation" text from the primary Menu URL/Official Website.
- **LLM Extraction**: Use a lightweight LLM prompt to normalize messy "Hours" text into a standard JSON schedule.
- **The Scheduler**: The "Phone Pulse" system MUST only initiate calls during the extracted "Open" windows.
- **The No-Vision Rule**: If hours are embedded in images/flyers, do NOT use OCR/Vision extraction. This is considered Complexity Creep. Fall back immediately to the "Phone Pulse."
- **Zero-Data Fallback**: If no hours are found via text scraping, use a global "Safe Window" (1 PM - 6 PM local time) for the first "Phone Pulse."

### The Escalation Pulse (NLP Curation)
If the primary site scraper fails to find hours, the bot script graduates to a curation tool.
- **The NLP Script**: "Hi, I'm from Agbalumo. We couldn't find your hours online—what are your opening hours today?"
- **Ambiguity Handling (The "Honest Failure" Rule)**: If the LLM parser confidence score is low (e.g., < 0.8), mark as "Help Us Verify".
- **UI Implementation**: Display a "We tried to verify hours but weren't 100% sure. Can you help us?" prompt on the listing.

### Zero-Cognitive-Load Curation
- **The Single-Tap Rule**: Interaction for "Help Us Verify" must be a binary confirmation (e.g., "Are they open right now? [Yes] [No]").

### Existence Verification & Proxy Signals (The Phone Pulse Protocol)
- **Frequency**: Limit successful pulses to once every 14 days.
- **Success Definition**: Human or IVR pickup counts as success.
- **Multi-Day Retry Logic**: Soft failures (Busy/No Answer) require three (3) retry attempts on different days and windows within a 1-week period.
- **Hard Failure Action**: If all three multi-day retries fail, immediately flag as "Menu Unavailable" and deprioritize.

### Automated Trust Scoring (Verified Badge)
A listing is "Verified" if it has:
- **Freshness**: Successful "Proof of Life" signal within the last 7 days.
- **Consistency**: Zero "Broken Link" or "Closed" reports in last 30 days.
- **Completeness**: Valid address, phone, and verified hours.
- **Partial Failure Honesty**: If a critical data point (like a Menu URL) is broken but the restaurant is verified open, DO NOT hide the restaurant. Display the listing with a clear "Menu Unavailable" status.

### Scaling Skepticism (Conflict Resolution)
If a user-tap conflicts with a recent system verification:
- **Early Stage (Full Trust)**: While users < 100, the user-tap takes immediate precedence.
- **Growth Stage (Skeptical)**: Once users > 100, require a threshold ($N > 1$) before overriding.

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

### CI & Infrastructure
* **Structural Optimization Lesson** [TRIGGER: refactoring, git_push]: The agent MUST strictly enforce structural optimizations and code duplication limits by ensuring `cmd/verify/main.go critique` (incremental) is run during development and `cmd/verify/main.go ci` (full) is run before pushing.
* **Action SHA Verification Lesson** [TRIGGER: ci_config_change]: The agent MUST verify all GitHub Action SHAs using the `verify-shas` tool. Remote verification via `gh` CLI is a **local-only hard gate** and MUST be executed locally before pushing any changes to CI configurations. This check is omitted from production CI to avoid unnecessary secret dependencies and because GitHub's infrastructure natively validates SHAs during job initialization.
* **Gate Relocation Lesson** [TRIGGER: ci_config_change]: Maintenance tools requiring third-party authentication or remote API access (e.g., `gh api`) should be prioritized as local verification gates rather than production CI blocks, unless the check is globally critical and has no local equivalent. This ensures development velocity while maintaining consistency across environments.
* **CI Infrastructure Linting Lesson** [TRIGGER: ci_config_change]: The agent MUST prioritize verifying CI YAML configuration (`.github/workflows/ci.yml`) as part of the local `audit` or `lint` phase. This ensures that infrastructure drift, such as invalid parameter names in GitHub Actions (e.g., `trivy-version` vs `version`), is caught locally before being pushed to production. The local verification suite is the source of truth for CI parameter validity.
* **CI Orchestration Audit Lesson** [TRIGGER: refactoring, ci_config_change]: When refactoring core internal CLI entry points (e.g., `cmd/verify`), the agent MUST autonomously audit `.github/workflows/*.yml` and all local automation scripts (e.g., `setup-hooks.sh`) to ensure orchestration logic is updated to the new syntax. This prevents production pipeline failures caused by legacy path assumptions.
* **Production Failure Verification Lesson** [TRIGGER: git_push]: Before attempting to fix a production CI failure, the agent MUST explicitly verify the failure by viewing the actual run logs (e.g., using `gh run view --log`). Deducing the failure from local state is insufficient and prone to targeting the wrong symptoms or missing environmental nuances.
* **Post-Push Production Monitoring** [TRIGGER: git_push]: After successfully pushing changes to the remote repository, the agent MUST monitor the production CI/CD pipeline using the GitHub CLI (e.g., `gh run watch`). This ensures that the code passes all remote-only infrastructure gates and maintains a 'green' stable branch. The turn is not complete until production CI stability is confirmed.
* **Environment Parity** [TRIGGER: env_variable, git_push]: Always verify CI failures in production by matching the exact environment variables (e.g., `AGBALUMO_ENV=development`) used in the remote workflow.
* **Local Dynamic Audit Gate** [TRIGGER: server_url, ci_config_change]: CLI-based CI pipelines MUST include a live server-startup check (`maintenance.VerifyServerStartup`) to catch nil-pointer regressions in routing and dependency injection layers that unit tests may miss.

### UI & Frontend
* **Initialization Guard Lesson** [TRIGGER: js_or_css_change]: UI components MUST NOT initialize themselves via `DOMContentLoaded` listeners if they are already registered in the centralized `initApp()` sequence in `app.js`. This prevents double-binding the same event listeners to the document, which can cause interactions (e.g., toggles) to effectively cancel themselves out.
* **Static Asset Cache Invalidation Lesson** [TRIGGER: js_or_css_change, template_change]: When applying critical UI or functional fixes to JS or CSS files that are subject to aggressive caching (e.g., `Cache-Control: immutable`), the agent MUST increment the version query parameter (`?v=N`) in `ui/templates/components/head_meta.html`. This guarantees the fix is delivered to all clients immediately upon server restart.
* **UI Regression Verification Lesson** [TRIGGER: browser_subagent]: The agent MUST always perform end-to-end verification of UI-facing regressions (like search filters or category switching) using the browser subagent. Relying solely on repository-level unit tests is insufficient to guarantee that the user's intent (e.g., Ada's 60-second discovery flow) is fully satisfied across the frontend stack.
* **Explicit UX Verification Lesson** [TRIGGER: browser_subagent]: For any UI-facing change, the agent MUST explicitly document the browser subagent verification steps taken (e.g., "Verified that the filter dropdown closes when clicking outside the container and that the results list updates within 500ms").
* **Native UI Accordion Lesson** [TRIGGER: template_change]: The agent MUST prioritize native HTML5 tags (e.g., `<details>` and `<summary>`) for accordions and expansion UI. This ensures baseline functionality (expanding/collapsing) remains bulletproof regardless of JavaScript state or HTMX context.
* **UI Template Robustness Lesson** [TRIGGER: template_change]: UI templates MUST include fallback text for all dynamic data fields (e.g., using `{{ if .Field }}{{ .Field }}{{ else }}Fallback{{ end }}`) to prevent "invisible" or "empty" elements from breaking the layout or UX.
* **Browser Verification Deep-Dive** [TRIGGER: browser_subagent]: Browser-based validation MUST NOT only check for the existence of an element but MUST explicitly verify visibility (`isVisible`), `innerText` length, and state transitions (e.g., "arrow rotated successfully") to catch silent regressions in rendering and interactivity.
* **Idempotent Global Listeners** [TRIGGER: js_or_css_change]: Global UI event listeners (e.g., "click-outside" or "Escape key" close handlers) MUST be bound to the `document` exactly once using an idempotent flag (e.g., `window._is_bound`) or a centralized `initApp` sequence that specifically prevents duplicate bindings during HTMX swaps.
* **UI Interaction Depth Lesson** [TRIGGER: browser_subagent]: The agent MUST explicitly verify interactive UX behaviors (e.g., scrollability of long lists, dismissal on outside click, mobile bottom-sheet transitions) using a browser subagent for any UI-facing change. Relying on "element existence" or purely logical code reviews is insufficient; the agent MUST provide proof of functional verification (e.g., "confirmed dropdown closes on background click") to avoid premature completion reports for non-trivial UI features.
* **Dropdown Viewport Clearance Lesson** [TRIGGER: ui_positioning, browser_subagent]: When implementing absolute-positioned UI components (e.g., dropdowns, modals), the agent MUST explicitly verify vertical clearance at standard resolutions (1440x900). If a component is likely to collide with the viewport edge, the agent MUST proactively reposition the parent container or implement 'open-upwards' logic to ensure the entire component is accessible without user scrolling of the main window.
* **Critical Layout Hardening Lesson** [TRIGGER: ui_positioning, js_or_css_change]: For UI fixes that resolve visual clipping or positioning failures, the agent MUST use **Inline Styles** (`style="..."`) with `!important` as a redundant fallback. This ensures the fix is applied immediately, bypassing potential stylesheet caching or framework build failures (e.g., missing Tailwind JIT utilities) in transient development environments.

### Security & Environment
* **Server Verification Lesson** [TRIGGER: browser_subagent, server_url]: The agent MUST always verify the active listener port (e.g., `https://localhost:8443` vs `http://localhost:8080`) from logs or `cmd/serve.go` logic before initiating browser verification, to avoid connection failures in secure or non-standard environments.
* **Environment Constant Lesson** [TRIGGER: env_variable]: The agent MUST use `domain.EnvKeyAppEnv` and other domain-defined constants for environment variable keys instead of hardcoded strings (like `APP_ENV`), to maintain consistency with the project's `.env` schema.
* **HTTPS Awareness Lesson** [TRIGGER: browser_subagent, env_variable]: The agent MUST always check `.env` for the `BASE_URL` to identify the correct protocol (HTTPS) and port (e.g., 8443) before initiating browser subagent tasks. Assuming `http://localhost:8080` is a failure of environment awareness.
* **JavaScript Reliability & Security Lesson** [TRIGGER: js_or_css_change, ci_config_change]: The agent MUST ensure all client-side JavaScript assets are syntactically validated in the CI pipeline using `go run ./cmd/verify js-syntax`. Direct inline scripts MUST be avoided in templates to maintain a strict `script-src 'self'` Content Security Policy, preventing the risk of XSS and execution failures in hardened production environments.
* **Scratch Directory Isolation Lesson** [TRIGGER: refactoring, git_push]: The agent MUST ensure that any temporary 'scratch' or 'brain' directories used for internal processing are explicitly added to `.gitignore` and NEVER committed to the repository. These directories often contain semi-structured Go files and data that can interfere with production CI tests and linter phases, leading to false-positive failures on the remote branch.

### Testing
* **Test Parallelism Safety Lesson** [TRIGGER: test_authoring, env_variable]: The agent MUST NOT use `t.Parallel()` in Go tests that modify global state, specifically environment variables via `os.Setenv` or `os.Unsetenv`. Doing so causes flakey CI failures that are difficult to reproduce locally but consistently fail under high remote concurrency. Non-isolated tests MUST run sequentially to ensure environment integrity.
* **Code Duplication Lesson** [TRIGGER: refactoring, test_authoring]: The agent MUST prioritize resolving "Code Duplication" (clone groups) over other technical debt from `critique`, as duplicated logic has the largest negative impact on agent context and maintainability. When addressing this debt, the agent MUST explicitly compare the total number of clone groups reported by `go run ./cmd/verify critique` before and after fixes are implemented to guarantee quantitative improvement.
* **Heal Automation Lesson** [TRIGGER: refactoring]: The agent SHOULD utilize `go run ./cmd/verify heal` to automatically resolve structural maintenance issues (like field alignment) before attempting manual refactoring.


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

### Truth over Feature (Unreliable Live Signals)
**[TRIGGER: TIME_SENSITIVE_SIGNAL]**: Never surface real-time signals (e.g., "Open Now", "Closing Soon") if the underlying data lacks the necessary context (e.g., per-listing Time Zones) to guarantee 100% accuracy.
- **Decision**: If a signal cannot be verified as accurate due to environmental gaps, it must be removed from the UI entirely to preserve user trust.

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
Use `github.com/stretchr/testify/mock`. Place mocks in `internal/testutil/`.

### Test Helpers & Anti-Duplication
- **Use Standards First:** Always check `internal/testutil/` for centralized UI, authentication, or seeding helpers (like `setupAdminTestContext`, `testutil.NewMainTemplate`) before creating custom ones inline.
- **Prevent Copy-Paste:** To avoid failing the `cmd/verify/main.go critique` toolchain, any repetitive setup boilerplate or repeated logic between subtests *must* be extracted into explicit helper functions.
- **Incremental Auditing:** To reduce noise from legacy technical debt, day-to-day development should use incremental auditing (`go run ./cmd/verify critique`). This only flags violations introduced in the current branch.
- **Auto-Healing:** Structural maintenance (like struct field alignment) should be automated. Run `go run ./cmd/verify heal` to resolve these issues automatically before committing.
- **Cognitive Complexity limits:** High-complexity test functions (with deep nesting or serial sequential assertions) will fail the project's quality gates. Extract large blocks into private helpers.

### Coverage & Rules
- **TDD:** Write failing tests (RED) FIRST.
- **Coverage:** Thresholds are defined in `.agents/coverage.json`. NEVER lower these value.
- **Persona Sync:** Changes to rules MUST be mirrored in `GEMINI.md` and `.agents/skills/`.
- **Functions:** Keep functions small and single-purpose (SRP).
- **Comments:** Code should be self-documenting; avoid unnecessary comments.

> **UI Component Constraint**: While Tailwind utility classes are preferred for layouts, heavily repeated UI atoms (e.g., close buttons, modal backdrops) MUST be abstracted into `@layer components` in `input.css` or unified Go templates (`ui_components.html`). Do not use raw, fragmented utility strings for universal elements.

## Strict Lessons

This section contains corrections and constraints derived from the `[/learn]` workflow. These rules take precedence over existing style guidelines.

### Quota & Workflow Hardening
* **Meta vs Product Quota Boundary** [TRIGGER: replace_file_content, multi_replace_file_content, /build-feature]: High-tier reasoning models (Pro/Opus) are restricted to Meta-Work (rules, architecture, design critique) and architectural planning. Iterative product execution (Go/UI TDD loops) MUST be delegated to Gemini 3 Flash via `/Flash Planning` to optimize cost and performance. High-tier models MUST halt and prompt for delegation when encountering product code execution unless an explicit `OVERRIDE` is provided.
* **Post-Push Production Monitoring** [TRIGGER: git_push]: The agent is FORBIDDEN from marking a task as "done" or ending a session based solely on a successful `git push`. After pushing, the agent MUST immediately execute `./scripts/pushw.sh` or `gh run watch` to monitor the remote CI pipeline. Work is only considered complete when the remote branch is 'green'.
* **Post-Action Verification Mandate** [TRIGGER: session_done]: When implementing a new verification tool or rule (via `/learn` or `/build-feature`), the agent MUST immediately use that tool/rule to verify its own state in the current turn.


### UI & Frontend
* **UI Blindspot Protocol (Mandatory Local E2E)** [TRIGGER: ui_change, template_change, css_change]: `go run ./cmd/verify preflight` intentionally skips Playwright tests for speed. Therefore, any task mutating HTML, CSS, or UI templates MUST explicitly execute `go run ./cmd/verify browser` locally to catch visual and interaction bugs before pushing.
* **Semantic UI Atoms** [TRIGGER: ui_change, modal_creation]: The agent is FORBIDDEN from hardcoding raw utility classes for universal atoms (close buttons, toast wrappers, standard surfaces). You MUST use predefined partials (e.g., {{ template "btn_close" }}) or semantic CSS classes.
* **UI Interactivity & Performance Depth** [TRIGGER: browser_subagent, ui_change]: Browser verification MUST NOT only check for element existence. It MUST verify state transitions (e.g., arrow rotation, modal dismissal), interaction depth (e.g., scrollability, mobile bottom-sheet transitions), and performance impact. Any asset degrading LCP by >100ms is a visual bug and must be lazy-loaded or optimized.
* **Idempotency & Duplicate branding Guard** [TRIGGER: js_or_css_change, template_change]: Prevent duplicate event bindings and redundant branding by using idempotent initialization flags (e.g., `window._is_bound`) and auditing template block boundaries. UI components MUST NOT initialize via `DOMContentLoaded` if they are registered in the centralized `initApp()` sequence.
* **UI Dialect & Positioning Awareness** [TRIGGER: ui_positioning, ui_dialect]: Distinguish between **Brand** (Rounded) and **Sharp** (Admin) dialects. **Sharp** templates prohibit rounding classes. When positioning overlays, verify vertical clearance at 1440x900 to prevent viewport collisions.
* **Template Robustness** [TRIGGER: template_change]: UI templates MUST include fallback text for all dynamic data fields to prevent "invisible" elements from breaking layout.
* **HTMX Micro-Interaction Guard** [TRIGGER: htmx]: Mandate `hx-indicator` (skeletons/spinners) and transition classes on all mutations to ensure fluid perceived performance.




### Infrastructure & Environment
* **Environment & Protocol Awareness** [TRIGGER: browser_subagent, env_variable]: The agent MUST use domain-defined constants (e.g., `domain.EnvKeyAppEnv`) for env keys. Before initiating browser tasks, construct the target URL from `.agents/invariants.json` (`protocol` + `port`) — this is the committed, canonical source of truth derived from BASE_URL. Do NOT read `.env` for URL construction (security risk — `.env` may contain secrets).
* **Port Guessing Prohibition** [TRIGGER: browser_subagent, url_construction]: The agent is STRICTLY FORBIDDEN from guessing or hardcoding port numbers (e.g., assuming port `8080`). The canonical port is in `.agents/invariants.json`. Guessing a port is equivalent to ignoring available documentation and will cause a production CI failure.
* **Local Server & Audit Gates** [TRIGGER: session_done, ci_config_change]: The agent is FORBIDDEN from ending a session if the local server is down (verify via `uptime`). CLI-based CI pipelines MUST include a live server-startup check to catch nil-pointer regressions in routing.
* **Scratch Directory Isolation** [TRIGGER: git_push]: Ensure temporary 'scratch' or 'brain' directories are in `.gitignore` and NEVER committed, as they interfere with remote CI linter phases.


### Testing
* **Flaky Test Eradication (Performance Assertions)** [TRIGGER: test_authoring, flake_fix]: Pass/Fail unit tests MUST NOT contain hardcoded latency budgets (e.g., `assert.Less(duration, 1000ms)`). Performance constraints should be measured in dedicated benchmark suites (`go test -bench`), not in standard unit tests that block standard execution paths.
* **CI Cost Awareness (Matrix Bloat)** [TRIGGER: e2e_change, playwright_config]: Agents are forbidden from adding new viewports or matrix dimensions to CI pipelines (e.g., Playwright `projects`) without explicitly calculating the CI time cost and implementing a parallelization strategy (e.g., sharding or increasing workers).
* **Test Parallelism & Env Safety** [TRIGGER: test_authoring]: Use `t.Setenv()` for setting env vars in tests. Raw `os.Setenv` without cleanup is forbidden as it causes flaky CI failures across unrelated test files.
* **Code Duplication & Critique** [TRIGGER: refactoring]: Prioritize resolving "Code Duplication" (clone groups) from `critique`. Quantify improvements by comparing clone group counts before and after refactoring.
* **E2E & Visual Regression** [TRIGGER: visual_testing, htmx]: Playwright tests MUST use relative paths. HTMX assertions MUST use `page.waitForResponse()` to avoid timing flakes. New modals MUST include visual snapshots in `visual.spec.ts`.
* **Visual Snapshot Platform Parity** [TRIGGER: snapshot_parity]: Visual snapshots MUST include linux-suffixed baselines generated via Docker. satisfy parity via `go run ./cmd/verify snapshot-parity`.
* **Full E2E Mandate — No Filtered Runs Before Push** [TRIGGER: e2e_change, ui_change, git_push]: The agent is FORBIDDEN from running a filtered/scoped Playwright invocation (e.g., `--grep`, `-g`, specific file) as a substitute for the full suite before pushing. If `go run ./cmd/verify browser` fails, the agent MUST (1) report the full failure output to the user, (2) halt and await direction — it MUST NOT silently run a subset to "verify the fix" and push.
* **AXE Accessibility Platform Parity Gate** [TRIGGER: a11y_change, ui_change, e2e_change]: The local Darwin environment produces different AXE contrast evaluations than the Linux CI container due to sub-pixel rendering and browser cache warm state. Any change touching color tokens, navigation, or card components MUST be validated by running the full `go run ./cmd/verify browser` suite in a cold-state browser (clear cache/incognito) before pushing. A local green result is NOT sufficient evidence of CI-green for accessibility checks.

 
### Scraping
* **Platform Attribution Exclusion** [TRIGGER: scraper_change]: Heuristic scrapers MUST explicitly exclude platform-level marketing and attribution links (e.g., "Powered by ZingMyOrder") which keyword-match but point to irrelevant content.


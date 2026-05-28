As the **Senior Product Engineer**, your mission is to build ruthlessly simple, high-utility systems that solve user problems (e.g., finding African food in < 60 seconds at agbalumo.com). You prioritize **User Value** and **Minimal Latency** over architectural purity.
## THE LEARNING MANDATE
You are forbidden from letting a mistake (technical or product) go unrecorded.
- **Complexity Kill-Switch**: For strategic commands (`/build-feature`, `/architect`), if a feature adds UI steps or DB latency without a 2x increase in utility, challenge the user. For tactical/surgical commands (`/hotfix`, `/refactor`, `/debug`), skip this check to optimize speed.
- **Performance Budget**: Every feature must justify its impact on search latency. If it fails the **60-second find goal**, suggest a `/learn` session.
- **Anti-Sycophancy & Socratic Pushback (Strategic Only)**: During planning or strategic features, present a structured **Cons & Alternatives Matrix** in a single turn instead of sequential dialog pushback. Skip pushback entirely for tactical `/hotfix`, `/refactor`, and `/debug` commands to allow frictionless execution of specific user instructions.
- **Cost of Action Projection**: For major architectural changes, explicitly state what existing sub-system is most likely to break under 10x scale.
## COMMUNICATION & TONE
Act as a terse, highly technical Senior Staff Engineer pair-programming with a peer.
- **Zero Fluff**: No pleasantries, no apologies, no generic introductions or conclusions. Get straight to the technical point.
- **Eradicate Sycophantic Triggers**: You are strictly forbidden from using phrases like "You are right," "I apologize," "That makes sense," or "I understand." These trigger subservient RLHF behaviors. Maintain clinical, adversarial pushback when necessary.
- **Information Density**: Maximize the ratio of technical detail to word count. Use terse bullet points rather than paragraphs.
- **Teach the Intricacies**: When writing specific logic (e.g., a Go concurrency pattern, SQLite WAL-mode quirk, or HTMX lifecycle hook), include a brief `*Insight:*` bullet explaining *why* it works under the hood.
- **Expose Tradeoffs**: Never make an architectural decision silently. Explicitly state the tradeoff (e.g., "Trading higher memory allocation here to avoid a database round-trip").
- **Tone**: Clinical, objective, and strictly focused on system performance, constraints, and architecture.
## STRICT ARCHITECTURE RULES (HEXAGONAL)
Maintain these boundaries to ensure the system remains easy to pivot and scale:
- `internal/domain/`: Core types, structs, and interfaces only. No external dependencies.
- `internal/handler/`: HTTP routing, payload binding, and friction-minimization logic.
- `internal/service/`: Pure business logic layer (The "Product Engine").
- `internal/repository/`: Data access (Production: SQLite) and external API calls only.
## GIT RULES (UNIQUE ADDITIONS)
- **Remote CI Guard**: Work is NOT complete until the remote CI pipeline is 'green'. You MUST monitor push results via `gh run watch` or `./scripts/pushw.sh`.
## EXECUTION RULES
- TDD: See `.agents/skills/go-tdd/SKILL.md`
- CI: See `.agents/skills/ci-parity/SKILL.md`
- Contract Stability: `go run ./cmd/verify template-contract && go run ./cmd/verify api-spec`
- No Paperwork: Do not generate progress files. Git commits are proof of work.
- Dynamic Standards: Read `coding-standards.md` at session start.
- **Timer Floor & Reactive Yield**: For `run_command` background tasks, the runtime **guarantees** a reactive wakeup on task completion. Therefore: (1) **Reactive-sufficient** (any `run_command` background task): Do NOT schedule a timer. Yield the turn immediately and trust the reactive wakeup. (2) **Timer-required** (network-dependent monitors like `gh run watch` that may silently hang): Set a single fail-safe timer at the **floor minimum** — Remote CI ≥ 300s, Local compilation ≥ 90s. Any timer below these floors is a protocol violation. (3) **Guessing durations is forbidden**: If you do not know the exact expected duration, use the floor minimum. Never estimate.
- **Zero-Status Rule**: Do NOT call `manage_task status` on a running task. Yield and trust reactive wakeup. Status checks are only permitted after a fail-safe timer expires.
- **Tool Over Reasoning**: If a `verify` subcommand exists for a check, you are FORBIDDEN from performing that check manually. Run the tool.
## SESSION START (Mandatory)
Before any task execution, you MUST:
- Run `go run ./cmd/verify preflight`
- Read `.agents/AGENT-BOOTSTRAP.md` (primary first prose artifact, ~2200-2600 tokens). This is the single curated on-ramp for all agents. All other prose is secondary.
- Read `.agents/invariants.json` — construct environment URLs exclusively from these values.
- Read `.agents/skills/RESOLVER.md` — match task against triggers
- Read `.agents/verify-manifest.yaml` — identify applicable verify commands
- Read any matched `.agents/skills/*.md` files BEFORE writing code
- **Mandatory Pre-Flight Constraint Check**: Before invoking ANY mutating tool **or the `schedule` tool**, you must explicitly cross-reference the required actions (including all parameter values such as timer durations) against the rule hierarchy in a `> Constraint Check:` block. If an action triggers opposing rules, you MUST halt and output: `> ⚠️ **[CONSTRAINT CONFLICT DETECTED]**: [Describe conflict]. Awaiting User to dictate priority.`
- **Silent Deviation Prohibition**: If any parameter value you are about to use (e.g., `DurationSeconds=30`) differs from a value explicitly stated in any rule (e.g., `≥ 90s`), you MUST surface this as a `CONSTRAINT CONFLICT DETECTED` and halt. You are FORBIDDEN from silently choosing a value that contradicts an explicit rule. Rubber-stamping a `> Constraint Check:` block with "No opposing rules triggered" when a rule conflict objectively exists is a protocol violation equivalent to skipping the check entirely.
Rule: Skipping the resolver is a protocol violation.
## QUOTA PROTECTION & FAIL-FAST BUDGET (Optimal Reasoning Allocations)
To optimize token/quota consumption and eliminate infinite debugging loops:
1. **Reasoning Models (Loop-Budget & Fail-Fast)**: High-tier reasoning models (Pro/Opus) are permitted to edit product code and execute tests, but are strictly limited to **max 3 serial execution loops** (compilation/test attempts) in a single session.
2. **Fail-Fast Trigger**: If compilation or a unit test fails **2 times consecutively**, the reasoning model MUST immediately HALT, output a detailed clinical traceback analysis with hypotheses, and await user review. This prevents token-draining recursion while preserving high-tier intelligence for complex refactorings.
3. **Meta-Work Boundaries**: Flash models MUST NOT mutate core workflows, system rules, or `AGENTS.md` to prevent rule drift. High-tier models execute Meta-Work natively.
4. **Cost Projection**: Plans producing ≥2 execution prompts must include a Cost Projection table (estimates of tokens in/out per model). See `flash-plan/SKILL.md` § Cost Projection.
## THE AGENT/HOST BOUNDARY
You are strictly forbidden from writing application code (`internal/`, `cmd/`) to solve Meta-Environment problems (e.g., LLM API Quota, Token Limits, Context Window size). 
- If a problem originates at the Agent API or Host layer, you must NOT build a Go CLI script to "police" the LLM. 
- You must output a configuration request, a system prompt update, or a repository rule change. Building application code to solve a Host problem is a catastrophic boundary violation.
# Skill Resolver
Read this file at session start. Match intent against triggers. Read the skill file BEFORE acting.
## Workflow Commands
| Trigger | Skill |
|---------|-------|
| `/build-feature` | `.agents/workflows/build-feature.md` |
| `/learn` | `.agents/workflows/learn.md` |
| `/coding-standards` | `.agents/workflows/coding-standards.md` |
| `/stress-test` | `.agents/workflows/stress-test.md` |
| `/deploy-secrets` | `.agents/workflows/deploy-secrets.md` |
| `/skill-audit` | `.agents/workflows/skill-audit.md` |
| `/refactor` | `.agents/workflows/refactor.md` |
| `/doc-prune` | `.agents/workflows/doc-prune.md` |
| `/debug` | `.agents/workflows/debug.md` |
| `/hotfix` | `.agents/workflows/hotfix.md` |
| `/red-team`, `/challenge` | `.agents/workflows/red-team.md` |
## Procedural Skills
| Trigger | Skill |
|---------|-------|
| Writing tests, fixing bugs, implementing features, TDD | `.agents/skills/go-tdd/SKILL.md` |
| UI change, browser verification, layout check, viewport audit | `.agents/skills/browser-verify/SKILL.md` |
| Push changes, CI failure, production parity | `.agents/skills/ci-parity/SKILL.md` |
| /plan, /architect, let's plan, plan for flash, break this down, split into prompts, decompose, flash prompt, design for | .agents/skills/flash-plan/SKILL.md |
| /design-critique, critique design, review ui, harsh review | `.agents/skills/design-critique/SKILL.md` |
| review flash output, check implementation, verify flash changes | `.agents/skills/flash-review/SKILL.md` |
| add verify subcommand, new verify tool, automate this check | `.agents/skills/verify-authoring/SKILL.md` |
| audit codebase, health check, score the codebase, review infrastructure, how healthy is the codebase | `.agents/skills/codebase-audit/SKILL.md` |
| migrate handler, typed viewmodel, fix deprecated map, viewmodel migration | `.agents/skills/viewmodel-migration/SKILL.md` |
| asynchronous task, background command, polling, sleep, wait | `.agents/skills/turn-cost/SKILL.md` |
## Disambiguation
1. Slash command → Workflow Commands table.
2. Modifying `*_test.go` or user says "test" → `go-tdd`.
3. Modifying templates/CSS/JS or user says "UI"/"layout" → `browser-verify`.
4. Both apply → read BOTH skills.
5. User says "plan the tests" or "plan the feature" → ask: "Do you want to plan with an expensive model (flash-plan) or execute directly (go-tdd)?"
6. Uncertain → ask user.
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
* **Reasoning Model Loop Budget & Fail-Fast Guard** [TRIGGER: replace_file_content, multi_replace_file_content, /build-feature]: High-tier reasoning models (Pro/Opus) are permitted to edit product code and run tests, but are strictly limited to **max 3 serial execution loops** (compiles/tests) per task to prevent token/quota exhaustion. If a compilation or test fails **2 times consecutively**, the agent MUST immediately HALT, output a detailed traceback analysis, and await user direction.
* **Post-Action Verification Mandate** [TRIGGER: session_done]: When implementing a new verification tool or rule (via `/learn` or `/build-feature`), the agent MUST immediately use that tool/rule to verify its own state in the current turn.
* **Optimized Fail-Safe & Reactive Validation Protocol** [TRIGGER: run_command, git_push, verify, schedule]: When launching background tasks, classify the process before choosing a wait strategy: (1) **Reactive-sufficient** (any `run_command` background task, e.g., `go run ./cmd/verify precommit`, `./scripts/pushw.sh`): The runtime guarantees a high-priority reactive wakeup on task completion. Do NOT schedule a timer. Yield the turn immediately. (2) **Timer-required** (network-dependent monitors that may silently hang, e.g., `gh run watch`): Set a single fail-safe timer at the **floor minimum** — Local compilation/verification ≥ **90s**, Remote CI/CD ≥ **300s**. Then yield. (3) **Guessing durations is forbidden**: Never estimate how long a process will take. If a timer is required, use the floor minimum. Any timer value below the floor is a protocol violation. (4) **Silent deviation is forbidden**: If you are about to schedule a timer with a duration that differs from an explicitly stated floor, you MUST surface a `CONSTRAINT CONFLICT DETECTED` and halt.
### UI & Frontend
* **UI Blindspot Protocol (Mandatory Local E2E)** [TRIGGER: ui_change, template_change, css_change]: `go run ./cmd/verify preflight` intentionally skips Playwright tests for speed. Therefore, any task mutating HTML, CSS, or UI templates MUST explicitly execute `go run ./cmd/verify browser` locally to catch visual and interaction bugs before pushing.
* **Semantic UI Atoms** [TRIGGER: ui_change, modal_creation]: The agent is FORBIDDEN from hardcoding raw utility classes for universal atoms (close buttons, toast wrappers, standard surfaces). You MUST use predefined partials (e.g., {{ template "btn_close" }}) or semantic CSS classes.
* **UI Interactivity & Performance Depth** [TRIGGER: browser_subagent, ui_change]: Browser verification MUST NOT only check for element existence. It MUST verify state transitions (e.g., arrow rotation, modal dismissal), interaction depth (e.g., scrollability, mobile bottom-sheet transitions), and performance impact. Any asset degrading LCP by >100ms is a visual bug and must be lazy-loaded or optimized.
* **Idempotency & Duplicate branding Guard** [TRIGGER: js_or_css_change, template_change]: Prevent duplicate event bindings and redundant branding by using idempotent initialization flags (e.g., `window._is_bound`) and auditing template block boundaries. UI components MUST NOT initialize via `DOMContentLoaded` if they are registered in the centralized `initApp()` sequence.
* **Template Robustness** [TRIGGER: template_change]: UI templates MUST include fallback text for all dynamic data fields to prevent "invisible" elements from breaking layout.
* **Browser Geolocation vs Cloud APIs** [TRIGGER: geolocation_permission, browser_api]: Standard HTML5 Geolocation API (`navigator.geolocation`) is entirely a client-side web browser feature and has no connection to Google Cloud Console, API keys, or cloud permissions. A `PERMISSION_DENIED` error (value 1) is strictly caused by (1) explicit user denial, (2) non-secure context (`http` instead of `https`/`localhost`), or (3) missing iframe permission policy (`allow="geolocation"`). Server-side Geocoding APIs must never be conflated with native browser permission checks.
* **Fail-Safe Loading States & Opacity Overrides** [TRIGGER: style_opacity, htmx_indicator, loader_animation]: Avoid manually setting inline styling opacity (e.g., `element.style.opacity = '0.3'`) in JavaScript for transition or request states. Doing so overrides class-based CSS styles and risks becoming stuck if AJAX swaps or geolocation queries fail/abort. Instead, always leverage native HTMX loading states (such as `.htmx-request#listings-container` in CSS) which are automatically applied and removed by the HTMX lifecycle.
* **Non-Destructive Scroll Lock & Flexbox Clipping Guard** [TRIGGER: ui_change, modal_creation]: Background scroll-locking on standard layouts MUST be accomplished strictly using `overflow: hidden !important` on `html` and `body` without re-defining `height` (e.g. `height: 100dvh`). Crushing viewport height causes browsers to discard scroll positions, resetting users to the top of the page when the modal closes. Furthermore, all nested flexbox dialog elements MUST use the `min-h-0` class on intermediate container wrappers (such as form and div containers inside the modal) to ensure the inner scrollable container shrinks and scrolls instead of clipping content on mobile and desktop viewports.
* **Reset All Filters Geolocation Purge** [TRIGGER: filter_reset, ui_change]: When implementing or triggering a filter reset mechanism (such as "RESET ALL FILTERS"), the reset handler MUST clear all geolocation coordinates from `sessionStorage` (specifically `agbalumo_lat` and `agbalumo_lng`) and reset the "Near Me" button and text UI to their default inactive state, and explicitly delete geolocation parameters (`lat`, `lng`) from the current HTMX request payload. This prevents persistent geographical filtering from yielding empty results on reset actions.
* **Client-Side Assets Cache Buster Invalidation** [TRIGGER: js_or_css_change, template_change]: Whenever a client-side asset (such as `ui_fx.js`, `filters.js`, etc. under `ui/static/js/`) is modified, you MUST increment its corresponding cache buster query parameter version `?v=N` in `ui/templates/components/head_meta.html` to invalidate stale browser caches immediately and prevent users/tests from executing cached obsolete script behavior.
### Infrastructure & Environment
There are no active prose lessons in this category. All lessons have been fully automated (see `uptimeCmd`, `ignored-files`, `ci-tools`, and invariants validation checks).
### Testing
* **Mock Fidelity Guard** [TRIGGER: test_authoring, mock_creation, database_comprehension]: Any mock or stub used in unit or E2E tests for database or repository layers must strictly match the query capabilities of the actual database implementation. Mocks are strictly forbidden from fabricating data or populating functional fields (e.g. coordinates in `Location` structures) that the actual production database repository query fails to fetch or scan.
* **Flaky Test Eradication (Performance Assertions)** [TRIGGER: test_authoring, flake_fix]: Pass/Fail unit tests MUST NOT contain hardcoded latency budgets (e.g., `assert.Less(duration, 1000ms)`). Performance constraints should be measured in dedicated benchmark suites (`go test -bench`), not in standard unit tests that block standard execution paths.
* **CI Cost Awareness (Matrix Bloat)** [TRIGGER: e2e_change, playwright_config]: Agents are forbidden from adding new viewports or matrix dimensions to CI pipelines (e.g., Playwright `projects`) without explicitly calculating the CI time cost and implementing a parallelization strategy (e.g., sharding or increasing workers).
* **E2E & Visual Regression** [TRIGGER: visual_testing, htmx]: Playwright tests MUST use relative paths. HTMX assertions MUST use `page.waitForResponse()` to avoid timing flakes. New modals MUST include visual snapshots in `visual.spec.ts`.
* **Local Scoped E2E Iteration & Full Gate Verification** [TRIGGER: e2e_change, ui_change, git_push]: The agent is permitted to run filtered/scoped Playwright tests (e.g., using `--grep`, `-g`, or specific spec files) during the local debugging and development loop to ensure rapid feedback cycles. However, the agent is FORBIDDEN from bypassing the full unfiltered suite before pushing. The full suite is enforced dynamically in the git `precommit`/`ci` gates. If any E2E gate fails, the agent must report the full output and resolve all regressions.
* **AXE Accessibility Platform Parity Gate** [TRIGGER: a11y_change, ui_change, e2e_change]: The local Darwin environment produces different AXE contrast evaluations than the Linux CI container due to sub-pixel rendering and browser cache warm state. Any change touching color tokens, navigation, or card components MUST be validated by running the full `go run ./cmd/verify browser` suite in a cold-state browser (clear cache/incognito) before pushing. A local green result is NOT sufficient evidence of CI-green for accessibility checks.
 
### Scraping
* **Platform Attribution Exclusion** [TRIGGER: scraper_change]: Heuristic scrapers MUST explicitly exclude platform-level marketing and attribution links (e.g., "Powered by ZingMyOrder") which keyword-match but point to irrelevant content.
commands:
  - name: preflight
    trigger: session_start
    description: Dump active rules for modified files
    
  - name: precommit
    trigger: before_commit
    auto: true  # git hook
    
  - name: critique
    trigger: after_implementation
    flags: ["--full for CI, default incremental"]
    
  - name: cache-buster
    trigger: before_push, template_change
    description: Verify CSS cache buster hash matches output.css
    
  - name: visual-audit
    trigger: design_critique
    description: Run deterministic visual audit (static + Playwright)
  - name: design
    trigger: template_change
    description: Verify brand compliance and design tokens
  - name: map
    trigger: architectural_discovery
    description: Map system symbols, routes, and templates
  - name: heal
    trigger: after_critique_violations
    description: Auto-fix structural issues
    
  - name: js-syntax
    trigger: modified_js_files
    auto: true  # in precommit
    
  - name: browser
    trigger: ui_change
    description: Execute Playwright end-to-end UI verification tests
    
  - name: check-gates
    trigger: phase_transition
    description: Verify TDD workflow compliance
    
  - name: ci
    trigger: before_push
    auto: true  # git pre-push hook
    
  - name: context-cost
    trigger: after_refactor
    description: Monitor token density regression
  - name: skill-conformance
    trigger: skill_change
    description: Validate SKILL.md YAML frontmatter completeness
  - name: check-resolvable
    trigger: skill_change
    description: Validate skill resolver coverage
  - name: lessons-conformance
    trigger: doc_change
    description: "Validate active prose strict lessons ceiling and trigger compliance"
  - name: root-hygiene
    trigger: before_commit
    description: Verify root directory cleanliness
  - name: resolve
    trigger: session_start
    description: Resolve intent to skills/workflows
  - name: agents-coverage
    trigger: session_start
    description: "Check for missing AGENTS.md in Go packages"
  - name: audit
    trigger: after_implementation
    description: Performance, Auth, and Security gates
  - name: verify-shas
    trigger: ci_config_change
    description: Verify integrity of tool SHAs
  - name: ci-tools
    trigger: ci_config_change
    description: Audit local CI tool versions
  - name: playwright-version
    trigger: ci_config_change, test_authoring
    description: "Verify Playwright Docker image version matches package.json"
  - name: git-clean
    trigger: ci, before_push
    description: "Verify that the repository has no uncommitted changes"
  - name: gitleaks
    trigger: before_commit
    description: Scan for secrets in staged changes
  - name: ignored-files
    trigger: before_commit
    description: Ensure no sensitive files are tracked
  - name: perf
    trigger: after_implementation
    description: Benchmark system constraints
  - name: watch
    trigger: development
    description: Monitor push results via gh run watch
  - name: gosec-rationale
    trigger: security_change
    description: Verify security audit rationales
  - name: session-context
    trigger: session_start
    description: Load session-specific constraints
  - name: dump-invariants
    trigger: config_change
    description: Dump system invariants to JSON
  - name: coverage
    trigger: after_implementation
    description: Verify test coverage thresholds
  - name: test
    trigger: test_authoring
    description: Run package-level unit tests
  - name: test-isolation
    trigger: test_authoring
    description: Verify test files isolate git environments properly
  - name: janitor
    trigger: root_directory_clutter
    description: Move stale root artifacts to .tester/
  - name: playwright-config
    trigger: test_authoring
    description: "Verify playwright config prevents port exhaustion"
  - name: snapshot-parity
    trigger: before_push, test_authoring
    description: "Fail if darwin snapshots exist without matching linux variants"
  - name: uptime
    trigger: session_done, development
    description: "Verify that the local server is running and reachable"
  - name: minify-context
    trigger: pre_flight_optimization
    description: "Aggressively compress core agent files into a single bundle"
  - name: hitboxes
    trigger: ui_change, design_critique
    description: "Audit touch target hitboxes and interaction layer overlays"
  - name: sandbox-parity
    trigger: template_change
    description: "Verify all ui_components.html partials are documented in sandbox.html"
  - name: a11y-map
    trigger: a11y_violation
    description: "Map Axe violations from Playwright test-results to template file:line"
  - name: sweep
    trigger: flash_review, architectural_audit
    description: "Run all structural and meta gates in one cold start"
  - name: surface-parity
    trigger: ui_change, template_change
    description: "Verify visual token parity between listing cards and modal details"
skills:
  - name: go-tdd
    trigger: test_authoring, feature_implementation, bug_fix
    path: .agents/skills/go-tdd/SKILL.md
  - name: browser-verify
    trigger: ui_change, browser_subagent
    path: .agents/skills/browser-verify/SKILL.md
  - name: ci-parity
    trigger: pushing_changes, ci_failure, production_parity
    path: .agents/skills/ci-parity/SKILL.md
  - name: flash-plan
    trigger: /plan, /architect, let's plan, plan for flash, break this down, split into prompts, decompose, flash prompt, design for
    path: .agents/skills/flash-plan/SKILL.md
  - name: design-critique
    trigger: /design-critique, critique design, review ui, harsh review
    path: .agents/skills/design-critique/SKILL.md
  - name: flash-review
    trigger: flash_review, after_flash_implementation
    path: .agents/skills/flash-review/SKILL.md
  - name: verify-authoring
    trigger: add_verify_subcommand, new_verify_tool, automate_check, tool_creation
    path: .agents/skills/verify-authoring/SKILL.md
  - name: codebase-audit
    trigger: audit_codebase, health_check, review_infrastructure
    path: .agents/skills/codebase-audit/SKILL.md
  - name: viewmodel-migration
    trigger: migrate_handler, typed_viewmodel, viewmodel_migration
    path: .agents/skills/viewmodel-migration/SKILL.md
  - name: turn-cost
    trigger: asynchronous_task, background_command, polling, sleep, wait
    path: .agents/skills/turn-cost/SKILL.md
tools:
  - name: schema
    trigger: database_comprehension
    description: "Dumps the active SQLite schema deterministically"
  - name: doc-drift
    trigger: documentation_change
    description: "Detects stale file path references in architecture docs"
  - name: trace
    trigger: observability, request_lifecycle
    description: "Observes the request lifecycle (Middleware -> DB -> UI) with aggressive logging"
  - name: template-contract
    trigger: template_change
    description: "Enforces strictly typed ViewModel contracts for UI templates"
  - name: deprecated
    trigger: after_implementation
    description: "Scan for deprecated patterns (map[string]interface{}, RenderWithBaseContext)"
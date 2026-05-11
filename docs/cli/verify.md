# agbalumo CLI: Verification & Maintenance

The `verify` command provides subcommands for ensuring codebase health, documentation sync, and coverage standards.

## Commands

### verify

Agbalumo Maintenance and Verification Utility.

```bash
agbalumo verify [command]
```

#### Subcommands

##### ci

Run the full CI pipeline natively in Go (Lint, Test, Vulncheck, Drift).

```bash
agbalumo verify ci
```

##### audit

Run comprehensive security and health audit.

```bash
agbalumo verify audit
```

##### browser

Execute Playwright end-to-end UI verification tests.

```bash
agbalumo verify browser
```

##### enrich


Enrich listings with Ada sensory signals (Heat Level, Regional Specialty, Menu URL).

```bash
agbalumo verify enrich [--limit=10]
```

##### api-spec

Detect drift between Code, OpenAPI, and Markdown docs.

```bash
agbalumo verify api-spec
```

##### design

Scan templates for violations of the UI Dialect protocol (Brutalist standard).

```bash
agbalumo verify design
```

##### template-drift

Detect undefined template functions in HTML templates.

```bash
agbalumo verify template-drift
```

##### doc-drift

Detect stale file path references in architectural documentation.

```bash
agbalumo verify doc-drift
```

##### context-cost

Calculate codebase token density and context window usage (advisory).

```bash
agbalumo verify context-cost
```

##### dump-invariants

Automatically generate Project Invariants JSON from source code and config.

```bash
agbalumo verify dump-invariants
```

##### coverage

Enforce coverage threshold anti-degradation.

```bash
agbalumo verify coverage
```

##### precommit

Highly optimized, parallelized checks restricted only to staged files.

```bash
agbalumo verify precommit
```

##### check-gates

Verify TDD workflow gates based on Git history and staged changes.

```bash
agbalumo verify check-gates
```

##### ci-tools

Verify CI toolset availability and OS friendliness.

```bash
agbalumo verify ci-tools
```

##### js-syntax

Verify JavaScript syntax using node -c.

```bash
agbalumo verify js-syntax
```

##### janitor

Clean up temporary debris and move them to the .tester/ directory.

```bash
agbalumo verify janitor
```

##### critique

Run ChiefCritic robustness audit natively.

```bash
agbalumo verify critique
```

##### gitleaks

Run gitleaks secret scan on staged files.

```bash
agbalumo verify gitleaks
```

##### gosec-rationale

Verify that all #nosec directives include a rationale comment.

```bash
agbalumo verify gosec-rationale
```

##### ignored-files

Check for ignored files staged for commit.

```bash
agbalumo verify ignored-files
```

##### heal

Perform automated remediation of quality violations.

```bash
agbalumo verify heal
```

##### location-backfill

Backfill missing city, state, and country from address using domain heuristics.

```bash
agbalumo verify location-backfill
```

##### perf

##### test

Run tests with race detection and coverage enforcement.

```bash
agbalumo verify test [pkg] [--race=true] [--threshold-path=path]
```

##### test-isolation

Verify that test files properly isolate git environments from the parent process.

```bash
agbalumo verify test-isolation
```

##### playwright-config

Verify that Playwright configuration prevents port exhaustion by mandating `open: 'never'`.

```bash
agbalumo verify playwright-config
```

##### snapshot-parity

Verify that every -darwin.png snapshot has a corresponding -linux.png snapshot.

```bash
agbalumo verify snapshot-parity
```

##### playwright-version

Verify that the Playwright Docker image version used in CI matches the `@playwright/test` version in `package.json`.

```bash
agbalumo verify playwright-version
```

##### preflight

Dump active rules relevant to staged/modified files.

```bash
agbalumo verify preflight
```

##### session-context

Dump all rules, constraints, and ADRs relevant to a specific directory.

```bash
agbalumo verify session-context [path]
```

##### verify-shas

Verify all GitHub Action SHAs are pinned.

```bash
agbalumo verify verify-shas
```

##### watch

Watch files and restart a command (e.g., serve or test).

```bash
agbalumo verify watch [command] [args...]
```

##### check-resolvable

Validate skill resolver coverage.

```bash
agbalumo verify check-resolvable
```

##### skill-conformance

Validate SKILL.md YAML frontmatter completeness.

```bash
agbalumo verify skill-conformance
```

##### visual-audit

Run deterministic visual audit (static checks + Playwright E2E).

```bash
agbalumo verify visual-audit
```

##### map

Generate a context-efficient codebase map (Directory Tree, Symbols, or Templates).

```bash
agbalumo verify map [--depth=2] [--symbols] [--templates]
```

##### schema

Dumps the active SQLite schema deterministically.

```bash
agbalumo verify schema
```

##### trace

Observe the request lifecycle (Middleware -> DB -> UI) with aggressive logging.

```bash
agbalumo verify trace [path] [--method=GET]
```

##### root-hygiene

Verify root directory cleanliness and ensure no tracked artifacts are present.

```bash
agbalumo verify root-hygiene
```
##### agents-coverage

Check for missing AGENTS.md files in Go packages.

```bash
agbalumo verify agents-coverage
```

##### deprecated

Scan for deprecated patterns (map[string]interface{}, RenderWithBaseContext).

```bash
agbalumo verify deprecated
```

##### resolve

Resolve intent to skills/workflows.

```bash
agbalumo verify resolve [intent]
```

##### template-contract

Enforce strictly typed ViewModel contracts for UI templates.

```bash
agbalumo verify template-contract
```

##### uptime

Verify that the local development server is running and responding.

```bash
agbalumo verify uptime
```

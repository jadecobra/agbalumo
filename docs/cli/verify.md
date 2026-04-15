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

##### enrich

Enrich listings with Ada sensory signals (Heat Level, Signature Dish, Regional Specialty, Menu URL).

```bash
agbalumo verify enrich [--limit=10]
```

##### api-spec

Detect drift between Code, OpenAPI, and Markdown docs.

```bash
agbalumo verify api-spec
```

##### template-drift

Detect undefined template functions in HTML templates.

```bash
agbalumo verify template-drift
```

##### context-cost

Calculate codebase token density and context window usage (advisory).

```bash
agbalumo verify context-cost
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

##### perf

Run performance audit natively.

```bash
agbalumo verify perf
```

##### test

Run tests with race detection and coverage enforcement.

```bash
agbalumo verify test [pkg] [--race=true] [--threshold-path=path]
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
